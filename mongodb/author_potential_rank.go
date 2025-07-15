package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthorPotentialRank struct {
	Rank         int    `json:"rank" bson:"rank"`
	UID          string `json:"uid" bson:"uid"`
	MarketLevel  int    `json:"market_level" bson:"market_level"`
	NickName     string `json:"nick_name" bson:"nick_name"`
	Avatar       string `json:"avatar" bson:"avatar"`
	DisplayID    string `json:"display_id" bson:"display_id"`
	VerifyType   int    `json:"verify_type" bson:"verify_type"`
	CategoryName string `json:"category_name" bson:"category_name"`
	FansCount    string `json:"fans_count" bson:"fans_count"`
	VideoCount   string `json:"video_count" bson:"video_count"`
	LiveCount    string `json:"live_count" bson:"live_count"`
	Score        string `json:"score" bson:"score"`
	MinTime      int    `json:"min_time" bson:"min_time"`
	MaxTime      int    `json:"max_time" bson:"max_time"`
}

type AuthorPotentialRankDAO struct {
	collection *mongo.Collection
}

// NewAuthorPotentialRankDAO 创建AuthorPotentialRank数据访问对象
func NewAuthorPotentialRankDAO(db *mongo.Database) *AuthorPotentialRankDAO {
	return &AuthorPotentialRankDAO{
		collection: db.Collection("author_potential_ranks"), // 集合名
	}
}

// Create 创建达人带货潜力排名
func (dao *AuthorPotentialRankDAO) Create(authorPotentialRank *AuthorPotentialRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), authorPotentialRank)
	if err != nil {
		log.Printf("Create author potential rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *AuthorPotentialRankDAO) BatchCreate(ctx context.Context, authorPotentialRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range authorPotentialRanks {
		authorPotentialRank := item.(*AuthorPotentialRank)
		filter := bson.M{"uid": authorPotentialRank.UID, "min_time": authorPotentialRank.MinTime, "max_time": authorPotentialRank.MaxTime}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(authorPotentialRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite author potential ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取达人带货潜力排名
func (dao *AuthorPotentialRankDAO) GetByID(ctx context.Context, uid string, minTime, maxTime int) (*AuthorPotentialRank, error) {
	var authorPotentialRank AuthorPotentialRank
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	err := dao.collection.FindOne(ctx, filter).Decode(&authorPotentialRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get author potential rank error: %v", err)
		return nil, err
	}
	return &authorPotentialRank, nil
}

// Update 更新达人带货潜力排名信息
func (dao *AuthorPotentialRankDAO) Update(uid string, minTime, maxTime int, updateData *AuthorPotentialRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal author potential rank error: %v", err)
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
		log.Printf("Update author potential rank error: %v", err)
	}
	return err
}

// Delete 删除达人带货潜力排名
func (dao *AuthorPotentialRankDAO) Delete(uid string, minTime, maxTime int) error {
	filter := bson.M{"uid": uid, "min_time": minTime, "max_time": maxTime}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete author potential rank error: %v", err)
	}
	return err
}

// ListAll 获取所有达人带货潜力排名（带分页）
func (dao *AuthorPotentialRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count author potential ranks error: %v", err)
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
		log.Printf("List author potential ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	authorPotentialRanks := make([]AuthorPotentialRank, 0)
	if err = cursor.All(ctx, &authorPotentialRanks); err != nil {
		return result, err
	}
	fmt.Println(authorPotentialRanks)
	result["list"] = authorPotentialRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
