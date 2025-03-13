package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回一个Gin路由器
func SetupRouter() *gin.Engine {
	// 初始化 Gin 引擎
	r := gin.Default()

	// 使用gin-contrib/cors库配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                                                                      // 允许所有源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},                                                       // 允许所有常用方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept"}, // 允许所有常用头
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 允许携带凭证
		MaxAge:           12 * time.Hour, // 预检请求结果缓存时间
	}))

	// 注册所有路由
	RegisterRoutes(r)

	return r
}

// GetRouter 获取配置好的Gin路由器（用于测试）
func GetRouter() *gin.Engine {
	return SetupRouter()
}
