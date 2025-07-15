package core

import (
	"collyDemo/handlers"
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

// TaskConfig 任务配置结构
type TaskConfig struct {
	URL     string
	Method  string
	Body    []byte
	Handler func(*colly.Response, *Account, *TaskDispatcher) error
	Meta    map[string]interface{}
}

// GetDefaultHeaders 获取默认请求头
func GetDefaultHeaders(token string) map[string]string {
	return map[string]string{
		"accept":             "*/*",
		"accept-language":    "zh-HK,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
		"authorization":      token,
		"origin":             "https://www.kaogujia.com",
		"priority":           "u=1, i",
		"referer":            "https://www.kaogujia.com/",
		"sec-ch-ua":          `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
		"version_code":       "3.1",
		"content-type":       "application/json",
	}
}

// GetMainTasks 获取主要任务配置
func GetMainTasks() []TaskConfig {
	return []TaskConfig{
		// 达人任务
		{
			URL:     "https://service.kaogujia.com/api/author/search?limit=50&page=1&sort_field=gmv&sort=0",
			Method:  "POST",
			Body:    []byte(`{"sort_field":"gmv","sort":0,"limit":50,"page":1}`),
			Handler: handlers.AuthorHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 品牌任务
		{
			URL:     "https://service.kaogujia.com/api/brand/search?limit=50&page=1&sort_field=gmv&sort=0",
			Method:  "POST",
			Body:    []byte(`{"period":1,"keyword":""}`),
			Handler: handlers.BrandHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 直播任务
		{
			URL:     "https://service.kaogujia.com/api/live/search?limit=50&page=1&sort_field=gmv&sort=0",
			Method:  "POST",
			Body:    []byte(`{"pub_time":{"min":"20250629","max":"20250705"},"keyword":"","keyword_type":1}`),
			Handler: handlers.LiveHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 商品任务
		{
			URL:     "https://service.kaogujia.com/api/sku/search?limit=50&page=1&sort_field=sales&sort=0",
			Method:  "POST",
			Body:    []byte(`{"period":1,"keyword":""}`),
			Handler: handlers.ProductHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 店铺任务
		{
			URL:     "https://service.kaogujia.com/api/shop/search?limit=50&page=1&sort_field=gmv&sort=0",
			Method:  "POST",
			Body:    []byte(`{"period":1,"keyword":""}`),
			Handler: handlers.StoreHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 视频任务
		{
			URL:     "https://service.kaogujia.com/api/video/search?limit=50&page=1&sort_field=like_count&sort=0",
			Method:  "POST",
			Body:    []byte(`{"date_code":{"min":"20250629","max":"20250705"},"keyword":"","video_type":1}`),
			Handler: handlers.VideoHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
	}
}

// GetRankTasks 获取排名任务配置
func GetRankTasks() []TaskConfig {
	return []TaskConfig{
		// 达人涨粉榜
		{
			URL:     "https://service.kaogujia.com/api/rank/author/fans/increase?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.AuthorFansIncreaseRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 达人掉粉榜
		{
			URL:     "https://service.kaogujia.com/api/rank/author/fans/decrease?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.AuthorFansDecreaseRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 达人带货潜力榜
		{
			URL:     "https://service.kaogujia.com/api/rank/author/potential?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.AuthorPotentialRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 商品热销榜
		{
			URL:     "https://service.kaogujia.com/api/rank/product/hot/sale?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.ProductHotSaleRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 商品实时销量榜
		{
			URL:     "https://service.kaogujia.com/api/rank/product/real/time/sales?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.ProductRealTimeSalesRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 直播达人带货榜
		{
			URL:     "https://service.kaogujia.com/api/rank/live/author/sales?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.LiveAuthorSalesRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 直播热推榜
		{
			URL:     "https://service.kaogujia.com/api/rank/live/hot/push?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.LiveHotPushRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 热门视频榜
		{
			URL:     "https://service.kaogujia.com/api/rank/video/hot?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.HotVideoRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 电商视频榜
		{
			URL:     "https://service.kaogujia.com/api/rank/video/ecommerce?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.EcommerceVideoRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 视频热推
		{
			URL:     "https://service.kaogujia.com/api/rank/video/hot/push?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.VideoHotPushHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 热销小店
		{
			URL:     "https://service.kaogujia.com/api/rank/shop/hot/sale?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.HotSaleShopHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 全站小时榜
		{
			URL:     "https://service.kaogujia.com/api/rank/site/hourly?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.SiteHourlyRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 带货小时榜
		{
			URL:     "https://service.kaogujia.com/api/rank/sales/hourly?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.SalesHourlyRankHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 实时热点
		{
			URL:     "https://service.kaogujia.com/api/hot/spot/real/time?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.RealTimeHotSpotHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 飙升热点
		{
			URL:     "https://service.kaogujia.com/api/hot/spot/soaring?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.SoaringHotSpotHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
		// 探测爆款
		{
			URL:     "https://service.kaogujia.com/api/explore/hot/burst?limit=50&page=1",
			Method:  "GET",
			Handler: handlers.ExploreHotBurstHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		},
	}
}

// CreatePaginationTask 创建分页任务
func CreatePaginationTask(baseURL string, method string, body []byte, handler func(*colly.Response, *Account, *TaskDispatcher) error, page, limit int64, headers map[string]string) *Task {
	url := fmt.Sprintf("%s&page=%d&limit=%d", baseURL, page, limit)
	return &Task{
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
		Handler: handler,
		Meta: map[string]interface{}{
			"page":  page,
			"limit": limit,
		},
	}
}

// CreateDetailTask 创建详情任务
func CreateDetailTask(baseURL, id string, handler func(*colly.Response, *Account, *TaskDispatcher) error, headers map[string]string, meta map[string]interface{}) *Task {
	url := fmt.Sprintf("%s/%s", baseURL, id)
	return &Task{
		URL:     url,
		Method:  "GET",
		Headers: headers,
		Handler: handler,
		Meta:    meta,
	}
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	dispatcher *TaskDispatcher
	headers    map[string]string
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(dispatcher *TaskDispatcher, token string) *TaskScheduler {
	return &TaskScheduler{
		dispatcher: dispatcher,
		headers:    GetDefaultHeaders(token),
	}
}

// AddMainTasks 添加主要任务
func (s *TaskScheduler) AddMainTasks() {
	tasks := GetMainTasks()
	for _, config := range tasks {
		task := &Task{
			URL:     config.URL,
			Method:  config.Method,
			Headers: s.headers,
			Body:    config.Body,
			Handler: config.Handler,
			Meta:    config.Meta,
		}
		s.dispatcher.AddTask(task)
	}
}

// AddRankTasks 添加排名任务
func (s *TaskScheduler) AddRankTasks() {
	tasks := GetRankTasks()
	for _, config := range tasks {
		task := &Task{
			URL:     config.URL,
			Method:  config.Method,
			Headers: s.headers,
			Body:    config.Body,
			Handler: config.Handler,
			Meta:    config.Meta,
		}
		s.dispatcher.AddTask(task)
	}
}

// SchedulePeriodicTasks 调度周期性任务
func (s *TaskScheduler) SchedulePeriodicTasks() {
	// 每小时执行一次排名任务
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.AddRankTasks()
			}
		}
	}()
}
