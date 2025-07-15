package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SalesHourlyRank struct {
	DateCode    int    `json:"date_code" bson:"date_code"`
	RoomID      string `json:"room_id" bson:"room_id"`
	Duration    int    `json:"duration" bson:"duration"`
	UID         string `json:"uid" bson:"uid"`
	MarketLevel int    `json:"market_level" bson:"market_level"`
	NickName    string `json:"nick_name" bson:"nick_name"`
	FansCount   string `json:"fans_count" bson:"fans_count"`
	Avatar      string `json:"avatar" bson:"avatar"`
	Title       string `json:"title" bson:"title"`
	Cover       string `json:"cover" bson:"cover"`
	Score       string `json:"score" bson:"score"`
	Sales       string `json:"sales" bson:"sales"`
	Gmv         string `json:"gmv" bson:"gmv"`
	Rpm         string `json:"rpm" bson:"rpm"`
	Rank        int    `json:"rank" bson:"rank"`
	QrCode      string `json:"qr_code" bson:"qr_code"`
}

type SalesHourlyRankDAO struct {
	collection *mongo.Collection
}

// NewSalesHourlyRankDAO 创建SalesHourlyRank数据访问对象
func NewSalesHourlyRankDAO(db *mongo.Database) *SalesHourlyRankDAO {
	return &SalesHourlyRankDAO{
		collection: db.Collection("sales_hourly_ranks"), // 集合名
	}
}

// Create 创建带货小时榜
func (dao *SalesHourlyRankDAO) Create(salesHourlyRank *SalesHourlyRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), salesHourlyRank)
	if err != nil {
		log.Printf("Create sales hourly rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *SalesHourlyRankDAO) BatchCreate(ctx context.Context, salesHourlyRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range salesHourlyRanks {
		salesHourlyRank := item.(*SalesHourlyRank)
		filter := bson.M{"room_id": salesHourlyRank.RoomID, "date_code": salesHourlyRank.DateCode}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(salesHourlyRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite sales hourly ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取带货小时榜
func (dao *SalesHourlyRankDAO) GetByID(ctx context.Context, roomID string, dateCode int) (*SalesHourlyRank, error) {
	var salesHourlyRank SalesHourlyRank
	filter := bson.M{"room_id": roomID, "date_code": dateCode}
	err := dao.collection.FindOne(ctx, filter).Decode(&salesHourlyRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get sales hourly rank error: %v", err)
		return nil, err
	}
	return &salesHourlyRank, nil
}

// Update 更新带货小时榜信息
func (dao *SalesHourlyRankDAO) Update(roomID string, dateCode int, updateData *SalesHourlyRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal sales hourly rank error: %v", err)
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
		log.Printf("Update sales hourly rank error: %v", err)
	}
	return err
}

// Delete 删除带货小时榜
func (dao *SalesHourlyRankDAO) Delete(roomID string, dateCode int) error {
	filter := bson.M{"room_id": roomID, "date_code": dateCode}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete sales hourly rank error: %v", err)
	}
	return err
}

// ListAll 获取所有带货小时榜（带分页）
func (dao *SalesHourlyRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count sales hourly ranks error: %v", err)
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
		log.Printf("List sales hourly ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	salesHourlyRanks := make([]SalesHourlyRank, 0)
	if err = cursor.All(ctx, &salesHourlyRanks); err != nil {
		return result, err
	}
	fmt.Println(salesHourlyRanks)
	result["list"] = salesHourlyRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
