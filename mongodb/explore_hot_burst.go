package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopStat struct {
	Sales      string     `json:"sales" bson:"sales"`
	Gmv        string     `json:"gmv" bson:"gmv"`
	LiveSales  string     `json:"live_sales" bson:"live_sales"`
	LiveGmv    string     `json:"live_gmv" bson:"live_gmv"`
	VideoSales string     `json:"video_sales" bson:"video_sales"`
	OtherSales string     `json:"other_sales" bson:"other_sales"`
	OtherGmv   string     `json:"other_gmv" bson:"other_gmv"`
	VideoGmv   string     `json:"video_gmv" bson:"video_gmv"`
	Lives      int        `json:"lives" bson:"lives"`
	Videos     int        `json:"videos" bson:"videos"`
	Users      int        `json:"users" bson:"users"`
	Uv         string     `json:"uv" bson:"uv"`
	Atv        string     `json:"atv" bson:"atv"`
	MarketType MarketType `json:"market_type" bson:"market_type"`
}

type ExploreHotBurst struct {
	ShopID      string   `json:"shop_id" bson:"shop_id"`
	Name        string   `json:"name" bson:"name"`
	Logo        string   `json:"logo" bson:"logo"`
	HasFlagship int      `json:"has_flagship" bson:"has_flagship"`
	Stat        ShopStat `json:"stat" bson:"stat"`
	DsrStr      string   `json:"dsr_str" bson:"dsr_str"`
	Dsr         float64  `json:"dsr" bson:"dsr"`
	Lv1         string   `json:"lv1" bson:"lv1"`
	SkuCount    string   `json:"sku_count" bson:"sku_count"`
	IsFollowed  bool     `json:"is_followed" bson:"is_followed"`
	Rank        int      `json:"rank" bson:"rank"`
}

type ExploreHotBurstDAO struct {
	collection *mongo.Collection
}

// NewExploreHotBurstDAO 创建ExploreHotBurst数据访问对象
func NewExploreHotBurstDAO(db *mongo.Database) *ExploreHotBurstDAO {
	return &ExploreHotBurstDAO{
		collection: db.Collection("explore_hot_bursts"), // 集合名
	}
}

// Create 创建探测爆款
func (dao *ExploreHotBurstDAO) Create(exploreHotBurst *ExploreHotBurst) error {
	_, err := dao.collection.InsertOne(context.TODO(), exploreHotBurst)
	if err != nil {
		log.Printf("Create explore hot burst error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *ExploreHotBurstDAO) BatchCreate(ctx context.Context, exploreHotBursts []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range exploreHotBursts {
		exploreHotBurst := item.(*ExploreHotBurst)
		filter := bson.M{"shop_id": exploreHotBurst.ShopID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(exploreHotBurst).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite explore hot bursts error: %v", err)
	}
	return err
}

// GetByID 根据ID获取探测爆款
func (dao *ExploreHotBurstDAO) GetByID(ctx context.Context, shopID string) (*ExploreHotBurst, error) {
	var exploreHotBurst ExploreHotBurst
	filter := bson.M{"shop_id": shopID}
	err := dao.collection.FindOne(ctx, filter).Decode(&exploreHotBurst)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get explore hot burst error: %v", err)
		return nil, err
	}
	return &exploreHotBurst, nil
}

// Update 更新探测爆款信息
func (dao *ExploreHotBurstDAO) Update(shopID string, updateData *ExploreHotBurst) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal explore hot burst error: %v", err)
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
		log.Printf("Update explore hot burst error: %v", err)
	}
	return err
}

// Delete 删除探测爆款
func (dao *ExploreHotBurstDAO) Delete(shopID string) error {
	filter := bson.M{"shop_id": shopID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete explore hot burst error: %v", err)
	}
	return err
}

// ListAll 获取所有探测爆款（带分页）
func (dao *ExploreHotBurstDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count explore hot bursts error: %v", err)
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
		log.Printf("List explore hot bursts error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	exploreHotBursts := make([]ExploreHotBurst, 0)
	if err = cursor.All(ctx, &exploreHotBursts); err != nil {
		return result, err
	}
	fmt.Println(exploreHotBursts)
	result["list"] = exploreHotBursts
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
