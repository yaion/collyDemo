package config

import (
	"time"
)

// ScheduleConfig 定时任务配置
type ScheduleConfig struct {
	// 主要数据采集任务配置
	MainTasks struct {
		Author  string `json:"author"`  // 达人数据采集频率
		Brand   string `json:"brand"`   // 品牌数据采集频率
		Live    string `json:"live"`    // 直播数据采集频率
		Product string `json:"product"` // 商品数据采集频率
		Store   string `json:"store"`   // 店铺数据采集频率
		Video   string `json:"video"`   // 视频数据采集频率
	} `json:"main_tasks"`

	// 排名数据采集任务配置
	RankTasks struct {
		AuthorFansIncreaseRank   string `json:"author_fans_increase_rank"`    // 达人涨粉榜
		AuthorFansDecreaseRank   string `json:"author_fans_decrease_rank"`    // 达人掉粉榜
		AuthorPotentialRank      string `json:"author_potential_rank"`        // 达人带货潜力榜
		ProductHotSaleRank       string `json:"product_hot_sale_rank"`        // 商品热销榜
		ProductRealTimeSalesRank string `json:"product_real_time_sales_rank"` // 商品实时销量榜
		LiveAuthorSalesRank      string `json:"live_author_sales_rank"`       // 直播达人带货榜
		LiveHotPushRank          string `json:"live_hot_push_rank"`           // 直播热推榜
		HotVideoRank             string `json:"hot_video_rank"`               // 热门视频榜
		EcommerceVideoRank       string `json:"ecommerce_video_rank"`         // 电商视频榜
		VideoHotPush             string `json:"video_hot_push"`               // 视频热推
		HotSaleShop              string `json:"hot_sale_shop"`                // 热销小店
		SiteHourlyRank           string `json:"site_hourly_rank"`             // 全站小时榜
		SalesHourlyRank          string `json:"sales_hourly_rank"`            // 带货小时榜
		RealTimeHotSpot          string `json:"real_time_hot_spot"`           // 实时热点
		SoaringHotSpot           string `json:"soaring_hot_spot"`             // 飙升热点
		ExploreHotBurst          string `json:"explore_hot_burst"`            // 探测爆款
	} `json:"rank_tasks"`

	// 系统配置
	System struct {
		MaxRetries     int           `json:"max_retries"`     // 最大重试次数
		RetryDelay     time.Duration `json:"retry_delay"`     // 重试延迟
		TaskTimeout    time.Duration `json:"task_timeout"`    // 任务超时时间
		MaxConcurrency int           `json:"max_concurrency"` // 最大并发数
	} `json:"system"`
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *ScheduleConfig {
	config := &ScheduleConfig{}

	// 主要数据采集任务默认配置
	config.MainTasks.Author = "6h"  // 每6小时采集一次达人数据
	config.MainTasks.Brand = "12h"  // 每12小时采集一次品牌数据
	config.MainTasks.Live = "2h"    // 每2小时采集一次直播数据
	config.MainTasks.Product = "4h" // 每4小时采集一次商品数据
	config.MainTasks.Store = "8h"   // 每8小时采集一次店铺数据
	config.MainTasks.Video = "3h"   // 每3小时采集一次视频数据

	// 排名数据采集任务默认配置
	config.RankTasks.AuthorFansIncreaseRank = "1h"   // 每小时采集一次
	config.RankTasks.AuthorFansDecreaseRank = "1h"   // 每小时采集一次
	config.RankTasks.AuthorPotentialRank = "1h"      // 每小时采集一次
	config.RankTasks.ProductHotSaleRank = "1h"       // 每小时采集一次
	config.RankTasks.ProductRealTimeSalesRank = "1h" // 每小时采集一次
	config.RankTasks.LiveAuthorSalesRank = "1h"      // 每小时采集一次
	config.RankTasks.LiveHotPushRank = "1h"          // 每小时采集一次
	config.RankTasks.HotVideoRank = "1h"             // 每小时采集一次
	config.RankTasks.EcommerceVideoRank = "1h"       // 每小时采集一次
	config.RankTasks.VideoHotPush = "1h"             // 每小时采集一次
	config.RankTasks.HotSaleShop = "1h"              // 每小时采集一次
	config.RankTasks.SiteHourlyRank = "1h"           // 每小时采集一次
	config.RankTasks.SalesHourlyRank = "1h"          // 每小时采集一次
	config.RankTasks.RealTimeHotSpot = "1h"          // 每小时采集一次
	config.RankTasks.SoaringHotSpot = "1h"           // 每小时采集一次
	config.RankTasks.ExploreHotBurst = "1h"          // 每小时采集一次

	// 系统默认配置
	config.System.MaxRetries = 3
	config.System.RetryDelay = 5 * time.Second
	config.System.TaskTimeout = 5 * time.Minute
	config.System.MaxConcurrency = 3

	return config
}

// GetMainTaskSchedules 获取主要任务的调度配置
func (c *ScheduleConfig) GetMainTaskSchedules() map[string]string {
	return map[string]string{
		"author":  c.MainTasks.Author,
		"brand":   c.MainTasks.Brand,
		"live":    c.MainTasks.Live,
		"product": c.MainTasks.Product,
		"store":   c.MainTasks.Store,
		"video":   c.MainTasks.Video,
	}
}

// GetRankTaskSchedules 获取排名任务的调度配置
func (c *ScheduleConfig) GetRankTaskSchedules() map[string]string {
	return map[string]string{
		"author_fans_increase_rank":    c.RankTasks.AuthorFansIncreaseRank,
		"author_fans_decrease_rank":    c.RankTasks.AuthorFansDecreaseRank,
		"author_potential_rank":        c.RankTasks.AuthorPotentialRank,
		"product_hot_sale_rank":        c.RankTasks.ProductHotSaleRank,
		"product_real_time_sales_rank": c.RankTasks.ProductRealTimeSalesRank,
		"live_author_sales_rank":       c.RankTasks.LiveAuthorSalesRank,
		"live_hot_push_rank":           c.RankTasks.LiveHotPushRank,
		"hot_video_rank":               c.RankTasks.HotVideoRank,
		"ecommerce_video_rank":         c.RankTasks.EcommerceVideoRank,
		"video_hot_push":               c.RankTasks.VideoHotPush,
		"hot_sale_shop":                c.RankTasks.HotSaleShop,
		"site_hourly_rank":             c.RankTasks.SiteHourlyRank,
		"sales_hourly_rank":            c.RankTasks.SalesHourlyRank,
		"real_time_hot_spot":           c.RankTasks.RealTimeHotSpot,
		"soaring_hot_spot":             c.RankTasks.SoaringHotSpot,
		"explore_hot_burst":            c.RankTasks.ExploreHotBurst,
	}
}
