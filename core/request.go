package core

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

var (
	dispatcher     *TaskDispatcher
	dispatcherOnce sync.Once
)

func ExecuteRequest(task *Task, account *Account, dispatcher *TaskDispatcher) error {
	// 在执行请求前应用速率限制
	if account.RateLimit != nil {
		account.RateLimit.Wait()
	}

	c := colly.NewCollector(
		colly.Async(true),
	)

	// 设置超时
	c.SetRequestTimeout(30 * time.Second)

	// 设置代理
	if account.Proxy != "" {
		if err := c.SetProxy(account.Proxy); err != nil {
			log.Printf("设置代理失败: %v", err)
			return err
		}
	}

	// 创建请求
	var body io.Reader
	if task.Body != nil {
		body = bytes.NewReader(task.Body)
	}

	request, err := http.NewRequest(task.Method, task.URL, body)
	if err != nil {
		return err
	}

	// 设置请求头
	for k, v := range task.Headers {
		if k == "authorization" {
			v = account.Token
		}
		request.Header.Set(k, v)
	}

	// 用于等待响应的 channel
	done := make(chan error, 1)

	// 注册响应处理
	c.OnResponse(func(r *colly.Response) {
		if err := task.Handler(r, account, dispatcher); err != nil {
			log.Printf("处理器错误: %v", err)
			done <- err
		} else {
			done <- nil
		}
	})

	// 处理请求错误
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("请求失败: %s, 错误: %v, 账号: %s", task.URL, err, account.UserName)
		done <- err
	})

	// 发送请求
	err = c.Request(request.Method, request.URL.String(), request.Body, nil, request.Header)
	if err != nil {
		return err
	}

	// 等待响应处理完成，带超时机制
	select {
	case err := <-done:
		return err
	case <-time.After(60 * time.Second):
		log.Printf("请求超时: %s", task.URL)
		return context.DeadlineExceeded
	}
}
