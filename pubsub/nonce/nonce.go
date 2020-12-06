package nonce

import (
	"math"
)

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	whS1, whS2, whS3 int = 100, 100, 100
)

// WichmannHill nonce generated randomly using the Wichmann-Hill algorithm.
// https://en.wikipedia.org/wiki/Wichmann%E2%80%93Hill
func WichmannHill() string {
	random := func() float64 {
		whS1 = (171 * whS1) % 30269
		whS2 = (172 * whS2) % 30307
		whS3 = (170 * whS3) % 30323
		val := float64(whS1) / float64(30269.0)
		val += float64(whS2) / float64(30307.0)
		val += float64(whS3) / float64(30323.0)
		return math.Mod(val, float64(1.0))
	}
	b := make([]rune, 32)
	for i := range b {
		b[i] = chars[int(math.Floor(random()*float64(len(chars))))]
	}
	return string(b)
}
