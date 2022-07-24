package irc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		raw      string
		expected *Message
	}{
		{
			":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv JOIN #jtv",
			&Message{
				Source: Source{
					Username: "justinfan16432",
					Nickname: "justinfan16432",
					Host:     "justinfan16432.tmi.twitch.tv",
				},
				Command: CMDJoin,
				Params:  []string{"#jtv"},
			},
		},
		{
			":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv JOIN",
			&Message{
				Source: Source{
					Username: "justinfan16432",
					Nickname: "justinfan16432",
					Host:     "justinfan16432.tmi.twitch.tv",
				},
				Command: CMDJoin,
			},
		},
		{
			"@badge-info=;badges=moments/1;client-nonce=4fb782293442bb3b0df16b4cb5eb21aa;color=#008000;display-name=justinfan16432;emotes=;first-msg=0;flags=;id=6198df9e-77af-4f4f-8d3c-d317802b7c0d;mod=0;returning-chatter=0;room-id=14027;subscriber=0;tmi-sent-ts=1656612693901;turbo=0;user-id=16933;user-type= :justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv PRIVMSG #jtv :Hello",
			&Message{
				Tags: Tags{
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
				Source: Source{
					Username: "justinfan16432",
					Nickname: "justinfan16432",
					Host:     "justinfan16432.tmi.twitch.tv",
				},
				Command: CMDPrivMessage,
				Params:  []string{"#jtv"},
				Text:    "Hello",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			test.expected.Raw = test.raw
			message, err := ParseMessage(test.raw)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.expected, message)
		})
	}
}

func TestParseMessageError(t *testing.T) {
	tests := []struct {
		raw      string
		expected error
	}{
		{"@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0", ErrPartialMessage},
		{"@", ErrInvalidTags},
		{":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv", ErrNoCommand},
		{": JOIN #jtv", ErrInvalidSource},
	}

	for _, test := range tests {
		t.Run(test.expected.Error(), func(t *testing.T) {
			_, err := ParseMessage(test.raw)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		raw      string
		expected Tags
	}{
		{
			"@emote-only=0;followers-only=1440;r9k=0;room-id=14027;slow=0;subs-only=0",
			Tags{
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
			Tags{
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
			Tags{
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
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			tags, err := ParseTags(test.raw)
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
		{"", ErrInvalidTags},
	}

	for _, test := range tests {
		t.Run(test.expected.Error(), func(t *testing.T) {
			_, err := ParseTags(test.raw)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}

func TestParseSource(t *testing.T) {
	tests := []struct {
		raw      string
		expected *Source
	}{
		{
			":tmi.twitch.tv",
			&Source{
				Host: "tmi.twitch.tv",
			},
		},
		{
			":justinfan16432!justinfan16432.tmi.twitch.tv",
			&Source{
				Nickname: "justinfan16432",
				Username: "justinfan16432",
				Host:     "justinfan16432.tmi.twitch.tv",
			},
		},
		{
			":justinfan16432@justinfan16432.tmi.twitch.tv",
			&Source{
				Nickname: "justinfan16432",
				Username: "justinfan16432",
				Host:     "justinfan16432.tmi.twitch.tv",
			},
		},
		{
			":justinfan16432!justinfan16432@justinfan16432.tmi.twitch.tv",
			&Source{
				Nickname: "justinfan16432",
				Username: "justinfan16432",
				Host:     "justinfan16432.tmi.twitch.tv",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			source, err := ParseSource(test.raw)
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
		{"", ErrInvalidSource},
	}

	for _, test := range tests {
		t.Run(test.expected.Error(), func(t *testing.T) {
			_, err := ParseSource(test.raw)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}
