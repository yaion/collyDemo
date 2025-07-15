package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoHotPush struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	Title       string  `json:"title" bson:"title"`
	MarketPrice int     `json:"market_price" bson:"market_price"`
	Price       int     `json:"price" bson:"price"`
	PriceStr    string  `json:"price_str" bson:"price_str"`
	Source      int     `json:"source" bson:"source"`
	Cover       string  `json:"cover" bson:"cover"`
	Url         string  `json:"url" bson:"url"`
	CosRatio    float64 `json:"cos_ratio" bson:"cos_ratio"`
	Sales       string  `json:"sales" bson:"sales"`
	Views       string  `json:"views" bson:"views"`
	VideoCount  string  `json:"video_count" bson:"video_count"`
	Rank        int     `json:"rank" bson:"rank"`
}

type VideoHotPushDAO struct {
	collection *mongo.Collection
}

// NewVideoHotPushDAO 创建VideoHotPush数据访问对象
func NewVideoHotPushDAO(db *mongo.Database) *VideoHotPushDAO {
	return &VideoHotPushDAO{
		collection: db.Collection("video_hot_pushes"), // 集合名
	}
}

// Create 创建视频热推
func (dao *VideoHotPushDAO) Create(videoHotPush *VideoHotPush) error {
	_, err := dao.collection.InsertOne(context.TODO(), videoHotPush)
	if err != nil {
		log.Printf("Create video hot push error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *VideoHotPushDAO) BatchCreate(ctx context.Context, videoHotPushes []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range videoHotPushes {
		videoHotPush := item.(*VideoHotPush)
		filter := bson.M{"product_id": videoHotPush.ProductID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(videoHotPush).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite video hot pushes error: %v", err)
	}
	return err
}

// GetByID 根据ID获取视频热推
func (dao *VideoHotPushDAO) GetByID(ctx context.Context, productID string) (*VideoHotPush, error) {
	var videoHotPush VideoHotPush
	filter := bson.M{"product_id": productID}
	err := dao.collection.FindOne(ctx, filter).Decode(&videoHotPush)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get video hot push error: %v", err)
		return nil, err
	}
	return &videoHotPush, nil
}

// Update 更新视频热推信息
func (dao *VideoHotPushDAO) Update(productID string, updateData *VideoHotPush) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal video hot push error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"product_id": productID}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update video hot push error: %v", err)
	}
	return err
}

// Delete 删除视频热推
func (dao *VideoHotPushDAO) Delete(productID string) error {
	filter := bson.M{"product_id": productID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete video hot push error: %v", err)
	}
	return err
}

// ListAll 获取所有视频热推（带分页）
func (dao *VideoHotPushDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count video hot pushes error: %v", err)
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
		log.Printf("List video hot pushes error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	videoHotPushes := make([]VideoHotPush, 0)
	if err = cursor.All(ctx, &videoHotPushes); err != nil {
		return result, err
	}
	fmt.Println(videoHotPushes)
	result["list"] = videoHotPushes
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
