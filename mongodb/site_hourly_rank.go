package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SiteHourlyRank struct {
	DateCode       int    `json:"date_code" bson:"date_code"`
	RoomID         string `json:"room_id" bson:"room_id"`
	DisplayID      string `json:"display_id" bson:"display_id"`
	UID            string `json:"uid" bson:"uid"`
	MarketLevel    int    `json:"market_level" bson:"market_level"`
	NickName       string `json:"nick_name" bson:"nick_name"`
	FansCount      string `json:"fans_count" bson:"fans_count"`
	Avatar         string `json:"avatar" bson:"avatar"`
	Title          string `json:"title" bson:"title"`
	Cover          string `json:"cover" bson:"cover"`
	CategoryName   string `json:"category_name" bson:"category_name"`
	UserCount      string `json:"user_count" bson:"user_count"`
	IsLive         int    `json:"is_live" bson:"is_live"`
	TotalUsers     string `json:"total_users" bson:"total_users"`
	AvgUsers       string `json:"avg_users" bson:"avg_users"`
	QrCode         string `json:"qr_code" bson:"qr_code"`
	Rank           int    `json:"rank" bson:"rank"`
	GapDescription string `json:"gap_description" bson:"gap_description"`
}

type SiteHourlyRankDAO struct {
	collection *mongo.Collection
}

// NewSiteHourlyRankDAO 创建SiteHourlyRank数据访问对象
func NewSiteHourlyRankDAO(db *mongo.Database) *SiteHourlyRankDAO {
	return &SiteHourlyRankDAO{
		collection: db.Collection("site_hourly_ranks"), // 集合名
	}
}

// Create 创建全站小时榜
func (dao *SiteHourlyRankDAO) Create(siteHourlyRank *SiteHourlyRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), siteHourlyRank)
	if err != nil {
		log.Printf("Create site hourly rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *SiteHourlyRankDAO) BatchCreate(ctx context.Context, siteHourlyRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range siteHourlyRanks {
		siteHourlyRank := item.(*SiteHourlyRank)
		filter := bson.M{"room_id": siteHourlyRank.RoomID, "date_code": siteHourlyRank.DateCode}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(siteHourlyRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite site hourly ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取全站小时榜
func (dao *SiteHourlyRankDAO) GetByID(ctx context.Context, roomID string, dateCode int) (*SiteHourlyRank, error) {
	var siteHourlyRank SiteHourlyRank
	filter := bson.M{"room_id": roomID, "date_code": dateCode}
	err := dao.collection.FindOne(ctx, filter).Decode(&siteHourlyRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get site hourly rank error: %v", err)
		return nil, err
	}
	return &siteHourlyRank, nil
}

// Update 更新全站小时榜信息
func (dao *SiteHourlyRankDAO) Update(roomID string, dateCode int, updateData *SiteHourlyRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal site hourly rank error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"room_id": roomID, "date_code": dateCode}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update site hourly rank error: %v", err)
	}
	return err
}

// Delete 删除全站小时榜
func (dao *SiteHourlyRankDAO) Delete(roomID string, dateCode int) error {
	filter := bson.M{"room_id": roomID, "date_code": dateCode}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete site hourly rank error: %v", err)
	}
	return err
}

// ListAll 获取所有全站小时榜（带分页）
func (dao *SiteHourlyRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count site hourly ranks error: %v", err)
		return result, err
	}
	result["total"] = total

	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	// 添加默认排序（按排名升序）
	findOptions.SetSort(bson.D{{Key: "rank", Value: 1}})

	cursor, err := dao.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("List site hourly ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	siteHourlyRanks := make([]SiteHourlyRank, 0)
	if err = cursor.All(ctx, &siteHourlyRanks); err != nil {
		return result, err
	}
	fmt.Println(siteHourlyRanks)
	result["list"] = siteHourlyRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
