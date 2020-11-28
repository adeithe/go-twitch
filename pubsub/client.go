package pubsub

import (
	"sync"
	"time"
)

// Client stores data about a PubSub shard manager
type Client struct {
	length        int
	topicsLength  int
	shards        map[int]*Conn
	awaitingClose int

	onShardConnect       []func(int)
	onShardMessage       []func(int, string, []byte)
	onShardLatencyUpdate []func(int, time.Duration)
	onShardReconnect     []func(int)
	onShardDisconnect    []func(int)

	mx sync.Mutex
	wg sync.WaitGroup
}

// IClient interface for methods used by the PubSub shard manager
type IClient interface {
	SetMaxShards(int)
	SetMaxTopicsPerShard(int)
	GetNumShards() int
	GetNumTopics() int
	GetNextShard() (*Conn, error)
	GetShard(int) (*Conn, error)
	Close()

	Listen(string, ...interface{}) error
	ListenWithAuth(string, string, ...interface{}) error
	Unlisten(...string) error

	OnShardConnect(func(int))
	OnShardMessage(func(int, string, []byte))
	OnShardLatencyUpdate(func(int, time.Duration))
	OnShardReconnect(func(int))
	OnShardDisconnect(func(int))
}

var _ IClient = &Client{}

// New PubSub Client
//
// The client uses a sharding system to comply with limits as listed on the Twitch PubSub Documentation.
// Twitch recommends no more than 10 simultaneous shards and no more than 50 topics per shard. These are set by default.
//
// If for any reason a shard attempts to listen to more topics than the server allows, it will attempt to correct for it.
//
// See: https://dev.twitch.tv/docs/pubsub
func New() *Client {
	return &Client{length: 10, topicsLength: 50}
}

// SetMaxShards set the maximum number of shards
//
// Default: 10
func (client *Client) SetMaxShards(max int) {
	if max < 1 {
		max = 10
	}
	client.length = max
}

// SetMaxTopicsPerShard set the maximum number of topics for each shard
//
// Default: 50
func (client *Client) SetMaxTopicsPerShard(max int) {
	if max < 1 {
		max = 50
	}
	client.mx.Lock()
	defer client.mx.Unlock()
	client.topicsLength = max
	for _, shard := range client.shards {
		shard.SetMaxTopics(client.topicsLength)
	}
}

// GetNumShards returns the number of active shards
func (client *Client) GetNumShards() int {
	return len(client.shards)
}

// GetNumTopics returns the number of topics being listened to across all shards
func (client *Client) GetNumTopics() (n int) {
	client.mx.Lock()
	defer client.mx.Unlock()
	for _, shard := range client.shards {
		n += shard.GetNumTopics()
	}
	return
}

// GetNextShard returns the first shard that can accept topics
func (client *Client) GetNextShard() (*Conn, error) {
	client.mx.Lock()
	shardID := len(client.shards)
	for id, conn := range client.shards {
		if conn.GetNumTopics() < conn.length {
			shardID = id
			break
		}
	}
	client.mx.Unlock()
	return client.GetShard(shardID)
}

// GetShard retrieves or creates a shard based on the provided id
func (client *Client) GetShard(id int) (*Conn, error) {
	if client.length < 1 {
		client.length = 10
	}
	if client.topicsLength < 1 {
		client.topicsLength = 50
	}
	if id < 0 || id > client.length-1 {
		return nil, ErrShardIDOutOfBounds
	}
	client.mx.Lock()
	defer client.mx.Unlock()
	if client.shards == nil {
		client.shards = make(map[int]*Conn)
	}
	if client.shards[id] == nil {
		conn := &Conn{length: client.topicsLength}
		conn.OnMessage(func(topic string, data []byte) {
			for _, f := range client.onShardMessage {
				go f(id, topic, data)
			}
		})
		conn.OnPong(func(latency time.Duration) {
			for _, f := range client.onShardLatencyUpdate {
				go f(id, latency)
			}
		})
		conn.OnReconnect(func() {
			for _, f := range client.onShardReconnect {
				go f(id)
			}
		})
		conn.OnDisconnect(func() {
			client.mx.Lock()
			defer client.mx.Unlock()
			for _, f := range client.onShardDisconnect {
				go f(id)
			}
			if client.awaitingClose > 0 {
				client.awaitingClose--
				delete(client.shards, id)
				client.wg.Done()
			}
		})
		client.shards[id] = conn
	}
	shard := client.shards[id]
	if !shard.IsConnected() {
		if err := shard.Connect(); err != nil {
			return nil, err
		}
		defer shard.Ping()
		for _, f := range client.onShardConnect {
			go f(id)
		}
	}
	return shard, nil
}

// Close all active shards
func (client *Client) Close() {
	client.mx.Lock()
	for _, shard := range client.shards {
		client.wg.Add(1)
		client.awaitingClose++
		shard.Close()
	}
	client.mx.Unlock()
	client.wg.Wait()
	client.awaitingClose = 0
}

// Listen to a topic on the best available shard
func (client *Client) Listen(topic string, args ...interface{}) error {
	topic = ParseTopic(topic, args...)
	shard, err := client.GetNextShard()
	if err != nil {
		return err
	}
	client.mx.Lock()
	defer client.mx.Unlock()
	if err := shard.Listen(topic); err != nil {
		if err == ErrShardTooManyTopics {
			client.SetMaxTopicsPerShard(shard.GetNumTopics())
			shard, err := client.GetNextShard()
			if err != nil {
				return err
			}
			return shard.Listen(topic)
		}
		return err
	}
	return nil
}

// ListenWithAuth starts listening to a topic on the best available shard using the provided authentication token
func (client *Client) ListenWithAuth(token string, topic string, args ...interface{}) error {
	topic = ParseTopic(topic, args...)
	shard, err := client.GetNextShard()
	if err != nil {
		return err
	}
	client.mx.Lock()
	defer client.mx.Unlock()
	if err := shard.ListenWithAuth(token, topic); err != nil {
		if err == ErrShardTooManyTopics {
			client.SetMaxTopicsPerShard(shard.GetNumTopics())
			shard, err := client.GetNextShard()
			if err != nil {
				return err
			}
			return shard.ListenWithAuth(token, topic)
		}
		return err
	}
	return nil
}

// Unlisten from the provided topics
//
// Will return the first error that occurs, if any
func (client *Client) Unlisten(topics ...string) error {
	client.mx.Lock()
	defer client.mx.Unlock()
	for _, shard := range client.shards {
		for _, topic := range topics {
			if shard.HasTopic(topic) {
				if err := shard.Unlisten(topic); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// OnShardConnect event called after a shard connects to the PubSub server
func (client *Client) OnShardConnect(f func(int)) {
	client.onShardConnect = append(client.onShardConnect, f)
}

// OnShardMessage event called after a shard gets a PubSub message
func (client *Client) OnShardMessage(f func(int, string, []byte)) {
	client.onShardMessage = append(client.onShardMessage, f)
}

// OnShardLatencyUpdate event called after a shards latency is updated
func (client *Client) OnShardLatencyUpdate(f func(int, time.Duration)) {
	client.onShardLatencyUpdate = append(client.onShardLatencyUpdate, f)
}

// OnShardReconnect event called after a shard reconnects to the PubSub server
func (client *Client) OnShardReconnect(f func(int)) {
	client.onShardReconnect = append(client.onShardReconnect, f)
}

// OnShardDisconnect event called after a shard is disconnected from the PubSub server
func (client *Client) OnShardDisconnect(f func(int)) {
	client.onShardDisconnect = append(client.onShardDisconnect, f)
}
