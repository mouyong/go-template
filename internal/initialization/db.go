package initialization

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitDatabaseConnection() error {
	// 检查数据库配置是否为空
	if AppConfig.DbHost == "" || AppConfig.DbDatabase == "" {
		fmt.Println("⏭️  数据库配置为空，跳过初始化")
		return nil
	}

	fmt.Printf("正在初始化 %s 数据库连接...\n", AppConfig.DbType)

	switch AppConfig.DbType {
	case "mysql":
		return initMySQL()
	case "postgres", "postgresql":
		return initPostgreSQL()
	case "mongodb", "mongo":
		return initMongoDB()
	default:
		return fmt.Errorf("不支持的数据库类型: %s (支持: mysql, postgres, mongodb)", AppConfig.DbType)
	}
}

func initMySQL() error {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DbUsername,
		AppConfig.DbPassword,
		AppConfig.DbHost,
		AppConfig.DbPort,
		AppConfig.DbDatabase,
	)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("MySQL 连接失败: %v", err)
	}

	fmt.Println("✅ MySQL 数据库连接成功")
	return nil
}

func initPostgreSQL() error {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		AppConfig.DbHost,
		AppConfig.DbUsername,
		AppConfig.DbPassword,
		AppConfig.DbDatabase,
		AppConfig.DbPort,
	)

	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("PostgreSQL 连接失败: %v", err)
	}

	fmt.Println("✅ PostgreSQL 数据库连接成功")
	return nil
}

func initMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建 MongoDB 连接字符串
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		AppConfig.DbUsername,
		AppConfig.DbPassword,
		AppConfig.DbHost,
		AppConfig.DbPort,
	)

	// 如果没有用户名密码，使用简单连接字符串
	if AppConfig.DbUsername == "" {
		uri = fmt.Sprintf("mongodb://%s:%d", AppConfig.DbHost, AppConfig.DbPort)
	}

	clientOptions := options.Client().ApplyURI(uri)

	var err error
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("MongoDB 连接失败: %v", err)
	}

	// 测试连接
	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("MongoDB Ping 失败: %v", err)
	}

	// 设置数据库
	MongoDB = MongoClient.Database(AppConfig.DbDatabase)

	fmt.Println("✅ MongoDB 数据库连接成功")
	return nil
}