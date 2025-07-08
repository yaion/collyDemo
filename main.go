package main

import (
	"collyDemo/core"
	"collyDemo/handlers"
	"collyDemo/mongodb"
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
		/*{
			ID:        "1",
			UserName:  "19282586479",
			Password:  "ee882244",
			Proxy:     "http://127.0.0.1:7890", // 替换为实际代理地址
			Token:     "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiZTlkMWEzYWUzNWI2NDA4ZDk5ODEyZDMzZjE3MjMyNGQiLCJzaWQiOjgzMDU3NTQsImlhdCI6MTc1MTc4ODkwMSwiZXhwIjoxNzUyMzkzNzAxLCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.ypEQ64GQgUspSGqr_RuuNVYrIztdpSt-m7Qc-Dw_mCigk_3At7hwFOf_XGAxDS82-sK1kLbz78vxobX5Snii3w",
			LastUsed:  time.Now().Add(-5 * time.Minute),
			RateLimit: core.NewRateLimiter(60),
			MinDelay:  2 * time.Second,
			MaxDelay:  5 * time.Second,
		},*/
		{
			ID:       "2",
			UserName: "19292545479",
			Password: "520134lwn",
			//Proxy:     "http://114.80.35.137:3081", // 替换为实际代理地址
			Token:     "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiNzgyMjQ0NDM2MWJjNGJjN2JlZjk0MTA5ZTIyNDRlZjQiLCJzaWQiOjgzMDU3NTYsImlhdCI6MTc1MTkwMTE3MSwiZXhwIjoxNzUyNTA1OTcxLCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.SEi6TkurfojUhtopsR7JMNN-VIeSz-0a0pOOiwtTda-qvYb0HyEHFs64frOl3gqlXkIrA9yWaZXQ2rYu8aLnXA",
			LastUsed:  time.Now().Add(-5 * time.Minute),
			RateLimit: core.NewRateLimiter(60),
			MinDelay:  2 * time.Second,
			MaxDelay:  5 * time.Second,
		},
		{
			ID:       "3",
			UserName: "19292481905",
			Password: "aa123619",
			//Proxy:     "http://60.188.102.225:18080", // 替换为实际代理地址
			Token:     "Bearer eyJhbGciOiJIUzUxMiJ9.eyJhdWQiOiIxMDAwIiwiaXNzIjoia2FvZ3VqaWEuY29tIiwianRpIjoiOTQ4NGIyMmFmODdkNDc2NzkxOWY3YzA1MzhhZjFjMjQiLCJzaWQiOjgzMDU3NTQsImlhdCI6MTc1MTkwMTEwNSwiZXhwIjoxNzUyNTA1OTA1LCJid2UiOjAsInR5cCI6MSwicF9id2UiOjB9.TFsB60-e1M2BgBWivwXb0ZmfTs-1nuT7k2XZLsERP7MGTxis_Z4wCBHlpkc89aPi6bYxaiAL3Wrg-MasRUnvpg",
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

	/*// 添加初始author任务
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
	})*/

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
	/*
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
		})*/
	// 启动爬虫，设置并发数为3
	dispatcher.Run(len(accounts))

	log.Println("所有任务处理完成")
}
