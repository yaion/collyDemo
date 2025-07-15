package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LiveAuthorSalesRank struct {
	Rank         int    `json:"rank" bson:"rank"`
	UID          string `json:"uid" bson:"uid"`
	MarketLevel  int    `json:"market_level" bson:"market_level"`
	NickName     string `json:"nick_name" bson:"nick_name"`
	Avatar       string `json:"avatar" bson:"avatar"`
	DisplayID    string `json:"display_id" bson:"display_id"`
	VerifyType   int    `json:"verify_type" bson:"verify_type"`
	CategoryName string `json:"category_name" bson:"category_name"`
	FansCount    string `json:"fans_count" bson:"fans_count"`
	LiveCount    string `json:"live_count" bson:"live_count"`
	Sales        string `json:"sales" bson:"sales"`
	Gmv          string `json:"gmv" bson:"gmv"`
	MinTime      int    `json:"min_time" bson:"min_time"`
	MaxTime      int    `json:"max_time" bson:"max_time"`
}

type LiveAuthorSalesRankDAO struct {
	collection *mongo.Collection
}

// NewLiveAuthorSalesRankDAO 创建LiveAuthorSalesRank数据访问对象
func NewLiveAuthorSalesRankDAO(db *mongo.Database) *LiveAuthorSalesRankDAO {
	return &LiveAuthorSalesRankDAO{
		collection: db.Collection("live_author_sales_ranks"), // 集合名
	}
}

// Create 创建直播达人带货排名
func (dao *LiveAuthorSalesRankDAO) Create(liveAuthorSalesRank *LiveAuthorSalesRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), liveAuthorSalesRank)
	if err != nil {
		log.Printf("Create live author sales rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *LiveAuthorSalesRankDAO) BatchCreate(ctx context.Context, liveAuthorSalesRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range liveAuthorSalesRanks {
		liveAuthorSalesRank := item.(*LiveAuthorSalesRank)
		filter := bson.M{"uid": liveAuthorSalesRank.UID, "min_time": liveAuthorSalesRank.MinTime, "max_time": liveAuthorSalesRank.MaxTime}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(liveAuthorSalesRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite live author sales ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取直播达人带货排名
func (dao *LiveAuthorSalesRankDAO) GetByID(ctx context.Context, uid string, minTime, maxTime int) (*LiveAuthorSalesRank, error) {
	var liveAuthorSalesRank LiveAuthorSalesRank
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	err := dao.collection.FindOne(ctx, filter).Decode(&liveAuthorSalesRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get live author sales rank error: %v", err)
		return nil, err
	}
	return &liveAuthorSalesRank, nil
}

// Update 更新直播达人带货排名信息
func (dao *LiveAuthorSalesRankDAO) Update(uid string, minTime, maxTime int, updateData *LiveAuthorSalesRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal live author sales rank error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update live author sales rank error: %v", err)
	}
	return err
}

// Delete 删除直播达人带货排名
func (dao *LiveAuthorSalesRankDAO) Delete(uid string, minTime, maxTime int) error {
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete live author sales rank error: %v", err)
	}
	return err
}

// ListAll 获取所有直播达人带货排名（带分页）
func (dao *LiveAuthorSalesRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count live author sales ranks error: %v", err)
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
		log.Printf("List live author sales ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	liveAuthorSalesRanks := make([]LiveAuthorSalesRank, 0)
	if err = cursor.All(ctx, &liveAuthorSalesRanks); err != nil {
		return result, err
	}
	fmt.Println(liveAuthorSalesRanks)
	result["list"] = liveAuthorSalesRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
