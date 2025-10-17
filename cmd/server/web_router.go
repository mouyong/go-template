package server

import (
	"embed"
	"net/http"
	"strings"

	"app/internal/common"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// SetWebRouter 配置前端静态文件服务
func SetWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	// 使用 embed 文件系统服务静态文件
	router.Use(static.Serve("/", common.EmbedFolder(buildFS, "build")))

	// 处理 SPA 路由
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.RequestURI

		// API 路由交给后端处理
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"err_code": 404,
				"err_msg":  "接口不存在",
			})
			return
		}

		// 静态资源文件
		if strings.HasPrefix(path, "/static") ||
			strings.HasPrefix(path, "/assets") ||
			strings.Contains(path, ".") {
			c.Status(http.StatusNotFound)
			return
		}

		// SPA 路由返回 index.html
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
	})
}
