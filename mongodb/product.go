package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Stat 统计信息
type Stat struct {
	Sales       string      `json:"sales" bson:"sales"`
	LiveSales   string      `json:"live_sales" bson:"live_sales"`
	VideoSales  string      `json:"video_sales" bson:"video_sales"`
	OtherSales  string      `json:"other_sales" bson:"other_sales"`
	Views       string      `json:"views" bson:"views"`
	ViewsVIP    string      `json:"views_vip" bson:"views_vip"`
	CVR         float64     `json:"cvr" bson:"cvr"`
	CVRStr      string      `json:"cvr_str" bson:"cvr_str"`
	Lives       interface{} `json:"lives" bson:"lives"`
	Videos      interface{} `json:"videos" bson:"videos"`
	ImageVideos string      `json:"image_videos" bson:"image_videos"`
	Users       string      `json:"users" bson:"users"`
	MarketType  MarketType  `json:"market_type" bson:"market_type"`
}

// Product 商品实体
type Product struct {
	ProductID       string   `json:"product_id" bson:"product_id"` // 使用product_id作为主键
	Cover           string   `json:"cover" bson:"cover"`
	Title           string   `json:"title" bson:"title"`
	Price           int      `json:"price" bson:"price"` // 单位分
	PriceStr        string   `json:"price_str" bson:"price_str"`
	ProductURL      string   `json:"product_url" bson:"product_url"`
	Source          int      `json:"source" bson:"source"`
	CosRatio        float64  `json:"cos_ratio" bson:"cos_ratio"`
	IsHighCosRatio  int      `json:"is_high_cos_ratio" bson:"is_high_cos_ratio"`
	ColonelCosRatio float64  `json:"colonel_cos_ratio" bson:"colonel_cos_ratio"`
	ColonelCosFee   int      `json:"colonel_cos_fee" bson:"colonel_cos_fee"`
	Gain            float64  `json:"gain" bson:"gain"`
	GainStr         string   `json:"gain_str" bson:"gain_str"`
	ActivityURL     string   `json:"activity_url" bson:"activity_url"`
	Stat            Stat     `json:"stat" bson:"stat"`
	H2Sales         *string  `json:"h2_sales,omitempty" bson:"h2_sales,omitempty"`
	TodaySales      *string  `json:"today_sales,omitempty" bson:"today_sales,omitempty"`
	SalesList       *string  `json:"sales_list,omitempty" bson:"sales_list,omitempty"` // 根据实际数据类型调整
	GoodRatio       *float64 `json:"good_ratio,omitempty" bson:"good_ratio,omitempty"`
	FindTime        int64    `json:"find_time" bson:"find_time"`     // Unix时间戳
	UpdateTime      int64    `json:"update_time" bson:"update_time"` // Unix时间戳
	IsSoldOut       bool     `json:"is_sold_out" bson:"is_sold_out"`
	IsFollowed      bool     `json:"is_followed" bson:"is_followed"`
	Day30Sales      string   `json:"day30_sales" bson:"day30_sales"`
	Day15Sales      string   `json:"day15_sales" bson:"day15_sales"`
	Day7Sales       string   `json:"day7_sales" bson:"day7_sales"`
	Day3Sales       string   `json:"day3_sales" bson:"day3_sales"`
	Day1Sales       string   `json:"day1_sales" bson:"day1_sales"`
	ShopID          string   `json:"shop_id" bson:"shop_id"`
	ShopName        string   `json:"shop_name" bson:"shop_name"`
	ShopCover       string   `json:"shop_cover" bson:"shop_cover"`
	DSR             float64  `json:"dsr" bson:"dsr"`
	DSRStr          string   `json:"dsr_str" bson:"dsr_str"`
	SynthesisScore  float64  `json:"synthesis_score" bson:"synthesis_score"`
	ServiceTags     []string `json:"service_tags" bson:"service_tags"`
}

// ProductDAO 商品数据访问对象
type ProductDAO struct {
	collection *mongo.Collection
}

// NewProductDAO 创建新的ProductDAO实例
func NewProductDAO(db *mongo.Database) *ProductDAO {
	return &ProductDAO{
		collection: db.Collection("products"),
	}
}

// Create 创建商品记录
func (dao *ProductDAO) Create(ctx context.Context, product *Product) error {
	_, err := dao.collection.InsertOne(ctx, product)
	return err
}

func (dao *ProductDAO) BatchCreate(ctx context.Context, products []interface{}) error {
	/*_, err := dao.collection.InsertMany(ctx, products)
	if err != nil {
		log.Printf("Create products error: %v", err)
	}
	return err*/
	var models []mongo.WriteModel
	for _, item := range products {
		product := item.(*Product)
		filter := bson.M{"product_id": product.ProductID} // 假设Brand有ID字段
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(product).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite brands error: %v", err)
	}
	return err
}

// GetByProductID 根据商品ID查询
func (dao *ProductDAO) GetByProductID(ctx context.Context, productID string) (*Product, error) {
	var product Product
	err := dao.collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &product, err
}

// Update 更新商品信息
func (dao *ProductDAO) Update(ctx context.Context, productID string, updateData bson.M) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": productID},
		bson.M{"$set": updateData},
		options.Update().SetUpsert(false),
	)
	return err
}

// Delete 删除商品记录
func (dao *ProductDAO) Delete(ctx context.Context, productID string) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": productID})
	return err
}

func (dao *ProductDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
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

	products := make([]Product, 0)
	if err = cursor.All(ctx, &products); err != nil {
		return result, err
	}
	result["list"] = products
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
