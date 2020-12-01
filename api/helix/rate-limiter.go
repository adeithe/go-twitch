package helix

import (
	"strconv"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/api/request"
)

// RateLimiter stores data used to keep track of the Twitch API Ratelimit.
type RateLimiter struct {
	Bucket    int
	Remaining int
	Reset     time.Time

	await chan bool
	open  bool
	close bool
	mx    sync.Mutex
}

// IRateLimiter contains all methods available to the RateLimiter.
type IRateLimiter interface {
	Enqueue(*request.HTTPRequest) (request.HTTPResponse, error)
	Close()
}

var _ IRateLimiter = &RateLimiter{}

// Enqueue queues a HTTP request for when the Twitch API will allow it to go through.
func (limiter *RateLimiter) Enqueue(req *request.HTTPRequest) (request.HTTPResponse, error) {
	limiter.mx.Lock()
	defer limiter.mx.Unlock()
	if !limiter.open {
		limiter.open = true
		limiter.await = make(chan bool)
		go limiter.ticker()
	}
	<-limiter.await
	limiter.Remaining--
	res, err := req.Do()
	if err != nil {
		limiter.update(req.Headers)
	}
	return res, err
}

// Close frees resources used by the RateLimiter.
func (limiter *RateLimiter) Close() {
	limiter.close = true
}

func (limiter *RateLimiter) update(headers map[string]string) {
	if bucket, err := strconv.Atoi(headers["Ratelimit-Limit"]); err == nil && limiter.Bucket != bucket {
		limiter.Bucket = bucket
	}
	if remaining, err := strconv.Atoi(headers["Ratelimit-Remaining"]); err == nil && limiter.Remaining != remaining {
		limiter.Remaining = remaining
	}
	if reset, err := strconv.ParseInt(headers["Ratelimit-Reset"], 10, 64); err == nil {
		limiter.Reset = time.Unix(reset, 0)
	}
}

func (limiter *RateLimiter) ticker() {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		if limiter.close {
			close(limiter.await)
			limiter.open = false
			limiter.close = false
			break
		}
		if limiter.Remaining <= 0 && limiter.Reset.Sub(time.Now().UTC()) > 0 {
			time.Sleep(time.Second)
			continue
		}
		timer.Reset(time.Second)
		select {
		case limiter.await <- true:
		case <-timer.C:
		}
	}
}
