package nonce

import (
	"math"
)

var s1, s2, s3 int = 100, 100, 100
var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// New nonce generated randomly using the Wichmann-Hill algorithm.
// https://en.wikipedia.org/wiki/Wichmann%E2%80%93Hill
func New() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = chars[int(math.Floor(random()*float64(len(chars))))]
	}
	return string(b)
}

func random() float64 {
	s1 = (171 * s1) % 30269
	s2 = (172 * s2) % 30307
	s3 = (170 * s3) % 30323
	val := float64(s1) / float64(30269.0)
	val += float64(s2) / float64(30307.0)
	val += float64(s3) / float64(30323.0)
	return math.Mod(val, float64(1.0))
}
