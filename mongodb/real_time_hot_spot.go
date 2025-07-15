package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrendPoint struct {
	X int64 `json:"x" bson:"x"`
	Y int64 `json:"y" bson:"y"`
}

type RealTimeHotSpot struct {
	Rank       int          `json:"rank" bson:"rank"`
	Name       string       `json:"name" bson:"name"`
	Index      string       `json:"index" bson:"index"`
	Flag       int          `json:"flag" bson:"flag"`
	Trend      []TrendPoint `json:"trend" bson:"trend"`
	CreateTime int64        `json:"create_time" bson:"create_time"`
}

type RealTimeHotSpotDAO struct {
	collection *mongo.Collection
}

// NewRealTimeHotSpotDAO 创建RealTimeHotSpot数据访问对象
func NewRealTimeHotSpotDAO(db *mongo.Database) *RealTimeHotSpotDAO {
	return &RealTimeHotSpotDAO{
		collection: db.Collection("real_time_hot_spots"), // 集合名
	}
}

// Create 创建实时热点
func (dao *RealTimeHotSpotDAO) Create(realTimeHotSpot *RealTimeHotSpot) error {
	_, err := dao.collection.InsertOne(context.TODO(), realTimeHotSpot)
	if err != nil {
		log.Printf("Create real time hot spot error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *RealTimeHotSpotDAO) BatchCreate(ctx context.Context, realTimeHotSpots []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range realTimeHotSpots {
		realTimeHotSpot := item.(*RealTimeHotSpot)
		filter := bson.M{"name": realTimeHotSpot.Name, "create_time": realTimeHotSpot.CreateTime}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(realTimeHotSpot).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite real time hot spots error: %v", err)
	}
	return err
}

// GetByID 根据名称和时间获取实时热点
func (dao *RealTimeHotSpotDAO) GetByID(ctx context.Context, name string, createTime int64) (*RealTimeHotSpot, error) {
	var realTimeHotSpot RealTimeHotSpot
	filter := bson.M{"name": name, "create_time": createTime}
	err := dao.collection.FindOne(ctx, filter).Decode(&realTimeHotSpot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get real time hot spot error: %v", err)
		return nil, err
	}
	return &realTimeHotSpot, nil
}

// Update 更新实时热点信息
func (dao *RealTimeHotSpotDAO) Update(name string, createTime int64, updateData *RealTimeHotSpot) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal real time hot spot error: %v", err)
		return err
	}

	var updateDoc bson.M
	if err = bson.Unmarshal(bsonData, &updateDoc); err != nil {
		log.Printf("Unmarshal to bson.M error: %v", err)
		return err
	}

	filter := bson.M{"name": name, "create_time": createTime}
	update := bson.M{"$set": updateDoc}

	_, err = dao.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Update real time hot spot error: %v", err)
	}
	return err
}

// Delete 删除实时热点
func (dao *RealTimeHotSpotDAO) Delete(name string, createTime int64) error {
	filter := bson.M{"name": name, "create_time": createTime}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete real time hot spot error: %v", err)
	}
	return err
}

// ListAll 获取所有实时热点（带分页）
func (dao *RealTimeHotSpotDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count real time hot spots error: %v", err)
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
		log.Printf("List real time hot spots error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	realTimeHotSpots := make([]RealTimeHotSpot, 0)
	if err = cursor.All(ctx, &realTimeHotSpots); err != nil {
		return result, err
	}
	fmt.Println(realTimeHotSpots)
	result["list"] = realTimeHotSpots
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
