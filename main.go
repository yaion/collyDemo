package main

import (
	"collyDemo/core"
	"collyDemo/handlers"
	"collyDemo/mongodb"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	// 初始化 mongo
	mongodb.InitMongo()
	rand.Seed(time.Now().UnixNano())

	// 初始化账号池
	accounts := []*core.Account{
		{
			ID:        "1",
			UserName:  "19282586479",
			Password:  "ee882244",
			Proxy:     "http://127.0.0.1:1080", // 替换为实际代理地址
			Token:     "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiZTlkMWEzYWUzNWI2NDA4ZDk5ODEyZDMzZjE3MjMyNGQiLCJzaWQiOjgzMDU3NTQsImlhdCI6MTc1MTc4ODkwMSwiZXhwIjoxNzUyMzkzNzAxLCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.ypEQ64GQgUspSGqr_RuuNVYrIztdpSt-m7Qc-Dw_mCigk_3At7hwFOf_XGAxDS82-sK1kLbz78vxobX5Snii3w",
			LastUsed:  time.Now().Add(-5 * time.Minute),
			RateLimit: core.NewRateLimiter(60),
			MinDelay:  2 * time.Second,
			MaxDelay:  5 * time.Second,
		},
		{
			ID:        "2",
			UserName:  "19292545479",
			Password:  "520134lwn",
			Proxy:     "http://127.0.0.1:1080", // 替换为实际代理地址
			Token:     "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiOGZjYzE3YmNjMjlkNGUyNGE1Yjc3NDY0NTgyNzJkNTgiLCJzaWQiOjgzMDU3NTYsImlhdCI6MTc1MTc4ODk0MCwiZXhwIjoxNzUyMzkzNzQwLCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.zzG5mr1AtQolGwo7CvngDeQIdmRSaqbmCG5Q-rm84EY7qqoj68sz_BxPrcqOX83woefMe28JD56Ud_KOIb7d_g",
			LastUsed:  time.Now().Add(-5 * time.Minute),
			RateLimit: core.NewRateLimiter(60),
			MinDelay:  2 * time.Second,
			MaxDelay:  5 * time.Second,
		},
		{
			ID:        "3",
			UserName:  "19292481905",
			Password:  "aa123619",
			Proxy:     "http://127.0.0.1:1080", // 替换为实际代理地址
			Token:     "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiMjEwOTNmOGQ1MjA2NDRjNDkwMjY3OGI1ZjlmM2EzNzAiLCJzaWQiOjgyOTk1MDcsImlhdCI6MTc1MTc4OTI3MSwiZXhwIjoxNzUyMzk0MDcxLCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.nFMTaVo1wRLBeLmEPdBoY7IsoAS9l_Cc4dsdqGgQpxrMFzwS1dyTtYTbmBkTNQiaUIBkvfHi5xgRGVD9RFsEJg",
			LastUsed:  time.Now().Add(-5 * time.Minute),
			RateLimit: core.NewRateLimiter(60),
			MinDelay:  2 * time.Second,
			MaxDelay:  5 * time.Second,
		},
	}
	accountPool := core.NewAccountPool(accounts, 3*time.Second)

	// 创建任务调度器
	dispatcher := core.NewTaskDispatcher(accountPool)

	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "zh-HK,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
		"authorization":      "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiNjI3MmYyN2EyZDU5NDc0YThhYzk1NTQyNzgyYjM4OWIiLCJzaWQiOjgyMjY4NzcsImlhdCI6MTc0OTkwMDQ0MywiZXhwIjoxNzUwNTA1MjQzLCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.sTul3qBenukj-HiOsTS_CnzHM0TV91cLA_U6dm6U5Z5ZFYgu6ZeTM3_Ai4AYdmvDN7q_SMoFjoQvv_LNo2VdzQ",
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

	// 添加初始author任务
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/author/search?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"sort_field":"gmv","sort":0,"limit":50,"page":1}`),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//直播达人带货榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/live/sku/author?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"sort_field":"gmv","sort":0,"limit":50,"page":1}`),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//直播达人带货榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/live/sku/author?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1,"m_lv1":0}`, time.Now().Format("20060102"))), //{"date_code":20250712,"period":1,"m_lv1":0}
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//带货达人潜力榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/author/potential?limit=50&page=1&sort_field=potential&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":30}`, time.Now().Format("20060102"))),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 达人涨粉榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/author/fans/increment?limit=50&page=1&sort_field=inc_fans_count&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"type":1,"date_code":%v,"period":1,"lv1":0}`, time.Now().Format("20060102"))),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 达人掉粉榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/author/fans/increment?limit=50&page=1&sort_field=inc_fans_count&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"type":2,"date_code":%v,"period":1,"lv1":0}`, time.Now().Format("20060102"))),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加初始品牌任务
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/brand/search?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"period":1,"keyword":""}`),
		Handler: handlers.BrandHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})
	//品牌热销榜
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/brand?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1,"m_lv1":0}`, time.Now().Format("20060102"))),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//品牌声量榜
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/brandsov?limit=50&page=1&sort_field=expose_count&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1,"m_lv1":0}`, time.Now().Format("20060102"))),
		Handler: handlers.AuthorHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加初始直播任务
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/live/search?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"pub_time":{"min":"20250629","max":"20250705"},"keyword":"","keyword_type":1}`),
		Handler: handlers.LiveHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 带货小时榜 每小时一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/official/live/sku/hour?limit=50&page=1&sort_field=score&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"hh":%v,"lv1":0,"m_lv1":0}`, time.Now().Format("20060102"), time.Now().Hour())),
		Handler: handlers.LiveHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 带货小时榜 每小时一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/official/live/hour?limit=50&page=1&sort_field=gap_description&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"hh":%v,"lv1":0,"m_lv1":0}`, time.Now().Format("20060102"), time.Now().Hour())),
		Handler: handlers.LiveHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加初始商品任务
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/sku/search?limit=50&page=1&sort_field=sales&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"period":1,"keyword":""}`),
		Handler: handlers.ProductHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加实时销量榜 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/sku/rta?limit=50&page=1&sort_field=h2_sales&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"m_lv1":0}`),
		Handler: handlers.ProductHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加商品热销榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/sku/pmt/2?limit=50&page=1&sort_field=sales&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1}`, time.Now().Format("20060102"))),
		Handler: handlers.ProductHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 直播热推榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/sku/rta?limit=50&page=1&sort_field=h2_sales&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1}`, time.Now().Format("20060102"))),
		Handler: handlers.ProductHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 视频热推榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/video/sku/2?limit=50&page=1&sort_field=sales&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1}`, time.Now().Format("20060102"))),
		Handler: handlers.ProductHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加初始店铺任务
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/shop/search?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"period":1,"keyword":""}`),
		Handler: handlers.StoreHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//热销小店榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/shop/hot?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1,"m_lv1":0}`, time.Now().Format("20060102"))),
		Handler: handlers.StoreHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//地区小店榜 一天一次 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/shop/area?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1,"city_id":0,"province_id":0}`, time.Now().Format("20060102"))),
		Handler: handlers.StoreHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 添加初始视频任务
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/video/search?limit=50&page=1&sort_field=like_count&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(`{"date_code":{"min":"20250629","max":"20250705"},"keyword":"","video_type":1}`),
		Handler: handlers.VideoHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//热门视频榜 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/video?limit=50&page=1&sort_field=like_count&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1,"hide_gov_media_vip":1,"lv1":0,"is_sku":0}`, time.Now().Format("20060102"))),
		Handler: handlers.VideoHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//电商视频榜 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/productvideo?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"date_code":%v,"period":1}`, time.Now().Format("20060102"))),
		Handler: handlers.VideoHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	//图文带货榜 todo
	dispatcher.AddTask(&core.Task{
		URL:     "https://service.kaogujia.com/api/rank/productvideo?limit=50&page=1&sort_field=gmv&sort=0",
		Method:  "POST",
		Headers: headers,
		Body:    []byte(fmt.Sprintf(`{"is_image":1,"date_code":%v,"period":1}`, time.Now().Format("20060102"))),
		Handler: handlers.VideoHandler,
		Meta: map[string]interface{}{
			"page":     1,
			"pageSize": 50,
		},
	})

	// 启动爬虫，设置并发数为3
	dispatcher.Run(len(accounts))

	log.Println("所有任务处理完成")
}
