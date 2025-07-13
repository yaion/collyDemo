package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Store 店铺实体
type Store struct {
	ShopID      string      `json:"shop_id" bson:"shop_id"` // 使用shop_id作为主键
	Name        string      `json:"name" bson:"name"`
	Logo        string      `json:"logo" bson:"logo"`
	HasFlagship int         `json:"has_flagship" bson:"has_flagship"` // 0-无，1-有
	DsrStr      string      `json:"dsr_str" bson:"dsr_str"`
	Dsr         float64     `json:"dsr" bson:"dsr"`
	Lv1         string      `json:"lv1" bson:"lv1"` // 一级类目
	SkuCount    string      `json:"sku_count" bson:"sku_count"`
	IsFollowed  bool        `json:"is_followed" bson:"is_followed"`
	Stat        interface{} `json:"stat" bson:"stat"`
}

// StoreDAO 店铺数据访问对象
type StoreDAO struct {
	collection *mongo.Collection
}

// NewStoreDAO 创建新的StoreDAO实例
func NewStoreDAO(db *mongo.Database) *StoreDAO {
	return &StoreDAO{
		collection: db.Collection("stores"),
	}
}

func (dao *StoreDAO) BatchCreate(ctx context.Context, stores []interface{}) error {
	_, err := dao.collection.InsertMany(ctx, stores)
	return err
}

// Create 创建店铺记录
func (dao *StoreDAO) Create(ctx context.Context, store *Store) error {
	_, err := dao.collection.InsertOne(ctx, store)
	return err
}

// GetByShopID 根据店铺ID查询
func (dao *StoreDAO) GetByShopID(ctx context.Context, shopID string) (*Store, error) {
	var store Store
	err := dao.collection.FindOne(ctx, bson.M{"_id": shopID}).Decode(&store)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &store, err
}

// Update 更新店铺信息
func (dao *StoreDAO) Update(ctx context.Context, shopID string, updateData bson.M) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": shopID},
		bson.M{"$set": updateData},
		options.Update().SetUpsert(false),
	)
	return err
}

// Delete 删除店铺记录
func (dao *StoreDAO) Delete(ctx context.Context, shopID string) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": shopID})
	return err
}

func (dao *StoreDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count authors error: %v", err)
		return result, err
	}
	result["total"] = total

	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	// 添加默认排序（按粉丝数降序）
	findOptions.SetSort(bson.D{{Key: "fans", Value: -1}})

	cursor, err := dao.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("List authors error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	stores := make([]Store, 0)
	if err = cursor.All(ctx, &stores); err != nil {
		return result, err
	}
	result["list"] = stores
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
