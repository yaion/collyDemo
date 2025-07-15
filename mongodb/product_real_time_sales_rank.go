package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRealTimeSalesRank struct {
	PromotionID string  `json:"promotion_id" bson:"promotion_id"`
	ProductID   string  `json:"product_id" bson:"product_id"`
	Title       string  `json:"title" bson:"title"`
	MarketPrice int     `json:"market_price" bson:"market_price"`
	Price       int     `json:"price" bson:"price"`
	PriceStr    string  `json:"price_str" bson:"price_str"`
	Cover       string  `json:"cover" bson:"cover"`
	Url         string  `json:"url" bson:"url"`
	Source      int     `json:"source" bson:"source"`
	CosRatio    float64 `json:"cos_ratio" bson:"cos_ratio"`
	H2Sales     string  `json:"h2_sales" bson:"h2_sales"`
	H2Views     string  `json:"h2_views" bson:"h2_views"`
	TodaySales  string  `json:"today_sales" bson:"today_sales"`
	Rank        int     `json:"rank" bson:"rank"`
}

type ProductRealTimeSalesRankDAO struct {
	collection *mongo.Collection
}

// NewProductRealTimeSalesRankDAO 创建ProductRealTimeSalesRank数据访问对象
func NewProductRealTimeSalesRankDAO(db *mongo.Database) *ProductRealTimeSalesRankDAO {
	return &ProductRealTimeSalesRankDAO{
		collection: db.Collection("product_real_time_sales_ranks"), // 集合名
	}
}

// Create 创建商品实时销量排名
func (dao *ProductRealTimeSalesRankDAO) Create(productRealTimeSalesRank *ProductRealTimeSalesRank) error {
	_, err := dao.collection.InsertOne(context.TODO(), productRealTimeSalesRank)
	if err != nil {
		log.Printf("Create product real time sales rank error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *ProductRealTimeSalesRankDAO) BatchCreate(ctx context.Context, productRealTimeSalesRanks []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range productRealTimeSalesRanks {
		productRealTimeSalesRank := item.(*ProductRealTimeSalesRank)
		filter := bson.M{"promotion_id": productRealTimeSalesRank.PromotionID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(productRealTimeSalesRank).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite product real time sales ranks error: %v", err)
	}
	return err
}

// GetByID 根据ID获取商品实时销量排名
func (dao *ProductRealTimeSalesRankDAO) GetByID(ctx context.Context, promotionID string) (*ProductRealTimeSalesRank, error) {
	var productRealTimeSalesRank ProductRealTimeSalesRank
	filter := bson.M{"promotion_id": promotionID}
	err := dao.collection.FindOne(ctx, filter).Decode(&productRealTimeSalesRank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get product real time sales rank error: %v", err)
		return nil, err
	}
	return &productRealTimeSalesRank, nil
}

// Update 更新商品实时销量排名信息
func (dao *ProductRealTimeSalesRankDAO) Update(promotionID string, updateData *ProductRealTimeSalesRank) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal product real time sales rank error: %v", err)
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
		log.Printf("Update product real time sales rank error: %v", err)
	}
	return err
}

// Delete 删除商品实时销量排名
func (dao *ProductRealTimeSalesRankDAO) Delete(promotionID string) error {
	filter := bson.M{"promotion_id": promotionID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete product real time sales rank error: %v", err)
	}
	return err
}

// ListAll 获取所有商品实时销量排名（带分页）
func (dao *ProductRealTimeSalesRankDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count product real time sales ranks error: %v", err)
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
		log.Printf("List product real time sales ranks error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	productRealTimeSalesRanks := make([]ProductRealTimeSalesRank, 0)
	if err = cursor.All(ctx, &productRealTimeSalesRanks); err != nil {
		return result, err
	}
	fmt.Println(productRealTimeSalesRanks)
	result["list"] = productRealTimeSalesRanks
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
