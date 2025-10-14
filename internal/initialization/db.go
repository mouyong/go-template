package initialization

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDatabaseConnection() error {
	// 检查数据库配置是否为空
	if AppConfig.DbHost == "" || AppConfig.DbDatabase == "" {
		fmt.Println("⏭️  数据库配置为空，跳过初始化")
		return nil
	}

	fmt.Println("正在初始化数据库连接...")

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
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	fmt.Println("✅ 数据库连接成功")
	return nil
}