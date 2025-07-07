package mongodb

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Brand 结构体定义
type Brand struct {
	BrandID    string      `json:"brand_id" bson:"brand_id"` // 使用brand_id作为主键
	Name       string      `json:"name" bson:"name"`
	IsBlack    bool        `json:"is_black" bson:"is_black"`
	ShopCount  string      `json:"shop_count" bson:"shop_count"`
	Logo       string      `json:"logo" bson:"logo"`
	RankInfo   interface{} `json:"rank_info" bson:"rank_info"`
	Stat       BrandStat   `json:"stat" bson:"stat"`
	IsFollowed bool        `json:"is_followed" bson:"is_followed"`
	FollowID   int         `json:"follow_id" bson:"follow_id"`
	CateName   string      `json:"cate_name" bson:"cate_name"`
}

// BrandStat 品牌统计信息
type BrandStat struct {
	SOVRankInfo      interface{} `json:"sov_rank_info" bson:"sov_rank_info"`
	BrandExpose      interface{} `json:"brand_expose" bson:"brand_expose"`
	BrandInteraction interface{} `json:"brand_interaction" bson:"brand_interaction"`
	SalesRankInfo    interface{} `json:"sales_rank_info" bson:"sales_rank_info"`
	Sales            string      `json:"sales" bson:"sales"`
	GMV              string      `json:"gmv" bson:"gmv"`
	LiveSales        string      `json:"live_sales" bson:"live_sales"`
	LiveGMV          string      `json:"live_gmv" bson:"live_gmv"`
	VideoSales       string      `json:"video_sales" bson:"video_sales"`
	VideoGMV         string      `json:"video_gmv" bson:"video_gmv"`
	SKUSaleCount     string      `json:"sku_sale_count" bson:"sku_sale_count"`
	Lives            string      `json:"lives" bson:"lives"`
	Videos           string      `json:"videos" bson:"videos"`
	Users            string      `json:"users" bson:"users"`
	MarketType       MarketType  `json:"market_type" bson:"market_type"`
}

// MarketType 市场类型信息
/*type BrandMarketType struct {
	MarketType int     `json:"market_type" bson:"market_type"`
	Ratio      float64 `json:"ratio" bson:"ratio"`
	LiveRatio  float64 `json:"live_ratio" bson:"live_ratio"`
	VideoRatio float64 `json:"video_ratio" bson:"video_ratio"`
	OtherRatio float64 `json:"other_ratio" bson:"other_ratio"`
	IsPure     int     `json:"is_pure" bson:"is_pure"`
}*/

// BrandDAO 品牌数据访问对象
type BrandDAO struct {
	collection *mongo.Collection
}

// NewBrandDAO 创建新的BrandDAO实例
func NewBrandDAO(db *mongo.Database) *BrandDAO {
	return &BrandDAO{
		collection: db.Collection("brands"),
	}
}

// Create 创建品牌记录
func (dao *BrandDAO) Create(ctx context.Context, brand *Brand) error {
	_, err := dao.collection.InsertOne(ctx, brand)
	return err
}

// Create 批量创建
func (dao *BrandDAO) BatchCreate(ctx context.Context, brands []interface{}) error {
	/*_, err := dao.collection.InsertMany(ctx, brands)
	if err != nil {
		log.Printf("Create brands error: %v", err)
	}
	return err*/
	var models []mongo.WriteModel
	for _, item := range brands {
		brand := item.(*Brand)
		filter := bson.M{"brand_id": brand.BrandID} // 假设Brand有ID字段
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(brand).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite brands error: %v", err)
	}
	return err
}

// GetByID 根据品牌ID查询
func (dao *BrandDAO) GetByID(ctx context.Context, brandID string) (*Brand, error) {
	var brand Brand
	err := dao.collection.FindOne(ctx, bson.M{"_id": brandID}).Decode(&brand)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &brand, err
}

// Update 更新品牌信息
func (dao *BrandDAO) Update(ctx context.Context, brandID string, updateData bson.M) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": brandID},
		bson.M{"$set": updateData},
		options.Update().SetUpsert(false),
	)
	return err
}

// Delete 删除品牌记录
func (dao *BrandDAO) Delete(ctx context.Context, brandID string) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": brandID})
	return err
}

func (dao *BrandDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
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

	// todo 添加默认排序（按粉丝数降序）
	//findOptions.SetSort(bson.D{{Key: "fans", Value: -1}})

	cursor, err := dao.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("List authors error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	brands := make([]Brand, 0)
	if err = cursor.All(ctx, &brands); err != nil {
		return result, err
	}
	result["list"] = brands
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
