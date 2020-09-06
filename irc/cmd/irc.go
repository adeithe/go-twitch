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
	UserState       IRCCommand = "USERSTATE"
	GlobalUserState IRCCommand = "GLOBALUSERSTATE"
	Join            IRCCommand = "JOIN"
	Part            IRCCommand = "PART"
	Ping            IRCCommand = "PING"
	Pong            IRCCommand = "PONG"
	Ready           IRCCommand = "376"
)
