package core

import (
	"github.com/gocolly/colly/v2"
	"sync"
	"time"
)

// 账号结构
type Account struct {
	ID        string
	UserName  string
	Password  string
	Proxy     string
	Token     string
	LastUsed  time.Time
	RateLimit *RateLimiter
	MinDelay  time.Duration // 最小延迟(2秒)
	MaxDelay  time.Duration // 最大延迟(3秒)
	mu        sync.Mutex
}

/*// 分页响应结构
type PagedResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// API1 响应数据结构
type API1Response struct {
	PagedResponse
	Items []API1Item `json:"items"`
}

type API1Item struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

// API2 响应数据结构
type API2Response struct {
	ID       string                 `json:"id"`
	Details  string                 `json:"details"`
	Metadata map[string]interface{} `json:"metadata"`
}*/

// 任务结构
type Task struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
	Handler func(*colly.Response, *Account, *TaskDispatcher) error
	Meta    map[string]interface{}
}
