package handlers

import (
	"collyDemo/core"
	"collyDemo/mongodb"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type StoreResult struct {
	IsAuthority bool             `json:"is_authority"`
	Items       []*mongodb.Store `json:"items"`
	Pagination  Pagination       `json:"pagination"`
	Sort        Sort             `json:"sort"`
}

func StoreHandler(r *colly.Response, acc *core.Account, d *core.TaskDispatcher) error {
	str, err := Handler(r)
	if err != nil {
		return err
	}
	result := new(StoreResult)
	err = json.Unmarshal([]byte(str), result)
	if err != nil {
		return err
	}
	if result.IsAuthority == false {
		return nil
	}
	client := mongodb.GetMongo()
	db := client.Database("kaogujia")
	dao := mongodb.NewStoreDAO(db)
	var docs []interface{}
	for _, author := range result.Items {
		docs = append(docs, author)
	}
	err = dao.BatchCreate(context.Background(), docs)
	if err != nil {
		log.Printf("Create author error: %v", err)
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
	if result.Pagination.TotalCount > result.Pagination.Page*result.Pagination.Limit {
		// 创建任务
		listTask := &core.Task{
			URL:     fmt.Sprintf("https://service.kaogujia.com/api/shop/search?limit=%v&page=%v&sort_field=gmv&sort=0", result.Pagination.Limit, result.Pagination.Page+1),
			Method:  "POST",
			Headers: headers,
			Body:    []byte(`{"period":1,"keyword":""}`),
			Handler: StoreHandler,
			Meta: map[string]interface{}{
				"page":     1,
				"pageSize": 50,
			},
		}
		d.AddTask(listTask)
	}
	// todo 创建info任务

	return nil
}
