package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductHotSaleRank struct {
	PromotionID string  `json:"promotion_id" bson:"promotion_id"`
	Title       string  `json:"title" bson:"title"`
	MarketPrice int     `json:"market_price" bson:"market_price"`
	Price       int     `json:"price" bson:"price"`
	PriceStr    string  `json:"price_str" bson:"price_str"`
	Cover       string  `json:"cover" bson:"cover"`
	Url         string  `json:"url" bson:"url"`
	Source      int     `json:"source" bson:"source"`
	CosRatio    float64 `json:"cos_ratio" bson:"cos_ratio"`
	Sales       string  `json:"sales" bson:"sales"`
	Views       string  `json:"views" bson:"views"`
	Rank        int     `json:"rank" bson:"rank"`
}

type ProductHotSaleRankDAO struct {
	collection *mongo.Collection
}

// NewProductHotSaleRankDAO 创建ProductHotSaleRank数据访问对象
func NewProductHotSaleRankDAO(db *mongo.Database) *ProductHotSaleRankDAO {
	return &ProductHotSaleRankDAO{
		collection: db.Collection("product_hot_sale_ranks"), // 集合名
	}
}

// Create 创建商品热销排名
func (dao *ProductHotSaleRankDAO) Create(productHotSaleRank *ProductHotSaleRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), productHotSaleRank)
	if err != nil {
		log.Printf("Create product hot sale rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *ProductHotSaleRankDAO) BatchCreate(ctx context.Context, productHotSaleRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range productHotSaleRanks {
		productHotSaleRank := item.(*ProductHotSaleRank)
		filter := bson.M{"promotion_id": productHotSaleRank.PromotionID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(productHotSaleRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite product hot sale ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取商品热销排名
func (dao *ProductHotSaleRankDAO) GetByID(ctx context.Context, promotionID string) (*ProductHotSaleRank, error) {
	var productHotSaleRank ProductHotSaleRank
	filter := bson.M{"promotion_id": promotionID}
	err := dao.collection.FindOne(ctx, filter).Decode(&productHotSaleRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get product hot sale rank error: %v", err)
		return nil, err
	}
	return &productHotSaleRank, nil
}

// Update 更新商品热销排名信息
func (dao *ProductHotSaleRankDAO) Update(promotionID string, updateData *ProductHotSaleRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal product hot sale rank error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"promotion_id": promotionID}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update product hot sale rank error: %v", err)
	}
	return err
}

// Delete 删除商品热销排名
func (dao *ProductHotSaleRankDAO) Delete(promotionID string) error {
	filter := bson.M{"promotion_id": promotionID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete product hot sale rank error: %v", err)
	}
	return err
}

// ListAll 获取所有商品热销排名（带分页）
func (dao *ProductHotSaleRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count product hot sale ranks error: %v", err)
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
		log.Printf("List product hot sale ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	productHotSaleRanks := make([]ProductHotSaleRank, 0)
	if err = cursor.All(ctx, &productHotSaleRanks); err != nil {
		return result, err
	}
	fmt.Println(productHotSaleRanks)
	result["list"] = productHotSaleRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
