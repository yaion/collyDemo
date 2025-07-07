package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// User 用户实体
type AuthorInfo struct {
	Uid                string           `json:"uid" bson:"uid"` // 使用uid作为主键
	LinkURL            string           `json:"link_url" bson:"link_url"`
	ShareURL           string           `json:"share_url" bson:"share_url"`
	NickName           string           `json:"nick_name" bson:"nick_name"`
	IsLiving           bool             `json:"is_living" bson:"is_living"`
	Avatar             string           `json:"avatar" bson:"avatar"`
	MarketLevel        int              `json:"market_level" bson:"market_level"`
	DisplayID          string           `json:"display_id" bson:"display_id"`
	Verify             string           `json:"verify" bson:"verify"`
	VerifyType         int              `json:"verify_type" bson:"verify_type"`
	Age                int              `json:"age" bson:"age"`
	Gender             int              `json:"gender" bson:"gender"`
	Location           *string          `json:"location" bson:"location,omitempty"`
	Tag                string           `json:"tag" bson:"tag"`
	Fans               string           `json:"fans" bson:"fans"`
	DSRScore           string           `json:"dsr_score" bson:"dsr_score"`
	DSRLevel           string           `json:"dsr_level" bson:"dsr_level"`
	FollowerOrderRatio float64          `json:"follower_order_ratio" bson:"follower_order_ratio"`
	MCN                *string          `json:"mcn" bson:"mcn,omitempty"`
	Ranking            string           `json:"ranking" bson:"ranking"`
	IsFollowed         bool             `json:"is_followed" bson:"is_followed"`
	RoomID             *string          `json:"room_id" bson:"room_id,omitempty"`
	RoomDateCode       interface{}      `json:"room_date_code" bson:"room_date_code,omitempty"`
	MarketCategories   []MarketCategory `json:"market_categories" bson:"market_categories"`
	MajorCategories    []string         `json:"major_categories" bson:"major_categories"`
	Signature          string           `json:"signature" bson:"signature"`
	GMV                string           `json:"gmv" bson:"gmv"`
	MarketType         MarketType       `json:"market_type" bson:"market_type"`
	MarketDays         int              `json:"market_days" bson:"market_days"`
	SKUs               int              `json:"skus" bson:"skus"`
	ShopID             *string          `json:"shop_id" bson:"shop_id,omitempty"`
	ShopName           string           `json:"shop_name" bson:"shop_name"`
	IsWatch            bool             `json:"is_watch" bson:"is_watch"`
	UpdateTime         int64            `json:"update_time" bson:"update_time"`
	VideoArmsID        string           `json:"video_arms_id" bson:"video_arms_id"`
	VideoArmsType      int              `json:"video_arms_type" bson:"video_arms_type"`
	LiveArmsID         int              `json:"live_arms_id" bson:"live_arms_id"`
	LiveArmsType       int              `json:"live_arms_type" bson:"live_arms_type"`
	IsMaterial         bool             `json:"is_material" bson:"is_material"`
}

// MarketCategory 市场分类
type MarketCategory struct {
	Name    string  `json:"name" bson:"name"`
	Percent float64 `json:"percent" bson:"percent"`
}

// AuthorInfoDao 用户数据访问对象
type AuthorInfoDao struct {
	collection *mongo.Collection
}

// NewAuthorInfo 创建新的UserDAO实例
func NewAuthorInfo(db *mongo.Database) *AuthorInfoDao {
	return &AuthorInfoDao{
		collection: db.Collection("authorInfo"),
	}
}

// Create 创建用户记录
func (dao *AuthorInfoDao) Create(ctx context.Context, authorInfo *AuthorInfo) error {
	_, err := dao.collection.InsertOne(ctx, authorInfo)
	return err
}

// Create 批量创建
func (dao *AuthorInfoDao) BatchCreate(ctx context.Context, authorInfos []interface{}) error {
	_, err := dao.collection.InsertMany(ctx, authorInfos)
	if err != nil {
		log.Printf("Create authorInfos error: %v", err)
	}
	return err
}

// GetByID 根据用户ID查询
func (dao *AuthorInfoDao) GetByID(ctx context.Context, userID string) (*AuthorInfo, error) {
	var authorInfo AuthorInfo
	err := dao.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&authorInfo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &authorInfo, err
}

// Update 更新用户信息
func (dao *AuthorInfoDao) Update(ctx context.Context, userID string, updateData bson.M) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": updateData},
		options.Update().SetUpsert(false),
	)
	return err
}

// Delete 删除用户记录
func (dao *AuthorInfoDao) Delete(ctx context.Context, userID string) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": userID})
	return err
}
