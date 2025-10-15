package server

import (
	"embed"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"go-api-template/internal/handlers"
	"go-api-template/internal/initialization"
	"go-api-template/pkg/rabbitmq"
)

var (
	buildFS   embed.FS
	indexPage []byte
)

var cmd = &cobra.Command{
	Use:   "server",
	Short: "run api server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Println(err)
			return
		}

		config := initialization.LoadConfig(cfg)

		fmt.Println("\n正在初始化...")

		// 可选初始化数据库
		if err := initialization.InitDatabaseConnection(); err != nil {
			fmt.Printf("⚠️  数据库: %v\n", err)
		}

		// 可选初始化 RabbitMQ
		if err := rabbitmq.NewRabbitmq(initialization.AppConfig.MqHost, initialization.AppConfig.MqPort); err != nil {
			fmt.Printf("⚠️  RabbitMQ: %v\n", err)
		}
		rabbitmq.ListenQueue()

		fmt.Println("✅ 初始化完成")

		r := gin.Default()

		// 配置可信代理
		if len(config.TrustedProxies) > 0 {
			r.SetTrustedProxies(config.TrustedProxies)
		} else {
			r.SetTrustedProxies(nil)
		}

		// 存储应用配置
		r.Use(func(c *gin.Context) {
			c.Keys = make(map[string]any)
			c.Keys["config"] = config
			c.Next()
		})

		// 基础路由
		r.GET("/api/health", handlers.Health)
		r.GET("/api/hello", handlers.Hello)
		r.POST("/api/echo", handlers.Echo)

		// 配置前端静态文件服务
		SetWebRouter(r, buildFS, indexPage)

		err = startHTTPServer(config, r)
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	},
}

func startHTTPServer(config initialization.Config, r *gin.Engine) (err error) {
	addr := fmt.Sprintf("%s:%d", config.HttpHost, config.HttpPort)

	fmt.Println("\n========================================")
	fmt.Printf("🚀 Server is running!\n\n")
	fmt.Printf("➜ Local:   http://localhost:%d/\n", config.HttpPort)
	fmt.Printf("➜ Network: http://127.0.0.1:%d/\n", config.HttpPort)
	fmt.Println("========================================\n")

	err = r.Run(addr)
	if err != nil {
		return fmt.Errorf("failed to start http server: %v", err)
	}
	return nil
}

func Register(rootCmd *cobra.Command, fs embed.FS, index []byte) error {
	buildFS = fs
	indexPage = index
	rootCmd.AddCommand(cmd)
	return nil
}
