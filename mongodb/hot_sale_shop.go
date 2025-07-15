package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotSaleShop struct {
	ShopID      string `json:"shop_id" bson:"shop_id"`
	Name        string `json:"name" bson:"name"`
	Logo        string `json:"logo" bson:"logo"`
	SkuCount    string `json:"sku_count" bson:"sku_count"`
	Dsr         string `json:"dsr" bson:"dsr"`
	ScoreRating string `json:"score_rating" bson:"score_rating"`
	Sales       string `json:"sales" bson:"sales"`
	Gmv         string `json:"gmv" bson:"gmv"`
	AuthorCount string `json:"author_count" bson:"author_count"`
	LiveCount   string `json:"live_count" bson:"live_count"`
	VideoCount  string `json:"video_count" bson:"video_count"`
	Rank        int    `json:"rank" bson:"rank"`
}

type HotSaleShopDAO struct {
	collection *mongo.Collection
}

// NewHotSaleShopDAO 创建HotSaleShop数据访问对象
func NewHotSaleShopDAO(db *mongo.Database) *HotSaleShopDAO {
	return &HotSaleShopDAO{
		collection: db.Collection("hot_sale_shops"), // 集合名
	}
}

// Create 创建热销小店
func (dao *HotSaleShopDAO) Create(hotSaleShop *HotSaleShop) error {
	_, err := dao.collection.InsertOne(context.TODO(), hotSaleShop)
	if err != nil {
		log.Printf("Create hot sale shop error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *HotSaleShopDAO) BatchCreate(ctx context.Context, hotSaleShops []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range hotSaleShops {
		hotSaleShop := item.(*HotSaleShop)
		filter := bson.M{"shop_id": hotSaleShop.ShopID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(hotSaleShop).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite hot sale shops error: %v", err)
	}
	return err
}

// GetByID 根据ID获取热销小店
func (dao *HotSaleShopDAO) GetByID(ctx context.Context, shopID string) (*HotSaleShop, error) {
	var hotSaleShop HotSaleShop
	filter := bson.M{"shop_id": shopID}
	err := dao.collection.FindOne(ctx, filter).Decode(&hotSaleShop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get hot sale shop error: %v", err)
		return nil, err
	}
	return &hotSaleShop, nil
}

// Update 更新热销小店信息
func (dao *HotSaleShopDAO) Update(shopID string, updateData *HotSaleShop) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal hot sale shop error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"shop_id": shopID}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update hot sale shop error: %v", err)
	}
	return err
}

// Delete 删除热销小店
func (dao *HotSaleShopDAO) Delete(shopID string) error {
	filter := bson.M{"shop_id": shopID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete hot sale shop error: %v", err)
	}
	return err
}

// ListAll 获取所有热销小店（带分页）
func (dao *HotSaleShopDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count hot sale shops error: %v", err)
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
		log.Printf("List hot sale shops error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	hotSaleShops := make([]HotSaleShop, 0)
	if err = cursor.All(ctx, &hotSaleShops); err != nil {
		return result, err
	}
	fmt.Println(hotSaleShops)
	result["list"] = hotSaleShops
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
