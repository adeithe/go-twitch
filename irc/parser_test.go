package irc

import "testing"

func TestParseErrors(t *testing.T) {
	tests := []struct {
		in   string
		want error
	}{
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;user-id=12345;user-type=", ErrPartialMessage},
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;user-id=12345;user-type= :tmi.twitch.tv", ErrNoCommand},
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;user-id=12345;user-type= :tmi.twitch.tv GLOBALUSERSTATE", nil},
		{":testuser!testuser@testuser.tmi.twitch.tv JOIN #channel", nil},
		{"@emote-only=0;followers-only=-1;r9k=0;rituals=0;room-id=123;slow=0;subs-only=1 :tmi.twitch.tv ROOMSTATE #channel", nil},
		{"@badge-info=;badges=;color=;display-name=TestUser;emotes=;flags=;id=abcd123-0123-4abc-defg-1234567890;mod=0;room-id=123;subscriber=0;tmi-sent-ts=1612256273447;turbo=0;user-id=12345;user-type= :testuser!testuser@testuser.tmi.twitch.tv PRIVMSG #channel :this is a message", nil},
		{"@badge-info=;badges=;color=;display-name=TestUser;emote-sets=0;mod=0;subscriber=0;user-type= :tmi.twitch.tv USERSTATE #channel", nil},
		{"@login=testuser;room-id=;target-msg-id=abcd123-0123-4abc-defg-1234567890;tmi-sent-ts=1612256431213 :tmi.twitch.tv CLEARMSG #channel :this is a message", nil},
		{"@ban-duration=60;room-id=123;target-user-id=12345;tmi-sent-ts=1612256572313 :tmi.twitch.tv CLEARCHAT #channel :testuser", nil},
		{"@room-id=123;target-user-id=12345;tmi-sent-ts=1612256572313 :tmi.twitch.tv CLEARCHAT #channel :testuser", nil},
		{"@msg-id=msg_banned :tmi.twitch.tv NOTICE #channel :You are permanently banned from talking in channel.", nil},
		{":testuser@testuser.tmi.twitch.tv PART #channel", nil},
		{":tmi.twitch.tv RECONNECT", nil},
		{":tmi.twitch.tv PING", nil},
	}
	conn := &Conn{}
	for i, test := range tests {
		t.Log(test.in)
		msg, err := NewParsedMessage(test.in)
		if err != nil {
			if test.want != nil {
				if test.want.Error() != err.Error() {
					t.Fatalf("Simulated line #%d failed, got: %s, want: %s", i, err, test.want)
				}
				continue
			}
			t.Fatalf("Simulated line #%d failed, got: %s, want: <nil>", i, err)
		}
		conn.handle(msg)
	}
}
