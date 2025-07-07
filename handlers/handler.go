package handlers

import (
	"collyDemo/pkg/utils"
	"encoding/json"
	"github.com/gocolly/colly/v2"
)

type Result struct {
	Code    int
	Message string
	Success bool
	Data    string
}

type Pagination struct {
	TotalCount int64 `json:"total_count"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
}

type Sort struct {
	SortField string `json:"sort_field"`
	Sort      int    `json:"sort"`
}

func Handler(r *colly.Response) (string, error) {
	result := new(Result)
	err := json.Unmarshal(r.Body, result)
	if err != nil {
		// todo 记录日志
		return "", err
	}

	str, err := utils.Decrypt(r.Request.URL.Path, result.Data)
	if err != nil {
		// todo 记录日志
		return "", err
	}

	return str, nil
}
