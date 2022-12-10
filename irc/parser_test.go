package irc_test

import (
	"fmt"
	"testing"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/stretchr/testify/assert"
)

func BenchmarkParseMessage(b *testing.B) {
	raw := []string{
		":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv JOIN",
		":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv JOIN #jtv",
		"@badge-info=;badges=moments/1;client-nonce=4fb782293442bb3b0df16b4cb5eb21aa;color=#008000;display-name=justinfan16432;emotes=;first-msg=0;flags=;id=6198df9e-77af-4f4f-8d3c-d317802b7c0d;mod=0;returning-chatter=0;room-id=14027;subscriber=0;tmi-sent-ts=1656612693901;turbo=0;user-id=16933;user-type= :justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv PRIVMSG #jtv :Hello",
		"@badge-info=subscriber/2;badges=subscriber/2,no_audio/1;color=#FF0000;display-name=Justinfan16432;emotes=;flags=;id=e1e8d818-3837-4381-b427-c4005ee29ba9;login=justinfan16432;mod=0;msg-id=resub;msg-param-cumulative-months=2;msg-param-months=0;msg-param-multimonth-duration=0;msg-param-multimonth-tenure=0;msg-param-should-share-streak=0;msg-param-sub-plan-name=Channel\\sSubscription\\s(jtv);msg-param-sub-plan=1000;msg-param-was-gifted=false;room-id=14027;subscriber=1;system-msg=Justinfan16432\\ssubscribed\\sat\\sTier\\s1.\\sThey've\\ssubscribed\\sfor\\s2\\smonths!;tmi-sent-ts=1656640802512;user-id=14028;user-type= :tmi.twitch.tv USERNOTICE #jtv :This is a resub message",
	}

	for i, raw := range raw {
		b.Run(fmt.Sprintf("#%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := irc.ParseMessage(raw); err != nil {
					b.Fail()
				}
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		raw      string
		expected *irc.Message
	}{
		{
			":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv JOIN",
			&irc.Message{
				Source: irc.Source{
					Username: "justinfan16432",
					Nickname: "justinfan16432",
					Host:     "justinfan16432.tmi.twitch.tv",
				},
				Command: irc.CMDJoin,
			},
		},
		{
			":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv JOIN #jtv",
			&irc.Message{
				Source: irc.Source{
					Username: "justinfan16432",
					Nickname: "justinfan16432",
					Host:     "justinfan16432.tmi.twitch.tv",
				},
				Command: irc.CMDJoin,
				Params:  []string{"#jtv"},
			},
		},
		{
			"@badge-info=;badges=moments/1;client-nonce=4fb782293442bb3b0df16b4cb5eb21aa;color=#008000;display-name=justinfan16432;emotes=;first-msg=0;flags=;id=6198df9e-77af-4f4f-8d3c-d317802b7c0d;mod=0;returning-chatter=0;room-id=14027;subscriber=0;tmi-sent-ts=1656612693901;turbo=0;user-id=16933;user-type= :justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv PRIVMSG #jtv :Hello",
			&irc.Message{
				Tags: irc.Tags{
					"badge-info":        "",
					"badges":            "moments/1",
					"client-nonce":      "4fb782293442bb3b0df16b4cb5eb21aa",
					"color":             "#008000",
					"display-name":      "justinfan16432",
					"emotes":            "",
					"first-msg":         "0",
					"flags":             "",
					"id":                "6198df9e-77af-4f4f-8d3c-d317802b7c0d",
					"mod":               "0",
					"returning-chatter": "0",
					"room-id":           "14027",
					"subscriber":        "0",
					"tmi-sent-ts":       "1656612693901",
					"turbo":             "0",
					"user-id":           "16933",
					"user-type":         "",
				},
				Source: irc.Source{
					Username: "justinfan16432",
					Nickname: "justinfan16432",
					Host:     "justinfan16432.tmi.twitch.tv",
				},
				Command: irc.CMDPrivMessage,
				Params:  []string{"#jtv"},
				Text:    "Hello",
			},
		},
		{

			"@badge-info=subscriber/2;badges=subscriber/2,no_audio/1;color=#FF0000;display-name=Justinfan16432;emotes=;flags=;id=e1e8d818-3837-4381-b427-c4005ee29ba9;login=justinfan16432;mod=0;msg-id=resub;msg-param-cumulative-months=2;msg-param-months=0;msg-param-multimonth-duration=0;msg-param-multimonth-tenure=0;msg-param-should-share-streak=0;msg-param-sub-plan-name=Channel\\sSubscription\\s(jtv);msg-param-sub-plan=1000;msg-param-was-gifted=false;room-id=14027;subscriber=1;system-msg=Justinfan16432\\ssubscribed\\sat\\sTier\\s1.\\sThey've\\ssubscribed\\sfor\\s2\\smonths!;tmi-sent-ts=1656640802512;user-id=14028;user-type= :tmi.twitch.tv USERNOTICE #jtv :This is a resub message",
			&irc.Message{
				Tags: irc.Tags{
					"badge-info":                    "subscriber/2",
					"badges":                        "subscriber/2,no_audio/1",
					"color":                         "#FF0000",
					"display-name":                  "Justinfan16432",
					"emotes":                        "",
					"flags":                         "",
					"id":                            "e1e8d818-3837-4381-b427-c4005ee29ba9",
					"login":                         "justinfan16432",
					"mod":                           "0",
					"msg-id":                        "resub",
					"msg-param-cumulative-months":   "2",
					"msg-param-months":              "0",
					"msg-param-multimonth-duration": "0",
					"msg-param-multimonth-tenure":   "0",
					"msg-param-should-share-streak": "0",
					"msg-param-sub-plan-name":       "Channel Subscription (jtv)",
					"msg-param-sub-plan":            "1000",
					"msg-param-was-gifted":          "false",
					"room-id":                       "14027",
					"subscriber":                    "1",
					"system-msg":                    "Justinfan16432 subscribed at Tier 1. They've subscribed for 2 months!",
					"tmi-sent-ts":                   "1656640802512",
					"user-id":                       "14028",
					"user-type":                     "",
				},
				Source:  irc.Source{Host: "tmi.twitch.tv"},
				Command: irc.CMDUserNotice,
				Params:  []string{"#jtv"},
				Text:    "This is a resub message",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			test.expected.Raw = test.raw
			message, err := irc.ParseMessage(test.raw)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.expected, message)
			assert.Equal(t, test.raw, message.String())
		})
	}
}

func TestParseMessageError(t *testing.T) {
	tests := []struct {
		raw      string
		expected error
	}{
		{"@", irc.ErrInvalidTags},
		{": JOIN #jtv", irc.ErrInvalidSource},
		{":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv", irc.ErrNoCommand},
		{"@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0", irc.ErrPartialMessage},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			_, err := irc.ParseMessage(test.raw)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}

func BenchmarkParseTags(b *testing.B) {
	raw := []string{
		"@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0",
		"@badge-info=subscriber/27;badges=subscriber/24,turbo/1;color=#F0F0F0;display-name=Kappa;emotes=301445381:3-10;first-msg=0;flags=0-0:P.6,12-12:P.6;id=df900783-2a71-414e-adc0-5fded36c1d55;mod=0;returning-chatter=0;room-id=14027;subscriber=1;tmi-sent-ts=1656612613186;turbo=1;user-id=14028;user-type=",
		"@badge-info=subscriber/2;badges=subscriber/2,no_audio/1;color=#FF0000;display-name=Justinfan16432;emotes=;flags=;id=e1e8d818-3837-4381-b427-c4005ee29ba9;login=justinfan16432;mod=0;msg-id=resub;msg-param-cumulative-months=2;msg-param-months=0;msg-param-multimonth-duration=0;msg-param-multimonth-tenure=0;msg-param-should-share-streak=0;msg-param-sub-plan-name=Channel\\sSubscription\\s(jtv);msg-param-sub-plan=1000;msg-param-was-gifted=false;room-id=14027;subscriber=1;system-msg=Justinfan16432\\ssubscribed\\sat\\sTier\\s1.\\sThey've\\ssubscribed\\sfor\\s2\\smonths!;tmi-sent-ts=1656640802512;user-id=14028;user-type=",
	}

	for i, raw := range raw {
		b.Run(fmt.Sprintf("#%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = irc.ParseTags(raw)
			}
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		raw      string
		expected irc.Tags
	}{
		{
			"@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0",
			irc.Tags{
				"emote-only":     "0",
				"followers-only": "1440",
				"r9k":            "0",
				"room-id":        "14027",
				"slow":           "0",
				"subs-only":      "0",
			},
		},
		{
			"@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0;with-extra-parts==0",
			irc.Tags{
				"emote-only":       "0",
				"followers-only":   "1440",
				"r9k":              "0",
				"room-id":          "14027",
				"slow":             "0",
				"subs-only":        "0",
				"with-extra-parts": "=0",
			},
		},
		{
			"@badge-info=subscriber/27;badges=subscriber/24,turbo/1;color=#F0F0F0;display-name=Kappa;emotes=301445381:3-10;first-msg=0;flags=0-0:P.6,12-12:P.6;id=df900783-2a71-414e-adc0-5fded36c1d55;mod=0;returning-chatter=0;room-id=14027;subscriber=1;tmi-sent-ts=1656612613186;turbo=1;user-id=14028;user-type=",
			irc.Tags{
				"badge-info":        "subscriber/27",
				"badges":            "subscriber/24,turbo/1",
				"color":             "#F0F0F0",
				"display-name":      "Kappa",
				"emotes":            "301445381:3-10",
				"first-msg":         "0",
				"flags":             "0-0:P.6,12-12:P.6",
				"id":                "df900783-2a71-414e-adc0-5fded36c1d55",
				"mod":               "0",
				"returning-chatter": "0",
				"room-id":           "14027",
				"subscriber":        "1",
				"tmi-sent-ts":       "1656612613186",
				"turbo":             "1",
				"user-id":           "14028",
				"user-type":         "",
			},
		},
		{
			"@badge-info=subscriber/2;badges=subscriber/2,no_audio/1;color=#FF0000;display-name=Justinfan16432;emotes=;flags=;id=e1e8d818-3837-4381-b427-c4005ee29ba9;login=justinfan16432;mod=0;msg-id=resub;msg-param-cumulative-months=2;msg-param-months=0;msg-param-multimonth-duration=0;msg-param-multimonth-tenure=0;msg-param-should-share-streak=0;msg-param-sub-plan-name=Channel\\sSubscription\\s(jtv);msg-param-sub-plan=1000;msg-param-was-gifted=false;room-id=14027;subscriber=1;system-msg=Justinfan16432\\ssubscribed\\sat\\sTier\\s1.\\sThey've\\ssubscribed\\sfor\\s2\\smonths!;tmi-sent-ts=1656640802512;user-id=14028;user-type=",
			irc.Tags{
				"badge-info":                    "subscriber/2",
				"badges":                        "subscriber/2,no_audio/1",
				"color":                         "#FF0000",
				"display-name":                  "Justinfan16432",
				"emotes":                        "",
				"flags":                         "",
				"id":                            "e1e8d818-3837-4381-b427-c4005ee29ba9",
				"login":                         "justinfan16432",
				"mod":                           "0",
				"msg-id":                        "resub",
				"msg-param-cumulative-months":   "2",
				"msg-param-months":              "0",
				"msg-param-multimonth-duration": "0",
				"msg-param-multimonth-tenure":   "0",
				"msg-param-should-share-streak": "0",
				"msg-param-sub-plan-name":       "Channel Subscription (jtv)",
				"msg-param-sub-plan":            "1000",
				"msg-param-was-gifted":          "false",
				"room-id":                       "14027",
				"subscriber":                    "1",
				"system-msg":                    "Justinfan16432 subscribed at Tier 1. They've subscribed for 2 months!",
				"tmi-sent-ts":                   "1656640802512",
				"user-id":                       "14028",
				"user-type":                     "",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			tags, err := irc.ParseTags(test.raw)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.expected, tags)
		})
	}
}

func TestParseTagsError(t *testing.T) {
	tests := []struct {
		raw      string
		expected error
	}{
		{"", irc.ErrInvalidTags},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			_, err := irc.ParseTags(test.raw)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}

func BenchmarkParseSource(b *testing.B) {
	raw := []string{
		":tmi.twitch.tv",
		":justinfan16432@justinfan16432.tmi.twitch.tv",
		":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv",
	}

	for i, raw := range raw {
		b.Run(fmt.Sprintf("#%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = irc.ParseSource(raw)
			}
		})
	}
}

func TestParseSource(t *testing.T) {
	tests := []struct {
		raw      string
		expected *irc.Source
	}{
		{
			":tmi.twitch.tv",
			&irc.Source{
				Host: "tmi.twitch.tv",
			},
		},
		{
			":justinfan16432@justinfan16432.tmi.twitch.tv",
			&irc.Source{
				Nickname: "justinfan16432",
				Username: "justinfan16432",
				Host:     "justinfan16432.tmi.twitch.tv",
			},
		},
		{
			":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv",
			&irc.Source{
				Nickname: "justinfan16432",
				Username: "justinfan16432",
				Host:     "justinfan16432.tmi.twitch.tv",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			source, err := irc.ParseSource(test.raw)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.expected, source)
		})
	}
}

func TestParseSourceError(t *testing.T) {
	tests := []struct {
		raw      string
		expected error
	}{
		{"", irc.ErrInvalidSource},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			_, err := irc.ParseSource(test.raw)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}
