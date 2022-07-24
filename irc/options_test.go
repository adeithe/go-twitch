package irc

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapabilities(t *testing.T) {
	tests := []struct {
		name     string
		opts     []ConnOption
		expected []Capability
	}{
		{
			name: "WithTags",
			opts: []ConnOption{
				WithTags(),
			},
			expected: []Capability{
				CapabilityCommands,
				CapabilityTags,
			},
		},
		{
			name: "WithCommands",
			opts: []ConnOption{
				WithCommands(),
			},
			expected: []Capability{
				CapabilityCommands,
			},
		},
		{
			name: "WithMembership",
			opts: []ConnOption{
				WithMembership(),
			},
			expected: []Capability{
				CapabilityMembership,
			},
		},
		{
			name: "WithCapability",
			opts: []ConnOption{
				WithCapability(CapabilityTags),
				WithCapability(CapabilityCommands),
				WithCapability(CapabilityMembership),
			},
			expected: []Capability{
				CapabilityTags,
				CapabilityCommands,
				CapabilityMembership,
			},
		},
		{
			name: "CapabilityDeduplication",
			opts: []ConnOption{
				WithTags(),
				WithCommands(),
				WithMembership(),
				WithCapability(CapabilityTags),
				WithCapability(CapabilityCommands),
				WithCapability(CapabilityMembership),
			},
			expected: []Capability{
				CapabilityCommands,
				CapabilityMembership,
				CapabilityTags,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(nil, test.opts...)
			assert.NotNil(t, client)
			assert.ElementsMatch(t, test.expected, client.Capabilities())
		})
	}
}

func TestWithAuth(t *testing.T) {
	tests := []struct {
		username, token string
	}{
		{"justinfan123", "Kappa123"},
		{"justinfan123", "oauth:Kappa123"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			client := New(nil, WithAuth(test.username, test.token))
			assert.NotNil(t, client)

			test.token = strings.TrimPrefix(test.token, "oauth:")
			assert.Equal(t, test.username, client.username)
			assert.Equal(t, test.token, client.token)
		})
	}
}

func TestWithAddr(t *testing.T) {
	expected, err := net.ResolveTCPAddr("tcp", "irc.chat.twitch.tv:6697")
	assert.NoError(t, err)

	client := New(nil, WithAddr(expected))
	if assert.NotNil(t, client) {
		assert.Equal(t, expected, client.Addr())
	}
}

func TestWithHostname(t *testing.T) {
	expected := "twitch.tv"
	client := New(nil, WithHostname(expected))
	if assert.NotNil(t, client) {
		assert.Equal(t, expected, client.hostname)
	}
}

func TestWithoutTLS(t *testing.T) {
	client := New(nil, WithoutTLS())
	if assert.NotNil(t, client) {
		assert.False(t, client.tls)
	}
}

func TestWithInsecure(t *testing.T) {
	client := New(nil, WithInsecure())
	if assert.NotNil(t, client) {
		assert.True(t, client.insecure)
	}
}
