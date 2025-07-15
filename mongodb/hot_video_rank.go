package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotVideoRank struct {
	FansCount    string `json:"fans_count" bson:"fans_count"`
	Avatar       string `json:"avatar" bson:"avatar"`
	Rank         int    `json:"rank" bson:"rank"`
	DateCode     int    `json:"date_code" bson:"date_code"`
	AwemeID      string `json:"aweme_id" bson:"aweme_id"`
	UID          string `json:"uid" bson:"uid"`
	MarketLevel  int    `json:"market_level" bson:"market_level"`
	NickName     string `json:"nick_name" bson:"nick_name"`
	Desc         string `json:"desc" bson:"desc"`
	Cover        string `json:"cover" bson:"cover"`
	PubTime      int64  `json:"pub_time" bson:"pub_time"`
	LikeCount    string `json:"like_count" bson:"like_count"`
	CommentCount string `json:"comment_count" bson:"comment_count"`
	ShareCount   string `json:"share_count" bson:"share_count"`
	CollectCount string `json:"collect_count" bson:"collect_count"`
	Url          string `json:"url" bson:"url"`
}

type HotVideoRankDAO struct {
	collection *mongo.Collection
}

// NewHotVideoRankDAO 创建HotVideoRank数据访问对象
func NewHotVideoRankDAO(db *mongo.Database) *HotVideoRankDAO {
	return &HotVideoRankDAO{
		collection: db.Collection("hot_video_ranks"), // 集合名
	}
}

// Create 创建热门视频榜
func (dao *HotVideoRankDAO) Create(hotVideoRank *HotVideoRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), hotVideoRank)
	if err != nil {
		log.Printf("Create hot video rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *HotVideoRankDAO) BatchCreate(ctx context.Context, hotVideoRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range hotVideoRanks {
		hotVideoRank := item.(*HotVideoRank)
		filter := bson.M{"aweme_id": hotVideoRank.AwemeID, "date_code": hotVideoRank.DateCode}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(hotVideoRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite hot video ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取热门视频榜
func (dao *HotVideoRankDAO) GetByID(ctx context.Context, awemeID string, dateCode int) (*HotVideoRank, error) {
	var hotVideoRank HotVideoRank
	filter := bson.M{"aweme_id": awemeID, "date_code": dateCode}
	err := dao.collection.FindOne(ctx, filter).Decode(&hotVideoRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get hot video rank error: %v", err)
		return nil, err
	}
	return &hotVideoRank, nil
}

// Update 更新热门视频榜信息
func (dao *HotVideoRankDAO) Update(awemeID string, dateCode int, updateData *HotVideoRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal hot video rank error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"aweme_id": awemeID, "date_code": dateCode}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update hot video rank error: %v", err)
	}
	return err
}

// Delete 删除热门视频榜
func (dao *HotVideoRankDAO) Delete(awemeID string, dateCode int) error {
	filter := bson.M{"aweme_id": awemeID, "date_code": dateCode}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete hot video rank error: %v", err)
	}
	return err
}

// ListAll 获取所有热门视频榜（带分页）
func (dao *HotVideoRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count hot video ranks error: %v", err)
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
		log.Printf("List hot video ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	hotVideoRanks := make([]HotVideoRank, 0)
	if err = cursor.All(ctx, &hotVideoRanks); err != nil {
		return result, err
	}
	fmt.Println(hotVideoRanks)
	result["list"] = hotVideoRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
