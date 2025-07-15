package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SoaringHotSpot struct {
	Rank       int          `json:"rank" bson:"rank"`
	Name       string       `json:"name" bson:"name"`
	Index      string       `json:"index" bson:"index"`
	Flag       int          `json:"flag" bson:"flag"`
	Trend      []TrendPoint `json:"trend" bson:"trend"`
	CreateTime int64        `json:"create_time" bson:"create_time"`
}

type SoaringHotSpotDAO struct {
	collection *mongo.Collection
}

// NewSoaringHotSpotDAO 创建SoaringHotSpot数据访问对象
func NewSoaringHotSpotDAO(db *mongo.Database) *SoaringHotSpotDAO {
	return &SoaringHotSpotDAO{
		collection: db.Collection("soaring_hot_spots"), // 集合名
	}
}

// Create 创建飙升热点
func (dao *SoaringHotSpotDAO) Create(soaringHotSpot *SoaringHotSpot) error {
	_, err := dao.collection.InsertOne(context.TODO(), soaringHotSpot)
	if err != nil {
		log.Printf("Create soaring hot spot error: %v", err)
	}
	return err
}

// BatchCreate 批量创建
func (dao *SoaringHotSpotDAO) BatchCreate(ctx context.Context, soaringHotSpots []interface{}) error {
	var models []mongo.WriteModel
	for _, item := range soaringHotSpots {
		soaringHotSpot := item.(*SoaringHotSpot)
		filter := bson.M{"name": soaringHotSpot.Name, "create_time": soaringHotSpot.CreateTime}
		model := mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(soaringHotSpot).
			SetUpsert(true)
		models = append(models, model)
	}

	_, err := dao.collection.BulkWrite(ctx, models)
	if err != nil {
		log.Printf("BulkWrite soaring hot spots error: %v", err)
	}
	return err
}

// GetByID 根据名称和时间获取飙升热点
func (dao *SoaringHotSpotDAO) GetByID(ctx context.Context, name string, createTime int64) (*SoaringHotSpot, error) {
	var soaringHotSpot SoaringHotSpot
	filter := bson.M{"name": name, "create_time": createTime}
	err := dao.collection.FindOne(ctx, filter).Decode(&soaringHotSpot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		log.Printf("Get soaring hot spot error: %v", err)
		return nil, err
	}
	return &soaringHotSpot, nil
}

// Update 更新飙升热点信息
func (dao *SoaringHotSpotDAO) Update(name string, createTime int64, updateData *SoaringHotSpot) error {
	bsonData, err := bson.Marshal(updateData)
	if err != nil {
		log.Printf("Marshal soaring hot spot error: %v", err)
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
		log.Printf("Update soaring hot spot error: %v", err)
	}
	return err
}

// Delete 删除飙升热点
func (dao *SoaringHotSpotDAO) Delete(name string, createTime int64) error {
	filter := bson.M{"name": name, "create_time": createTime}
	_, err := dao.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Delete soaring hot spot error: %v", err)
	}
	return err
}

// ListAll 获取所有飙升热点（带分页）
func (dao *SoaringHotSpotDAO) ListAll(ctx context.Context, filter bson.M, page, limit int64) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	// 获取总条数
	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Count soaring hot spots error: %v", err)
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
		log.Printf("List soaring hot spots error: %v", err)
		return result, err
	}
	defer cursor.Close(context.TODO())

	soaringHotSpots := make([]SoaringHotSpot, 0)
	if err = cursor.All(ctx, &soaringHotSpots); err != nil {
		return result, err
	}
	fmt.Println(soaringHotSpots)
	result["list"] = soaringHotSpots
	result["page"] = page
	result["limit"] = limit
	return result, nil
}
