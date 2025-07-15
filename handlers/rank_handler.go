package handlers

import (
	"collyDemo/core"
	"collyDemo/mongodb"
	"collyDemo/pkg/utils"
	"context"
	"encoding/json"
	"log"

	"github.com/gocolly/colly/v2"
)

// 达人涨粉榜
type AuthorFansIncreaseRankResult struct {
	IsAuthority bool                              `json:"is_authority"`
	Items       []*mongodb.AuthorFansIncreaseRank `json:"items"`
	Pagination  Pagination                        `json:"pagination"`
	Sort        Sort                              `json:"sort"`
}

func AuthorFansIncreaseRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理达人涨粉榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(AuthorFansIncreaseRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("author_fans_increase_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewAuthorFansIncreaseRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create author fans increase rank error: %v", err)
		return err
	}

	log.Printf("达人涨粉榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 达人掉粉榜
type AuthorFansDecreaseRankResult struct {
	IsAuthority bool                              `json:"is_authority"`
	Items       []*mongodb.AuthorFansDecreaseRank `json:"items"`
	Pagination  Pagination                        `json:"pagination"`
	Sort        Sort                              `json:"sort"`
}

func AuthorFansDecreaseRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理达人掉粉榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(AuthorFansDecreaseRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("author_fans_decrease_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewAuthorFansDecreaseRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create author fans decrease rank error: %v", err)
		return err
	}

	log.Printf("达人掉粉榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 达人带货潜力榜
type AuthorPotentialRankResult struct {
	IsAuthority bool                           `json:"is_authority"`
	Items       []*mongodb.AuthorPotentialRank `json:"items"`
	Pagination  Pagination                     `json:"pagination"`
	Sort        Sort                           `json:"sort"`
}

func AuthorPotentialRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理达人带货潜力榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(AuthorPotentialRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("author_potential_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewAuthorPotentialRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create author potential rank error: %v", err)
		return err
	}

	log.Printf("达人带货潜力榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 商品热销榜
type ProductHotSaleRankResult struct {
	IsAuthority bool                          `json:"is_authority"`
	Items       []*mongodb.ProductHotSaleRank `json:"items"`
	Pagination  Pagination                    `json:"pagination"`
	Sort        Sort                          `json:"sort"`
}

func ProductHotSaleRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理商品热销榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(ProductHotSaleRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("product_hot_sale_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewProductHotSaleRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create product hot sale rank error: %v", err)
		return err
	}

	log.Printf("商品热销榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 商品实时销量榜
type ProductRealTimeSalesRankResult struct {
	IsAuthority bool                                `json:"is_authority"`
	Items       []*mongodb.ProductRealTimeSalesRank `json:"items"`
	Pagination  Pagination                          `json:"pagination"`
	Sort        Sort                                `json:"sort"`
}

func ProductRealTimeSalesRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理商品实时销量榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(ProductRealTimeSalesRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("product_real_time_sales_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewProductRealTimeSalesRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create product real time sales rank error: %v", err)
		return err
	}

	log.Printf("商品实时销量榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 直播达人带货榜
type LiveAuthorSalesRankResult struct {
	IsAuthority bool                           `json:"is_authority"`
	Items       []*mongodb.LiveAuthorSalesRank `json:"items"`
	Pagination  Pagination                     `json:"pagination"`
	Sort        Sort                           `json:"sort"`
}

func LiveAuthorSalesRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理直播达人带货榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(LiveAuthorSalesRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("live_author_sales_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewLiveAuthorSalesRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create live author sales rank error: %v", err)
		return err
	}

	log.Printf("直播达人带货榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 直播热推榜
type LiveHotPushRankResult struct {
	IsAuthority bool                       `json:"is_authority"`
	Items       []*mongodb.LiveHotPushRank `json:"items"`
	Pagination  Pagination                 `json:"pagination"`
	Sort        Sort                       `json:"sort"`
}

func LiveHotPushRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理直播热推榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(LiveHotPushRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("live_hot_push_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewLiveHotPushRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create live hot push rank error: %v", err)
		return err
	}

	log.Printf("直播热推榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 热门视频榜
type HotVideoRankResult struct {
	IsAuthority bool                    `json:"is_authority"`
	Items       []*mongodb.HotVideoRank `json:"items"`
	Pagination  Pagination              `json:"pagination"`
	Sort        Sort                    `json:"sort"`
}

func HotVideoRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理热门视频榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(HotVideoRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("hot_video_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewHotVideoRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create hot video rank error: %v", err)
		return err
	}

	log.Printf("热门视频榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 电商视频榜
type EcommerceVideoRankResult struct {
	IsAuthority bool                          `json:"is_authority"`
	Items       []*mongodb.EcommerceVideoRank `json:"items"`
	Pagination  Pagination                    `json:"pagination"`
	Sort        Sort                          `json:"sort"`
}

func EcommerceVideoRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理电商视频榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(EcommerceVideoRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("ecommerce_video_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewEcommerceVideoRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create ecommerce video rank error: %v", err)
		return err
	}

	log.Printf("电商视频榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 视频热推
type VideoHotPushResult struct {
	IsAuthority bool                    `json:"is_authority"`
	Items       []*mongodb.VideoHotPush `json:"items"`
	Pagination  Pagination              `json:"pagination"`
	Sort        Sort                    `json:"sort"`
}

func VideoHotPushHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理视频热推: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(VideoHotPushResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("video_hot_push.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewVideoHotPushDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create video hot push error: %v", err)
		return err
	}

	log.Printf("视频热推处理完成: %s", r.Request.URL.String())
	return nil
}

// 热销小店
type HotSaleShopResult struct {
	IsAuthority bool                   `json:"is_authority"`
	Items       []*mongodb.HotSaleShop `json:"items"`
	Pagination  Pagination             `json:"pagination"`
	Sort        Sort                   `json:"sort"`
}

func HotSaleShopHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理热销小店: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(HotSaleShopResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("hot_sale_shop.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewHotSaleShopDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create hot sale shop error: %v", err)
		return err
	}

	log.Printf("热销小店处理完成: %s", r.Request.URL.String())
	return nil
}

// 全站小时榜
type SiteHourlyRankResult struct {
	IsAuthority bool                      `json:"is_authority"`
	Items       []*mongodb.SiteHourlyRank `json:"items"`
	Pagination  Pagination                `json:"pagination"`
	Sort        Sort                      `json:"sort"`
}

func SiteHourlyRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理全站小时榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(SiteHourlyRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("site_hourly_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewSiteHourlyRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create site hourly rank error: %v", err)
		return err
	}

	log.Printf("全站小时榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 带货小时榜
type SalesHourlyRankResult struct {
	IsAuthority bool                       `json:"is_authority"`
	Items       []*mongodb.SalesHourlyRank `json:"items"`
	Pagination  Pagination                 `json:"pagination"`
	Sort        Sort                       `json:"sort"`
}

func SalesHourlyRankHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理带货小时榜: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(SalesHourlyRankResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("sales_hourly_rank.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewSalesHourlyRankDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create sales hourly rank error: %v", err)
		return err
	}

	log.Printf("带货小时榜处理完成: %s", r.Request.URL.String())
	return nil
}

// 实时热点
type RealTimeHotSpotResult struct {
	IsAuthority bool                       `json:"is_authority"`
	Items       []*mongodb.RealTimeHotSpot `json:"items"`
	Pagination  Pagination                 `json:"pagination"`
	Sort        Sort                       `json:"sort"`
}

func RealTimeHotSpotHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理实时热点: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(RealTimeHotSpotResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("real_time_hot_spot.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewRealTimeHotSpotDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create real time hot spot error: %v", err)
		return err
	}

	log.Printf("实时热点处理完成: %s", r.Request.URL.String())
	return nil
}

// 飙升热点
type SoaringHotSpotResult struct {
	IsAuthority bool                      `json:"is_authority"`
	Items       []*mongodb.SoaringHotSpot `json:"items"`
	Pagination  Pagination                `json:"pagination"`
	Sort        Sort                      `json:"sort"`
}

func SoaringHotSpotHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理飙升热点: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(SoaringHotSpotResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("soaring_hot_spot.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewSoaringHotSpotDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create soaring hot spot error: %v", err)
		return err
	}

	log.Printf("飙升热点处理完成: %s", r.Request.URL.String())
	return nil
}

// 探测爆款
type ExploreHotBurstResult struct {
	IsAuthority bool                       `json:"is_authority"`
	Items       []*mongodb.ExploreHotBurst `json:"items"`
	Pagination  Pagination                 `json:"pagination"`
	Sort        Sort                       `json:"sort"`
}

func ExploreHotBurstHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理探测爆款: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(ExploreHotBurstResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("explore_hot_burst.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}

	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewExploreHotBurstDAO(db)
	var docs []interface{}
	for _, item := range result.Items {
		docs = append(docs, item)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create explore hot burst error: %v", err)
		return err
	}

	log.Printf("探测爆款处理完成: %s", r.Request.URL.String())
	return nil
}
