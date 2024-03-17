package tools

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrRateLimited = errors.New("rate limited")
)

func NewRateLimiter(seconds int64) *RateLimiter {
	rl := &RateLimiter{
		wg:      sync.WaitGroup{},
		done:    make(chan bool),
		seconds: seconds,
		keys:    make(map[string]time.Time),
	}

	go rl.run()

	return rl
}

type RateLimiter struct {
	wg   sync.WaitGroup
	done chan bool

	seconds int64
	keys    map[string]time.Time
}

func (rl *RateLimiter) Set(key string) {
	rl.operation(func(l *RateLimiter) {
		l.keys[key] = time.Now()
	})
}

func (rl *RateLimiter) Check(key string) (err error) {
	rl.operation(func(l *RateLimiter) {
		_, ok := l.keys[key]
		if ok {
			err = ErrRateLimited
		}
	})

	return
}

func (rl *RateLimiter) Deinit() {
	rl.done <- true
}

func (rl *RateLimiter) clear(t time.Time) {
	rl.operation(func(l *RateLimiter) {
		keys := []string{}

		for key, time := range l.keys {
			if time.Unix()+l.seconds <= t.Unix() {
				keys = append(keys, key)
			}
		}

		for _, key := range keys {
			delete(l.keys, key)
		}
	})
}

func (rl *RateLimiter) operation(op func(*RateLimiter)) {
	rl.wg.Wait()
	rl.wg.Add(1)
	op(rl)
	rl.wg.Done()
}

func (rl *RateLimiter) run() {
	ticker := time.NewTicker(time.Second * time.Duration(rl.seconds))
	exit := false

	for !exit {
		select {
		case t := <-ticker.C:
			rl.clear(t)
		case <-rl.done:
			exit = true
		}
	}
}
