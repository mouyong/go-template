package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"go-api-template/internal/handlers"
	"go-api-template/internal/initialization"
	"go-api-template/pkg/rabbitmq"
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

		fmt.Println("\n正在初始化")

		// 可选初始化数据库
		if err := initialization.InitDatabaseConnection(); err != nil {
			log.Printf("数据库初始化失败: %v\n", err)
		}

		// 可选初始化 RabbitMQ
		if err := rabbitmq.NewRabbitmq(initialization.AppConfig.MqHost, initialization.AppConfig.MqPort); err != nil {
			log.Printf("RabbitMQ 初始化失败: %v\n", err)
		}
		rabbitmq.ListenQueue()

		fmt.Println("初始化完成")

		r := gin.Default()
		// 存储应用配置
		r.Use(func(c *gin.Context) {
			c.Keys = make(map[string]any)
			c.Keys["config"] = config
			c.Next()
		})

		// 基础路由
		r.GET("/", handlers.Health)
		r.GET("/api/hello", handlers.Hello)
		r.POST("/api/echo", handlers.Echo)

		err = startHTTPServer(config, r)
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	},
}

func startHTTPServer(config initialization.Config, r *gin.Engine) (err error) {
	fmt.Println("")
	log.Printf("http server started: http://localhost:%d/\n", config.HttpPort)
	err = r.Run(fmt.Sprintf("%s:%d", config.HttpHost, config.HttpPort))
	if err != nil {
		return fmt.Errorf("failed to start http server: %v", err)
	}
	return nil
}

func Register(rootCmd *cobra.Command) error {
	rootCmd.AddCommand(cmd)
	return nil
}
