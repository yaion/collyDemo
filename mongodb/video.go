package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// VideoAnchor 视频锚点信息
type VideoAnchor struct {
	AnchorType  int     `json:"anchor_type" bson:"anchor_type"`
	AnchorName  string  `json:"anchor_name" bson:"anchor_name"`
	AnchorTitle *string `json:"anchor_title,omitempty" bson:"anchor_title,omitempty"`
}

// Video 视频实体
type Video struct {
	DateCode         int         `json:"date_code" bson:"date_code"`
	AwemeID          string      `json:"aweme_id" bson:"aweme_id"` // 使用aweme_id作为主键
	Uid              string      `json:"uid" bson:"uid"`
	MarketLevel      int         `json:"market_level" bson:"market_level"`
	NickName         string      `json:"nick_name" bson:"nick_name"`
	Desc             string      `json:"desc" bson:"desc"`
	Cover            string      `json:"cover" bson:"cover"`
	PubTime          int64       `json:"pub_time" bson:"pub_time"`     // Unix时间戳
	LikeCount        string      `json:"like_count" bson:"like_count"` // 含单位字符串
	PlayCount        *string     `json:"play_count,omitempty" bson:"play_count,omitempty"`
	CollectCount     string      `json:"collect_count" bson:"collect_count"`
	InteractiveRatio string      `json:"interactive_ratio" bson:"interactive_ratio"`
	CommentCount     string      `json:"comment_count" bson:"comment_count"`
	ShareCount       string      `json:"share_count" bson:"share_count"`
	URL              string      `json:"url" bson:"url"`
	FansCount        string      `json:"fans_count" bson:"fans_count"`
	Avatar           string      `json:"avatar" bson:"avatar"`
	IsSentence       bool        `json:"is_sentence" bson:"is_sentence"`
	Sentence         *string     `json:"sentence,omitempty" bson:"sentence,omitempty"`
	ProductID        *string     `json:"product_id,omitempty" bson:"product_id,omitempty"`
	ProductTitle     *string     `json:"product_title,omitempty" bson:"product_title,omitempty"`
	ProductCover     *string     `json:"product_cover,omitempty" bson:"product_cover,omitempty"`
	ProductPrice     int         `json:"product_price" bson:"product_price"`
	PriceStr         string      `json:"price_str" bson:"price_str"`
	Sales            string      `json:"sales" bson:"sales"`
	GMV              string      `json:"gmv" bson:"gmv"`
	VideoAnchor      VideoAnchor `json:"video_anchor" bson:"video_anchor"`
}

// VideoDAO 视频数据访问对象
type VideoDAO struct {
	collection *mongo.Collection
}

// NewVideoDAO 创建新的VideoDAO实例
func NewVideoDAO(db *mongo.Database) *VideoDAO {
	return &VideoDAO{
		collection: db.Collection("videos"),
	}
}

// Create 创建视频记录
func (dao *VideoDAO) Create(ctx context.Context, video *Video) error {
	_, err := dao.collection.InsertOne(ctx, video)
	return err
}

// BatchCreate 批量创建视频记录
func (dao *VideoDAO) BatchCreate(ctx context.Context, videos []interface{}) error {
	/*_, err := dao.collection.InsertMany(ctx, videos)
	if err != nil {
		log.Printf("Create videos error: %v", err)
	}
	return nil*/
	var models []mongo.WriteModel
	for _, item := range videos {
		video := item.(*Video)
		filter := bson.M{"aweme_id": video.AwemeID} // 假设Video有ID字段
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(video).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite videos error: %v", err)
	}
	return err
}

// GetByAwemeID 根据视频ID查询
func (dao *VideoDAO) GetByAwemeID(ctx context.Context, awemeID string) (*Video, error) {
	var video Video
	err := dao.collection.FindOne(ctx, bson.M{"_id": awemeID}).Decode(&video)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &video, err
}

// Update 更新视频信息
func (dao *VideoDAO) Update(ctx context.Context, awemeID string, updateData bson.M) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": awemeID},
		bson.M{"$set": updateData},
	)
	return err
}

// Delete 删除视频记录
func (dao *VideoDAO) Delete(ctx context.Context, awemeID string) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": awemeID})
	return err
}

// 分页查询示例
func (dao *VideoDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
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
	defer cursor.Close(ctx)

	videos := make([]Video, 0)
	if err = cursor.All(ctx, &videos); err != nil {
		return result, err
	}
	result["list"] = videos
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
