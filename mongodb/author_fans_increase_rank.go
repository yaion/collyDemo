package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthorFansIncreaseRank struct {
	Rank         int    `json:"rank" bson:"rank"`
	UID          string `json:"uid" bson:"uid"`
	MarketLevel  int    `json:"market_level" bson:"market_level"`
	NickName     string `json:"nick_name" bson:"nick_name"`
	Avatar       string `json:"avatar" bson:"avatar"`
	DisplayID    string `json:"display_id" bson:"display_id"`
	VerifyType   int    `json:"verify_type" bson:"verify_type"`
	CategoryName string `json:"category_name" bson:"category_name"`
	IncFansCount string `json:"inc_fans_count" bson:"inc_fans_count"`
	FansCount    string `json:"fans_count" bson:"fans_count"`
	VideoCount   string `json:"video_count" bson:"video_count"`
	MinTime      int    `json:"min_time" bson:"min_time"`
	MaxTime      int    `json:"max_time" bson:"max_time"`
}

type AuthorFansIncreaseRankDAO struct {
	collection *mongo.Collection
}

// NewAuthorFansIncreaseRankDAO 创建AuthorFansIncreaseRank数据访问对象
func NewAuthorFansIncreaseRankDAO(db *mongo.Database) *AuthorFansIncreaseRankDAO {
	return &AuthorFansIncreaseRankDAO{
		collection: db.Collection("author_fans_increase_ranks"), // 集合名
	}
}

// Create 创建达人涨粉排名
func (dao *AuthorFansIncreaseRankDAO) Create(authorFansIncreaseRank *AuthorFansIncreaseRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), authorFansIncreaseRank)
	if err != nil {
		log.Printf("Create author fans increase rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *AuthorFansIncreaseRankDAO) BatchCreate(ctx context.Context, authorFansIncreaseRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range authorFansIncreaseRanks {
		authorFansIncreaseRank := item.(*AuthorFansIncreaseRank)
		filter := bson.M{"uid": authorFansIncreaseRank.UID, "min_time": authorFansIncreaseRank.MinTime, "max_time": authorFansIncreaseRank.MaxTime}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(authorFansIncreaseRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite author fans increase ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取达人涨粉排名
func (dao *AuthorFansIncreaseRankDAO) GetByID(ctx context.Context, uid string, minTime, maxTime int) (*AuthorFansIncreaseRank, error) {
	var authorFansIncreaseRank AuthorFansIncreaseRank
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	err := dao.collection.FindOne(ctx, filter).Decode(&authorFansIncreaseRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get author fans increase rank error: %v", err)
		return nil, err
	}
	return &authorFansIncreaseRank, nil
}

// Update 更新达人涨粉排名信息
func (dao *AuthorFansIncreaseRankDAO) Update(uid string, minTime, maxTime int, updateData *AuthorFansIncreaseRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal author fans increase rank error: %v", err)
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
		log.Printf("Update author fans increase rank error: %v", err)
	}
	return err
}

// Delete 删除达人涨粉排名
func (dao *AuthorFansIncreaseRankDAO) Delete(uid string, minTime, maxTime int) error {
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete author fans increase rank error: %v", err)
	}
	return err
}

// ListAll 获取所有达人涨粉排名（带分页）
func (dao *AuthorFansIncreaseRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count author fans increase ranks error: %v", err)
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
		log.Printf("List author fans increase ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	authorFansIncreaseRanks := make([]AuthorFansIncreaseRank, 0)
	if err = cursor.All(ctx, &authorFansIncreaseRanks); err != nil {
		return result, err
	}
	fmt.Println(authorFansIncreaseRanks)
	result["list"] = authorFansIncreaseRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
