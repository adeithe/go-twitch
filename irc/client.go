package irc

import (
	"sync"
	"time"
)

// Client is a sharded connection to the Twitch IRC service
type Client struct {
	length        int
	shards        map[int]*Conn
	awaitingClose int

	onShardMessage              []func(int, ChatMessage)
	onShardServerNotice         []func(int, ServerNotice)
	onShardLatencyUpdate        []func(int, time.Duration)
	onShardChannelJoin          []func(int, string, string)
	onShardChannelLeave         []func(int, string, string)
	onShardChannelUpdate        []func(int, RoomState)
	onShardChannelUserNotice    []func(int, UserNotice)
	onShardChannelMessageDelete []func(int, ChatMessageDelete)
	onShardChannelBan           []func(int, ChatBan)
	onShardRawMessage           []func(int, Message)
	onShardReconnect            []func(int)
	onShardDisconnect           []func(int)

	mx sync.Mutex
	wg sync.WaitGroup
}

// IClient is a generic IRC shard provider
type IClient interface {
	SetMaxChannelsPerShard(int)
	GetNextShard() (*Conn, error)
	GetShard(int) (*Conn, error)
	GetChannel(string) (RoomState, bool)

	Join(...string) error
	Leave(...string) error
	Close()

	OnShardMessage(func(int, ChatMessage))
	OnShardServerNotice(func(int, ServerNotice))
	OnShardLatencyUpdate(func(int, time.Duration))
	OnShardChannelJoin(func(int, string, string))
	OnShardChannelLeave(func(int, string, string))
	OnShardChannelUpdate(func(int, RoomState))
	OnShardChannelUserNotice(func(int, UserNotice))
	OnShardChannelMessageDelete(func(int, ChatMessageDelete))
	OnShardChannelBan(func(int, ChatBan))
	OnShardRawMessage(func(int, Message))
	OnShardReconnect(func(int))
	OnShardDisconnect(func(int))
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
	for id, conn := range client.shards {
		if len(conn.channels) < client.length {
			shardID = id
			break
		}
	}
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
	if client.shards[id] == nil {
		conn := &Conn{isShard: true}
		client.addEventHandlers(id, conn)
		client.shards[id] = conn
	}
	shard := client.shards[id]
	return shard, nil
}

// GetChannel returns true if any shard is listening to the provided channel along with its RoomState
func (client *Client) GetChannel(channel string) (RoomState, bool) {
	client.mx.Lock()
	defer client.mx.Unlock()
	for _, shard := range client.shards {
		if state, ok := shard.GetChannel(channel); ok {
			return state, ok
		}
	}
	return RoomState{}, false
}

// Join attempts to join a channel on an available shard
func (client *Client) Join(channels ...string) error {
	for _, channel := range channels {
		if _, ok := client.GetChannel(channel); !ok {
			shard, err := client.GetNextShard()
			if err != nil {
				return err
			}
			if err := shard.Join(channel); err != nil {
				return err
			}
		}
	}
	return nil
}

// Leave attempts to leave a channel
func (client *Client) Leave(channels ...string) error {
	client.mx.Lock()
	for _, shard := range client.shards {
		for _, channel := range channels {
			if _, ok := shard.GetChannel(channel); ok {
				if err := shard.Leave(channel); err != nil {
					return err
				}
			}
		}
	}
	client.mx.Unlock()
	return nil
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

// OnShardMessage event called after a shard receives a chat message
func (client *Client) OnShardMessage(f func(int, ChatMessage)) {
	client.onShardMessage = append(client.onShardMessage, f)
}

// OnShardServerNotice event called when the IRC server sends a notice message
func (client *Client) OnShardServerNotice(f func(int, ServerNotice)) {
	client.onShardServerNotice = append(client.onShardServerNotice, f)
}

// OnShardLatencyUpdate event called after a shards latency is updated
func (client *Client) OnShardLatencyUpdate(f func(int, time.Duration)) {
	client.onShardLatencyUpdate = append(client.onShardLatencyUpdate, f)
}

// OnShardChannelJoin event called after a user joins a chatroom
func (client *Client) OnShardChannelJoin(f func(int, string, string)) {
	client.onShardChannelJoin = append(client.onShardChannelJoin, f)
}

// OnShardChannelLeave event called after a user leaves a chatroom
func (client *Client) OnShardChannelLeave(f func(int, string, string)) {
	client.onShardChannelLeave = append(client.onShardChannelLeave, f)
}

// OnShardChannelUpdate event called after a chatrooms state has been modified
func (client *Client) OnShardChannelUpdate(f func(int, RoomState)) {
	client.onShardChannelUpdate = append(client.onShardChannelUpdate, f)
}

// OnShardChannelUserNotice event called after a generic user event occurrs in a channels chatroom
func (client *Client) OnShardChannelUserNotice(f func(int, UserNotice)) {
	client.onShardChannelUserNotice = append(client.onShardChannelUserNotice, f)
}

// OnShardChannelMessageDelete event called after a message was deleted in a channels chatroom
func (client *Client) OnShardChannelMessageDelete(f func(int, ChatMessageDelete)) {
	client.onShardChannelMessageDelete = append(client.onShardChannelMessageDelete, f)
}

// OnShardChannelBan event called after a user was banned or timed out in a channels chatoom
func (client *Client) OnShardChannelBan(f func(int, ChatBan)) {
	client.onShardChannelBan = append(client.onShardChannelBan, f)
}

// OnShardRawMessage event called after a shard receives an IRC message
func (client *Client) OnShardRawMessage(f func(int, Message)) {
	client.onShardRawMessage = append(client.onShardRawMessage, f)
}

// OnShardReconnect event called after a shards connection is reopened
func (client *Client) OnShardReconnect(f func(int)) {
	client.onShardReconnect = append(client.onShardReconnect, f)
}

// OnShardDisconnect event called after a shards connection is closed
func (client *Client) OnShardDisconnect(f func(int)) {
	client.onShardDisconnect = append(client.onShardDisconnect, f)
}

//nolint: gocyclo
//gocyclo:ignore
func (client *Client) addEventHandlers(id int, conn *Conn) {
	conn.OnMessage(func(msg ChatMessage) {
		for _, f := range client.onShardMessage {
			go f(id, msg)
		}
	})
	conn.OnServerNotice(func(notice ServerNotice) {
		for _, f := range client.onShardServerNotice {
			go f(id, notice)
		}
	})
	conn.OnLatencyUpdate(func(latency time.Duration) {
		for _, f := range client.onShardLatencyUpdate {
			go f(id, latency)
		}
	})
	conn.OnChannelJoin(func(channel, username string) {
		for _, f := range client.onShardChannelJoin {
			go f(id, channel, username)
		}
	})
	conn.OnChannelLeave(func(channel, username string) {
		for _, f := range client.onShardChannelLeave {
			go f(id, channel, username)
		}
	})
	conn.OnChannelUpdate(func(state RoomState) {
		for _, f := range client.onShardChannelUpdate {
			go f(id, state)
		}
	})
	conn.OnChannelUserNotice(func(notice UserNotice) {
		for _, f := range client.onShardChannelUserNotice {
			go f(id, notice)
		}
	})
	conn.OnChannelMessageDelete(func(delete ChatMessageDelete) {
		for _, f := range client.onShardChannelMessageDelete {
			go f(id, delete)
		}
	})
	conn.OnChannelBan(func(ban ChatBan) {
		for _, f := range client.onShardChannelBan {
			go f(id, ban)
		}
	})
	conn.OnRawMessage(func(msg Message) {
		for _, f := range client.onShardRawMessage {
			go f(id, msg)
		}
	})
	conn.OnReconnect(func() {
		for _, f := range client.onShardReconnect {
			go f(id)
		}
	})
	conn.OnDisconnect(func() {
		for _, f := range client.onShardDisconnect {
			go f(id)
		}
		if client.awaitingClose > 0 {
			client.mx.Lock()
			defer client.mx.Unlock()
			client.awaitingClose--
			delete(client.shards, id)
			client.wg.Done()
		}
	})
}
