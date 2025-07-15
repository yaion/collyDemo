package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EcommerceVideoRank struct {
	Rank            int    `json:"rank" bson:"rank"`
	DateCode        int    `json:"date_code" bson:"date_code"`
	UID             string `json:"uid" bson:"uid"`
	MarketLevel     int    `json:"market_level" bson:"market_level"`
	AwemeID         string `json:"aweme_id" bson:"aweme_id"`
	Desc            string `json:"desc" bson:"desc"`
	Cover           string `json:"cover" bson:"cover"`
	PubTime         int64  `json:"pub_time" bson:"pub_time"`
	Url             string `json:"url" bson:"url"`
	ProductID       string `json:"product_id" bson:"product_id"`
	ProductTitle    string `json:"product_title" bson:"product_title"`
	ProductPrice    int    `json:"product_price" bson:"product_price"`
	ProductPriceStr string `json:"product_price_str" bson:"product_price_str"`
	ProductCover    string `json:"product_cover" bson:"product_cover"`
	Sales           string `json:"sales" bson:"sales"`
	Gmv             string `json:"gmv" bson:"gmv"`
}

type EcommerceVideoRankDAO struct {
	collection *mongo.Collection
}

// NewEcommerceVideoRankDAO 创建EcommerceVideoRank数据访问对象
func NewEcommerceVideoRankDAO(db *mongo.Database) *EcommerceVideoRankDAO {
	return &EcommerceVideoRankDAO{
		collection: db.Collection("ecommerce_video_ranks"), // 集合名
	}
}

// Create 创建电商视频榜
func (dao *EcommerceVideoRankDAO) Create(ecommerceVideoRank *EcommerceVideoRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), ecommerceVideoRank)
	if err != nil {
		log.Printf("Create ecommerce video rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *EcommerceVideoRankDAO) BatchCreate(ctx context.Context, ecommerceVideoRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range ecommerceVideoRanks {
		ecommerceVideoRank := item.(*EcommerceVideoRank)
		filter := bson.M{"aweme_id": ecommerceVideoRank.AwemeID, "date_code": ecommerceVideoRank.DateCode}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(ecommerceVideoRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite ecommerce video ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取电商视频榜
func (dao *EcommerceVideoRankDAO) GetByID(ctx context.Context, awemeID string, dateCode int) (*EcommerceVideoRank, error) {
	var ecommerceVideoRank EcommerceVideoRank
	filter := bson.M{"aweme_id": awemeID, "date_code": dateCode}
	err := dao.collection.FindOne(ctx, filter).Decode(&ecommerceVideoRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get ecommerce video rank error: %v", err)
		return nil, err
	}
	return &ecommerceVideoRank, nil
}

// Update 更新电商视频榜信息
func (dao *EcommerceVideoRankDAO) Update(awemeID string, dateCode int, updateData *EcommerceVideoRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal ecommerce video rank error: %v", err)
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
		log.Printf("Update ecommerce video rank error: %v", err)
	}
	return err
}

// Delete 删除电商视频榜
func (dao *EcommerceVideoRankDAO) Delete(awemeID string, dateCode int) error {
	filter := bson.M{"aweme_id": awemeID, "date_code": dateCode}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete ecommerce video rank error: %v", err)
	}
	return err
}

// ListAll 获取所有电商视频榜（带分页）
func (dao *EcommerceVideoRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count ecommerce video ranks error: %v", err)
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
		log.Printf("List ecommerce video ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	ecommerceVideoRanks := make([]EcommerceVideoRank, 0)
	if err = cursor.All(ctx, &ecommerceVideoRanks); err != nil {
		return result, err
	}
	fmt.Println(ecommerceVideoRanks)
	result["list"] = ecommerceVideoRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
