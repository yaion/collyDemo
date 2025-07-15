package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthorFansDecreaseRank struct {
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

type AuthorFansDecreaseRankDAO struct {
	collection *mongo.Collection
}

// NewAuthorFansDecreaseRankDAO 创建AuthorFansDecreaseRank数据访问对象
func NewAuthorFansDecreaseRankDAO(db *mongo.Database) *AuthorFansDecreaseRankDAO {
	return &AuthorFansDecreaseRankDAO{
		collection: db.Collection("author_fans_decrease_ranks"), // 集合名
	}
}

// Create 创建达人掉粉排名
func (dao *AuthorFansDecreaseRankDAO) Create(authorFansDecreaseRank *AuthorFansDecreaseRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), authorFansDecreaseRank)
	if err != nil {
		log.Printf("Create author fans decrease rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *AuthorFansDecreaseRankDAO) BatchCreate(ctx context.Context, authorFansDecreaseRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range authorFansDecreaseRanks {
		authorFansDecreaseRank := item.(*AuthorFansDecreaseRank)
		filter := bson.M{"uid": authorFansDecreaseRank.UID, "min_time": authorFansDecreaseRank.MinTime, "max_time": authorFansDecreaseRank.MaxTime}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(authorFansDecreaseRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite author fans decrease ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取达人掉粉排名
func (dao *AuthorFansDecreaseRankDAO) GetByID(ctx context.Context, uid string, minTime, maxTime int) (*AuthorFansDecreaseRank, error) {
	var authorFansDecreaseRank AuthorFansDecreaseRank
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	err := dao.collection.FindOne(ctx, filter).Decode(&authorFansDecreaseRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get author fans decrease rank error: %v", err)
		return nil, err
	}
	return &authorFansDecreaseRank, nil
}

// Update 更新达人掉粉排名信息
func (dao *AuthorFansDecreaseRankDAO) Update(uid string, minTime, maxTime int, updateData *AuthorFansDecreaseRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal author fans decrease rank error: %v", err)
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
		log.Printf("Update author fans decrease rank error: %v", err)
	}
	return err
}

// Delete 删除达人掉粉排名
func (dao *AuthorFansDecreaseRankDAO) Delete(uid string, minTime, maxTime int) error {
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete author fans decrease rank error: %v", err)
	}
	return err
}

// ListAll 获取所有达人掉粉排名（带分页）
func (dao *AuthorFansDecreaseRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count author fans decrease ranks error: %v", err)
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
		log.Printf("List author fans decrease ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	authorFansDecreaseRanks := make([]AuthorFansDecreaseRank, 0)
	if err = cursor.All(ctx, &authorFansDecreaseRanks); err != nil {
		return result, err
	}
	fmt.Println(authorFansDecreaseRanks)
	result["list"] = authorFansDecreaseRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
