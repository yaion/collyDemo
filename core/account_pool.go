package core

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type AccountPool struct {
	accounts   []*Account
	interval   time.Duration
	currentIdx int
	mu         sync.Mutex
}

func NewAccountPool(accounts []*Account, interval time.Duration) *AccountPool {
	return &AccountPool{
		accounts: accounts,
		interval: interval,
	}
}

func (p *AccountPool) GetAccount() *Account {
	p.mu.Lock()
	defer p.mu.Unlock()

	startIdx := p.currentIdx

	for {
		p.currentIdx = (p.currentIdx + 1) % len(p.accounts)
		acc := p.accounts[p.currentIdx]

		acc.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(acc.LastUsed)

		// 计算所需延迟时间
		requiredDelay := acc.MinDelay + time.Duration(rand.Int63n(int64(acc.MaxDelay-acc.MinDelay)))

		if elapsed >= requiredDelay {
			acc.LastUsed = now
			acc.mu.Unlock()
			log.Printf("获取账号成功: %s, 延迟: %v", acc.UserName, elapsed)
			return acc
		}
		acc.mu.Unlock()

		// 如果转了一圈都没找到可用账号
		if p.currentIdx == startIdx {
			// 找出最早可用的账号和等待时间
			var minWaitTime time.Duration = time.Hour
			for _, acc := range p.accounts {
				acc.mu.Lock()
				elapsed := time.Since(acc.LastUsed)
				requiredDelay := acc.MinDelay + time.Duration(rand.Int63n(int64(acc.MaxDelay-acc.MinDelay)))
				waitTime := requiredDelay - elapsed
				if waitTime > 0 && waitTime < minWaitTime {
					minWaitTime = waitTime
				}
				acc.mu.Unlock()
			}

			if minWaitTime > 0 && minWaitTime < time.Hour {
				log.Printf("所有账号都在冷却中，等待 %v", minWaitTime)
				// 解锁后再等待，避免长时间持有锁
				p.mu.Unlock()
				time.Sleep(minWaitTime)
				p.mu.Lock()
				// 等待后重新开始循环
				continue
			}
			break
		}
	}

	// 如果还是没有可用账号，递归重试
	return p.GetAccount()
}
