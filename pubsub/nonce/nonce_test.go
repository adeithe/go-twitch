package nonce_test

import (
	"testing"

	"github.com/adeithe/go-twitch/pubsub/nonce"
	"github.com/stretchr/testify/require"
)

func TestNonce_WichmannHill(t *testing.T) {
	expected := []string{
		"QGj7ZX885miaWXziC2VqYnUlNQhpcUeQ",
		"szFJVHGmyTRVuXdrZN1uT9Slhjx98SHq",
		"ELsgOfsOrjYYA5XsSz0XPckLCJd2RMqR",
		"nJB2KriBHVcfZocLJnIf4iVs3KGF7byE",
		"Uq1qCVDeLT9UnPULckKsziLbiFsYXcQl",
		"QQLGpQSC2w505zD9nBp1tlKi3H6IULUH",
		"s2bolhdH9nyHtREMJTTNWSJcW5EhvWuk",
		"sso3wlzM5BWZ7jCLe7LwD9bCxMyj0II4",
		"IfaIX6sOHxOu7MPp1SkpET2O23FsiW65",
		"qfVc68AhIucS2W3ptlC5OVxpoPxXdwIj",
	}

	for i := range 10 {
		require.Equal(t, expected[i], nonce.WichmannHill())
	}
}

func BenchmarkNonce_WichmannHill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nonce.WichmannHill()
	}
}
