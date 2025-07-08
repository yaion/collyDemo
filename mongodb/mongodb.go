package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var MongoClient *mongo.Client

func InitMongo() error {

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// 连接 MongoDB（带超时控制）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//defer client.Disconnect(ctx) // 确保关闭连接

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping failed: ", err)
		return err
	}
	MongoClient = client
	fmt.Println("Connected to MongoDB!")
	return nil
}

func GetMongo() *mongo.Client {
	return MongoClient
}
