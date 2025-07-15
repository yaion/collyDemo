package handlers

import (
	"collyDemo/core"
	"collyDemo/mongodb"
	"collyDemo/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type VideoResult struct {
	IsAuthority bool             `json:"is_authority"`
	Items       []*mongodb.Video `json:"items"`
	Pagination  Pagination       `json:"pagination"`
	Sort        Sort             `json:"sort"`
}

func VideoHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理视频列表: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(VideoResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("video.json", str)
		return err
	}
	if result.IsAuthority == false {
		return nil
	}
	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewVideoDAO(db)
	var docs []interface{}
	for _, video := range result.Items {
		docs = append(docs, video)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create video error: %v", err)
		return err
	}
	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "zh-HK,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7",
		"authorization":      acc.Token,
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

	// 处理分页
	if result.Pagination.TotalCount > result.Pagination.Page*result.Pagination.Limit {
		// 创建下一页任务
		listTask := &core.Task{
			URL:     fmt.Sprintf("https://service.kaogujia.com/api/video/search?limit=%v&page=%v&sort_field=like_count&sort=0", result.Pagination.Limit, result.Pagination.Page+1),
			Method:  "POST",
			Headers: headers,
			Body:    []byte(`{"date_code":{"min":"20250629","max":"20250705"},"keyword":"","video_type":1}`),
			Handler: VideoHandler,
			Meta: map[string]interface{}{
				"page":  result.Pagination.Page + 1,
				"limit": result.Pagination.Limit,
			},
		}
		d.AddTask(listTask)
	}

	// 获取视频详情数据
	for _, item := range result.Items {
		log.Printf("处理视频详情: ID=%s, Desc=%s", item.AwemeID, item.Desc)
		// 创建详情任务
		infoTask := &core.Task{
			URL:     fmt.Sprintf("https://service.kaogujia.com/api/video/detail/%s", item.AwemeID),
			Method:  "GET",
			Headers: headers,
			Handler: VideoInfoHandler,
			Meta: map[string]interface{}{
				"aweme_id": item.AwemeID,
			},
		}
		d.AddTask(infoTask)
	}

	log.Printf("视频列表处理完成: %s", r.Request.URL.String())
	return nil
}

func VideoInfoHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	log.Printf("处理视频详情: %s", r.Request.URL.String())
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(mongodb.Video)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		log.Printf("Unmarshal error: %v,str : %v", err, str)
		utils.WriteToFile("video_info.json", str)
		return err
	}
	//  插入详情数据
	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewVideoDAO(db)
	err = dao.Create(context.Background(), result)
	if err != nil {
		log.Printf("Create video info error: %v", err)
		return err
	}

	log.Printf("视频详情处理完成: %s", r.Request.URL.String())
	return nil
}
