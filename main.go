package main

import (
	"collyDemo/config"
	"collyDemo/core"
	"collyDemo/handlers"
	"collyDemo/mongodb"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化 mongo
	mongodb.InitMongo()
	rand.Seed(time.Now().UnixNano())

	// 加载配置
	scheduleConfig := config.GetDefaultConfig()

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

	// 创建任务配置调度器
	taskScheduler := core.NewTaskScheduler(dispatcher, accounts[0].Token)

	// 注册所有处理器
	registerHandlers(taskScheduler)

	// 创建定时任务调度器
	scheduler := core.NewScheduler(dispatcher, taskScheduler)

	// 使用配置文件初始化定时任务
	mainSchedules := scheduleConfig.GetMainTaskSchedules()
	rankSchedules := scheduleConfig.GetRankTaskSchedules()
	scheduler.InitTasksWithConfig(mainSchedules, rankSchedules)

	// 启动定时任务调度器
	scheduler.Start()

	// 启动爬虫工作池，设置并发数为账号数量
	go dispatcher.Run(scheduleConfig.System.MaxConcurrency)

	// 启动任务状态监控
	go monitorTaskStatus(dispatcher, scheduler)

	// 等待中断信号
	waitForInterrupt()

	// 优雅关闭
	log.Println("正在关闭系统...")
	scheduler.Stop()
	dispatcher.Stop()
	log.Println("系统已关闭")
}

// registerHandlers 注册所有处理器
func registerHandlers(taskScheduler *core.TaskScheduler) {
	// 主要数据处理器
	taskScheduler.RegisterHandler("author", handlers.AuthorHandler)
	taskScheduler.RegisterHandler("brand", handlers.BrandHandler)
	taskScheduler.RegisterHandler("live", handlers.LiveHandler)
	taskScheduler.RegisterHandler("product", handlers.ProductHandler)
	taskScheduler.RegisterHandler("store", handlers.StoreHandler)
	taskScheduler.RegisterHandler("video", handlers.VideoHandler)

	// 详情处理器
	taskScheduler.RegisterHandler("author_info", handlers.AuthorInfoHandler)
	taskScheduler.RegisterHandler("brand_info", handlers.BrandInfoHandler)
	taskScheduler.RegisterHandler("product_info", handlers.ProductInfoHandler)
	taskScheduler.RegisterHandler("live_info", handlers.LiveInfoHandler)
	taskScheduler.RegisterHandler("video_info", handlers.VideoInfoHandler)
	taskScheduler.RegisterHandler("store_info", handlers.StoreInfoHandler)

	// 排名数据处理器
	taskScheduler.RegisterHandler("author_fans_increase_rank", handlers.AuthorFansIncreaseRankHandler)
	taskScheduler.RegisterHandler("author_fans_decrease_rank", handlers.AuthorFansDecreaseRankHandler)
	taskScheduler.RegisterHandler("author_potential_rank", handlers.AuthorPotentialRankHandler)
	taskScheduler.RegisterHandler("product_hot_sale_rank", handlers.ProductHotSaleRankHandler)
	taskScheduler.RegisterHandler("product_real_time_sales_rank", handlers.ProductRealTimeSalesRankHandler)
	taskScheduler.RegisterHandler("live_author_sales_rank", handlers.LiveAuthorSalesRankHandler)
	taskScheduler.RegisterHandler("live_hot_push_rank", handlers.LiveHotPushRankHandler)
	taskScheduler.RegisterHandler("hot_video_rank", handlers.HotVideoRankHandler)
	taskScheduler.RegisterHandler("ecommerce_video_rank", handlers.EcommerceVideoRankHandler)
	taskScheduler.RegisterHandler("video_hot_push", handlers.VideoHotPushHandler)
	taskScheduler.RegisterHandler("hot_sale_shop", handlers.HotSaleShopHandler)
	taskScheduler.RegisterHandler("site_hourly_rank", handlers.SiteHourlyRankHandler)
	taskScheduler.RegisterHandler("sales_hourly_rank", handlers.SalesHourlyRankHandler)
	taskScheduler.RegisterHandler("real_time_hot_spot", handlers.RealTimeHotSpotHandler)
	taskScheduler.RegisterHandler("soaring_hot_spot", handlers.SoaringHotSpotHandler)
	taskScheduler.RegisterHandler("explore_hot_burst", handlers.ExploreHotBurstHandler)

	log.Println("所有处理器注册完成")
}

// monitorTaskStatus 监控任务状态
func monitorTaskStatus(dispatcher *core.TaskDispatcher, scheduler *core.Scheduler) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			queueLen, active := dispatcher.TaskStatus()
			taskStatus := scheduler.GetTaskStatus()

			log.Printf("=== 系统状态监控 ===")
			log.Printf("任务队列长度: %d", queueLen)
			log.Printf("活跃任务数: %d", active)
			log.Printf("定时任务状态:")

			for id, status := range taskStatus {
				statusMap := status.(map[string]interface{})
				log.Printf("  %s: %s (下次执行: %s)",
					id,
					statusMap["name"],
					statusMap["next_run"].(time.Time).Format("2006-01-02 15:04:05"))
			}
			log.Printf("==================")
		}
	}
}

// waitForInterrupt 等待中断信号
func waitForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
