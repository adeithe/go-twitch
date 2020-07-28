package cmd

// IRCCommand is the command sent in an IRCMessage
type IRCCommand string

const (
	PrivMessage     IRCCommand = "PRIVMSG"
	ClearChat       IRCCommand = "CLEARCHAT"
	ClearMessage    IRCCommand = "CLEARMSG"
	HostTarget      IRCCommand = "HOSTTARGET"
	Notice          IRCCommand = "NOTICE"
	Reconnect       IRCCommand = "RECONNECT"
	RoomState       IRCCommand = "ROOMSTATE"
	UserNotice      IRCCommand = "USERNOTICE"
	GlobalUserState IRCCommand = "GLOBALUSERSTATE"
	UserState       IRCCommand = "USERSTATE"
	Ping            IRCCommand = "PING"
	Pong            IRCCommand = "PONG"
	Ready           IRCCommand = "376"
)