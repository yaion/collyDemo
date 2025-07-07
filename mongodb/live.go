package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Live 直播数据实体
type Live struct {
	RoomID        string      `json:"room_id" bson:"room_id"` // 使用room_id作为主键
	DateCode      int         `json:"date_code" bson:"date_code"`
	IsLive        int         `json:"is_live" bson:"is_live"`
	MarketLevel   int         `json:"market_level" bson:"market_level"`
	UID           string      `json:"uid" bson:"uid"`
	NickName      string      `json:"nick_name" bson:"nick_name"`
	FansCount     string      `json:"fans_count" bson:"fans_count"`
	Avatar        string      `json:"avatar" bson:"avatar"`
	Title         string      `json:"title" bson:"title"`
	LiveCateName  string      `json:"live_cate_name" bson:"live_cate_name"`
	Cover         string      `json:"cover" bson:"cover"`
	PubTime       int64       `json:"pub_time" bson:"pub_time"`
	RPM           string      `json:"rpm" bson:"rpm"`
	GMV           string      `json:"gmv" bson:"gmv"`
	SKUs          string      `json:"skus" bson:"skus"`
	Sales         string      `json:"sales" bson:"sales"`
	StayDuration  int         `json:"stay_duration" bson:"stay_duration"`
	TotalUsers    string      `json:"total_users" bson:"total_users"`
	PeakUsers     string      `json:"peak_users" bson:"peak_users"`
	Duration      int         `json:"duration" bson:"duration"`
	Through       float64     `json:"through" bson:"through"`
	UV            string      `json:"uv" bson:"uv"`
	AUP           string      `json:"aup" bson:"aup"`
	InteractRatio float64     `json:"interact_ratio" bson:"interact_ratio"`
	Traffic       []Traffic   `json:"traffic" bson:"traffic,omitempty"`
	MaxTraffic    *MaxTraffic `json:"max_traffic,omitempty" bson:"max_traffic,omitempty"`
	JumpURL       string      `json:"jump_url" bson:"jump_url"`
	ExposedNum    string      `json:"exposed_num" bson:"exposed_num"`
	IsCart        bool        `json:"is_cart" bson:"is_cart"`
	IsRedirect    bool        `json:"is_redirect" bson:"is_redirect"`
}

// Traffic 流量来源
type Traffic struct {
	Name  string  `json:"name" bson:"name"`
	Value float64 `json:"value" bson:"value"`
}

// MaxTraffic 最大流量来源
type MaxTraffic struct {
	Name  string  `json:"name" bson:"name"`
	Value float64 `json:"value" bson:"value"`
}

// LiveDAO 直播数据访问对象
type LiveDAO struct {
	collection *mongo.Collection
}

// NewLiveDAO 创建新的LiveDAO实例
func NewLiveDAO(db *mongo.Database) *LiveDAO {
	return &LiveDAO{
		collection: db.Collection("live"),
	}
}

// Create 创建直播记录
func (dao *LiveDAO) Create(ctx context.Context, live *Live) error {
	_, err := dao.collection.InsertOne(ctx, live)
	return err
}

func (dao *LiveDAO) BatchCreate(ctx context.Context, lives []interface{}) error {
	/*_, err := dao.collection.InsertMany(ctx, lives)
	if err != nil {
		log.Printf("Create lives error: %v", err)
	}
	return err*/
	var models []mongo.WriteModel
	for _, item := range lives {
		live := item.(*Live)
		filter := bson.M{"room_id": live.RoomID} // 使用RoomID作为主键
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(live).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite lives error: %v", err)
	}
	return err
}

// GetByRoomID 根据房间ID查询
func (dao *LiveDAO) GetByRoomID(ctx context.Context, roomID string) (*Live, error) {
	var live Live
	err := dao.collection.FindOne(ctx, bson.M{"_id": roomID}).Decode(&live)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &live, err
}

// Update 更新直播信息
func (dao *LiveDAO) Update(ctx context.Context, roomID string, updateData bson.M) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": roomID},
		bson.M{"$set": updateData},
		options.Update().SetUpsert(false),
	)
	return err
}

// Delete 删除直播记录
func (dao *LiveDAO) Delete(ctx context.Context, roomID string) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": roomID})
	return err
}

func (dao *LiveDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
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

	lives := make([]Live, 0)
	if err = cursor.All(ctx, &lives); err != nil {
		return result, err
	}
	fmt.Println(lives)
	result["list"] = lives
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
