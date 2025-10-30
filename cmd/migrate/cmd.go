package migrate

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	_ "app/db/migrations" // 导入迁移文件
	"app/internal/initialization"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "执行所有未执行的数据库迁移",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigration("up"); err != nil {
			log.Fatalf("迁移失败: %v", err)
		}
		fmt.Println("✅ 迁移完成")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "回滚最后一次数据库迁移",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigration("down"); err != nil {
			log.Fatalf("回滚失败: %v", err)
		}
		fmt.Println("✅ 回滚完成")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看数据库迁移状态",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigration("status"); err != nil {
			log.Fatalf("查询状态失败: %v", err)
		}
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置所有数据库迁移",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigration("reset"); err != nil {
			log.Fatalf("重置失败: %v", err)
		}
		fmt.Println("✅ 重置完成")
	},
}

var cmd = &cobra.Command{
	Use:   "migrate",
	Short: "数据库迁移管理",
}

func Register(rootCmd *cobra.Command) error {
	cmd.AddCommand(upCmd)
	cmd.AddCommand(downCmd)
	cmd.AddCommand(statusCmd)
	cmd.AddCommand(resetCmd)
	rootCmd.AddCommand(cmd)
	return nil
}

func runMigration(command string) error {
	// 获取配置文件路径
	cfg, err := cmd.Flags().GetString("config")
	if err != nil {
		cfg = "./config.yaml"
	}

	// 加载配置
	config := initialization.LoadConfig(cfg)

	// 检查业务数据库配置
	if config.DbHost == "" || config.DbDatabase == "" {
		return fmt.Errorf("业务数据库配置为空, 请检查 DB_* 配置项")
	}

	// 构建 DSN (使用业务数据库配置)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.DbUsername,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbDatabase,
	)

	fmt.Printf("📦 使用业务数据库: %s@%s:%d/%s\n",
		config.DbUsername, config.DbHost, config.DbPort, config.DbDatabase)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	defer db.Close()

	// 设置 goose 使用 MySQL 方言
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("设置数据库方言失败: %w", err)
	}

	// 执行迁移命令
	switch command {
	case "up":
		return goose.Up(db, ".")
	case "down":
		return goose.Down(db, ".")
	case "status":
		return goose.Status(db, ".")
	case "reset":
		return goose.Reset(db, ".")
	default:
		return fmt.Errorf("未知命令: %s", command)
	}
}
