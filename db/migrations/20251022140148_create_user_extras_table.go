package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

// User 用户表模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	// 从 sql.Tx 获取底层的 *sql.DB
	// 使用 GORM 连接到同一个数据库
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: tx,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create gorm instance: %w", err)
	}

	// 使用 GORM AutoMigrate 创建表
	if err := db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: tx,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create gorm instance: %w", err)
	}

	// 删除表
	if err := db.Migrator().DropTable(&User{}); err != nil {
		return fmt.Errorf("failed to drop table: %w", err)
	}

	return nil
}
