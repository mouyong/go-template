package handlers

import "github.com/gin-gonic/gin"

// Health 健康检查
func Health(c *gin.Context) {
	resp := NewResp(c)
	resp.successWithData(gin.H{
		"status": "ok",
		"service": "app",
	}, nil)
}

// Hello GET 示例
func Hello(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	resp := NewResp(c)
	resp.successWithData(gin.H{
		"message": "Hello, " + name + "!",
	}, nil)
}

// Echo POST 示例
func Echo(c *gin.Context) {
	var json map[string]interface{}
	if err := c.ShouldBindJSON(&json); err != nil {
		resp := NewResp(c)
		resp.fail("Invalid JSON")
		return
	}

	resp := NewResp(c)
	resp.successWithData(gin.H{
		"echo": json,
	}, nil)
}
