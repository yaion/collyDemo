package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Author struct {
	Aup                string     `json:"aup" bson:"aup"`
	Avatar             string     `json:"avatar" bson:"avatar"`
	AvgLikeCount       string     `json:"avg_like_count" bson:"avg_like_count"`
	AvgLiveGmv         string     `json:"avg_live_gmv" bson:"avg_live_gmv"`
	AvgLiveSales       string     `json:"avg_live_sales" bson:"avg_live_sales"`
	AvgPlayCount       string     `json:"avg_play_count" bson:"avg_play_count"`
	AvgStayDuration    int        `json:"avg_stay_duration" bson:"avg_stay_duration"`
	AvgThrough         float64    `json:"avg_through" bson:"avg_through"`
	AvgTotalUsers      string     `json:"avg_total_users" bson:"avg_total_users"`
	AvgUsers           string     `json:"avg_users" bson:"avg_users"`
	AvgVideoDuration   int        `json:"avg_video_duration" bson:"avg_video_duration"`
	AvgVideoGmv        string     `json:"avg_video_gmv" bson:"avg_video_gmv"`
	AvgVideoSales      string     `json:"avg_video_sales" bson:"avg_video_sales"`
	CateAvgLiveGmv     string     `json:"cate_avg_live_gmv" bson:"cate_avg_live_gmv"`
	CateAvgLiveSales   string     `json:"cate_avg_live_sales" bson:"cate_avg_live_sales"`
	CateAvgVideoGmv    string     `json:"cate_avg_video_gmv" bson:"cate_avg_video_gmv"`
	CateAvgVideoSales  string     `json:"cate_avg_video_sales" bson:"cate_avg_video_sales"`
	CateGmv            string     `json:"cate_gmv" bson:"cate_gmv"`
	CateLiveGmv        string     `json:"cate_live_gmv" bson:"cate_live_gmv"`
	CateLiveSales      string     `json:"cate_live_sales" bson:"cate_live_sales"`
	CateSales          string     `json:"cate_sales" bson:"cate_sales"`
	CateVideoGmv       string     `json:"cate_video_gmv" bson:"cate_video_gmv"`
	CateVideoSales     string     `json:"cate_video_sales" bson:"cate_video_sales"`
	Cvr                float64    `json:"cvr" bson:"cvr"`
	CvrStr             string     `json:"cvr_str" bson:"cvr_str"`
	DisplayID          string     `json:"display_id" bson:"display_id"`
	Fans               string     `json:"fans" bson:"fans"`
	FollowID           int        `json:"follow_id" bson:"follow_id"`
	Gender             int        `json:"gender" bson:"gender"`
	Gmv                string     `json:"gmv" bson:"gmv"`
	IncFans            string     `json:"inc_fans" bson:"inc_fans"`
	IsFollowed         bool       `json:"is_followed" bson:"is_followed"`
	IsInclude          bool       `json:"is_include" bson:"is_include"`
	IsLiving           bool       `json:"is_living" bson:"is_living"`
	IsShop             bool       `json:"is_shop" bson:"is_shop"`
	LiveAup            string     `json:"live_aup" bson:"live_aup"`
	LiveGmv            string     `json:"live_gmv" bson:"live_gmv"`
	LiveInteractRatio  float64    `json:"live_interact_ratio" bson:"live_interact_ratio"`
	LiveSales          string     `json:"live_sales" bson:"live_sales"`
	Lives              int        `json:"lives" bson:"lives"`
	MarketLevel        int        `json:"market_level" bson:"market_level"`
	MarketType         MarketType `json:"market_type" bson:"market_type"`
	NickName           string     `json:"nick_name" bson:"nick_name"`
	Rpm                string     `json:"rpm" bson:"rpm"`
	SkuVideos          int        `json:"sku_videos" bson:"sku_videos"`
	Skus               int        `json:"skus" bson:"skus"`
	UID                string     `json:"uid" bson:"uid"`
	UnInclude          bool       `json:"un_include" bson:"un_include"`
	UV                 string     `json:"uv" bson:"uv"`
	VideoAup           string     `json:"video_aup" bson:"video_aup"`
	VideoGmv           string     `json:"video_gmv" bson:"video_gmv"`
	VideoInteractRatio float64    `json:"video_interact_ratio" bson:"video_interact_ratio"`
	VideoSales         string     `json:"video_sales" bson:"video_sales"`
	VideoSkus          int        `json:"video_skus" bson:"video_skus"`
	Videos             int        `json:"videos" bson:"videos"`
}

type MarketType struct {
	IsPure     int     `json:"is_pure" bson:"is_pure"`
	LiveRatio  float64 `json:"live_ratio" bson:"live_ratio"`
	MarketType int     `json:"market_type" bson:"market_type"`
	OtherRatio float64 `json:"other_ratio" bson:"other_ratio"`
	Ratio      float64 `json:"ratio" bson:"ratio"`
	VideoRatio float64 `json:"video_ratio" bson:"video_ratio"`
}

type AuthorDAO struct {
	collection *mongo.Collection
}

// NewAuthorDAO 创建Author数据访问对象
func NewAuthorDAO(db *mongo.Database) *AuthorDAO {
	return &AuthorDAO{
		collection: db.Collection("authors"), // 集合名
	}
}

// Create 创建作者
func (dao *AuthorDAO) Create(author *Author) error {
	_, err := dao.collection.InsertOne(context.TODO(), author)
	if err != nil {
		log.Printf("Create author error: %v", err)
	}
	return err
}

// Create 批量创建
func (dao *AuthorDAO) BatchCreate(ctx context.Context, authors []interface{}) error {
	/*_, err := dao.collection.InsertMany(ctx, authors)
	if err != nil {
		log.Printf("Create author error: %v", err)
	}
	return err*/
	var models []mongo.WriteModel
	for _, item := range authors {
		author := item.(*Author)
		filter := bson.M{"uid": author.UID}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(author).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite authorInfos error: %v", err)
	}
	return err
}

// GetByID 根据ID获取作者
func (dao *AuthorDAO) GetByID(ctx context.Context, uid string) (*Author, error) {
	var author Author
	filter := bson.M{"uid": uid}
	err := dao.collection.FindOne(ctx, filter).Decode(&author)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get author error: %v", err)
		return nil, err
	}
	return &author, nil
}

// Update 更新作者信息
func (dao *AuthorDAO) Update(authorID string, updateData *Author) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal author error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"uid": authorID} // 注意：根据结构体应使用 uid 而非 author_id
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update author error: %v", err)
	}
	return err
}

// Delete 删除作者
func (dao *AuthorDAO) Delete(authorID string) error {
	filter := bson.M{"author_id": authorID}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete author error: %v", err)
	}
	return err
}

// ListAll 获取所有作者（带分页）
func (dao *AuthorDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
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

	authors := make([]Author, 0)
	if err = cursor.All(ctx, &authors); err != nil {
		return result, err
	}
	fmt.Println(authors)
	result["list"] = authors
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
