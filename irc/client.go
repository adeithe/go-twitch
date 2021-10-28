package irc

import (
	"sync"
)

// Client is a sharded connection to the Twitch IRC service
type Client struct {
	length        int
	shards        map[int]*Conn
	awaitingClose int

	ServerNotice chan ServerNotice
	RawMessage   chan Message

	mx sync.Mutex
	wg sync.WaitGroup
}

// IClient is a generic IRC shard provider
type IClient interface {
	SetMaxChannelsPerShard(int)
	GetNextShard() (*Conn, error)
	GetShard(int) (*Conn, error)

	Close()
}

var _ IClient = &Client{}

// New IRC Client
//
// The client uses a sharding system to allow applications to listen to large numbers of Twitch chatrooms with
// minimized backlogs of message handling. The client will separate channels into groups of 100 by default.
//
// See: https://dev.twitch.tv/docs/irc
func New() *Client {
	return &Client{length: 100}
}

// SetMaxChannelsPerShard sets the maximum number of channels a shard can listen to at a time
//
// Default: 100
func (client *Client) SetMaxChannelsPerShard(max int) {
	if max < 1 {
		max = 100
	}
	client.length = max
}

// GetNextShard returns the first shard that can join channels
func (client *Client) GetNextShard() (*Conn, error) {
	client.mx.Lock()
	shardID := len(client.shards)
	// TODO: Compare shards connected channels to max per shard
	client.mx.Unlock()
	return client.GetShard(shardID)
}

// GetShard retrieves or creates a shard based on the provided id
func (client *Client) GetShard(id int) (*Conn, error) {
	client.mx.Lock()
	defer client.mx.Unlock()
	if id < 0 {
		return nil, ErrShardIDOutOfBounds
	}
	if client.length < 1 {
		client.SetMaxChannelsPerShard(0)
	}
	if client.shards == nil {
		client.shards = make(map[int]*Conn)
	}
	if _, ok := client.shards[id]; !ok {
		conn := &Conn{}
		conn.ServerNotice = client.ServerNotice
		conn.RawMessage = client.RawMessage
		client.shards[id] = conn
	}
	shard := client.shards[id]
	return shard, nil
}

// Close disconnect all active shards
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
