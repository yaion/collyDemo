package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LiveHotPushRank struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	PromotionID string  `json:"promotion_id" bson:"promotion_id"`
	Title       string  `json:"title" bson:"title"`
	MarketPrice int     `json:"market_price" bson:"market_price"`
	Price       int     `json:"price" bson:"price"`
	PriceStr    string  `json:"price_str" bson:"price_str"`
	Cover       string  `json:"cover" bson:"cover"`
	Url         string  `json:"url" bson:"url"`
	CosRatio    float64 `json:"cos_ratio" bson:"cos_ratio"`
	Sales       string  `json:"sales" bson:"sales"`
	Gmv         string  `json:"gmv" bson:"gmv"`
	Lives       string  `json:"lives" bson:"lives"`
	Rank        int     `json:"rank" bson:"rank"`
	Source      int     `json:"source" bson:"source"`
}

type LiveHotPushRankDAO struct {
	collection *mongo.Collection
}

// NewLiveHotPushRankDAO 创建LiveHotPushRank数据访问对象
func NewLiveHotPushRankDAO(db *mongo.Database) *LiveHotPushRankDAO {
	return &LiveHotPushRankDAO{
		collection: db.Collection("live_hot_push_ranks"), // 集合名
	}
}

// Create 创建直播热推榜
func (dao *LiveHotPushRankDAO) Create(liveHotPushRank *LiveHotPushRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), liveHotPushRank)
	if err != nil {
		log.Printf("Create live hot push rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *LiveHotPushRankDAO) BatchCreate(ctx context.Context, liveHotPushRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range liveHotPushRanks {
		liveHotPushRank := item.(*LiveHotPushRank)
		filter := bson.M{"product_id": liveHotPushRank.ProductID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(liveHotPushRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite live hot push ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取直播热推榜
func (dao *LiveHotPushRankDAO) GetByID(ctx context.Context, productID string) (*LiveHotPushRank, error) {
	var liveHotPushRank LiveHotPushRank
	filter := bson.M{"product_id": productID}
	err := dao.collection.FindOne(ctx, filter).Decode(&liveHotPushRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get live hot push rank error: %v", err)
		return nil, err
	}
	return &liveHotPushRank, nil
}

// Update 更新直播热推榜信息
func (dao *LiveHotPushRankDAO) Update(productID string, updateData *LiveHotPushRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal live hot push rank error: %v", err)
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
		log.Printf("Update live hot push rank error: %v", err)
	}
	return err
}

// Delete 删除直播热推榜
func (dao *LiveHotPushRankDAO) Delete(productID string) error {
	filter := bson.M{"product_id": productID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete live hot push rank error: %v", err)
	}
	return err
}

// ListAll 获取所有直播热推榜（带分页）
func (dao *LiveHotPushRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count live hot push ranks error: %v", err)
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
		log.Printf("List live hot push ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	liveHotPushRanks := make([]LiveHotPushRank, 0)
	if err = cursor.All(ctx, &liveHotPushRanks); err != nil {
		return result, err
	}
	fmt.Println(liveHotPushRanks)
	result["list"] = liveHotPushRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
