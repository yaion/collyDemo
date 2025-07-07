package core

import (
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	limit      int
	tokens     int
	maxTokens  int
	lastRefill time.Time
	refillRate time.Duration
	mu         sync.Mutex
}

func NewRateLimiter(limitPerMinute int) *RateLimiter {
	refillRate := time.Minute / time.Duration(limitPerMinute)
	return &RateLimiter{
		limit:      limitPerMinute,
		tokens:     limitPerMinute,
		maxTokens:  limitPerMinute,
		lastRefill: time.Now(),
		refillRate: refillRate,
	}
}

func (rl *RateLimiter) Limit() int {
	return rl.limit
}

func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}

	// 如果令牌不足，等待
	if rl.tokens <= 0 {
		waitTime := rl.refillRate
		log.Printf("速率限制等待: %v", waitTime)
		time.Sleep(waitTime)
		rl.tokens = 1
		rl.lastRefill = time.Now()
	}

	// 消耗一个令牌
	rl.tokens--
}
