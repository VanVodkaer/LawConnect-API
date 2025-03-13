package router

import (
	"github.com/VanVodkaer/LawConnect-API/internal/handler"
	"github.com/VanVodkaer/LawConnect-API/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RouteGroups 保存所有的路由组
type RouteGroups struct {
	Public *gin.RouterGroup // 公共路由组
	Auth   *gin.RouterGroup // 认证路由组
	API    *gin.RouterGroup // API路由组
	Admin  *gin.RouterGroup // 管理员路由组
}

// 全局变量，保存所有路由组的引用
var Groups RouteGroups

// InitGroups 初始化所有路由组
func InitGroups(r *gin.Engine) {
	// 初始化路由组
	Groups.Public = r.Group("/public")
	Groups.Auth = r.Group("/auth")
	Groups.API = r.Group("/api")
	Groups.API.Use(middleware.JWTAuth()) // API路由组需要JWT验证
	Groups.Admin = Groups.API.Group("/admin")
	Groups.Admin.Use(middleware.AdminRequired()) // 管理员路由组需要管理员权限
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 初始化路由组
	InitGroups(r)

	// 注册各组路由
	registerPublicRoutes()
	registerAuthRoutes()
	registerAPIRoutes()
	registerAdminRoutes()
}

// registerPublicRoutes 注册公共路由
func registerPublicRoutes() {
	// 示例路由，取消注释即可启用
	// Groups.Public.GET("/article", handler.GetArticles)
	// Groups.Public.GET("/article/:id", handler.GetArticleByID)
}

// registerAuthRoutes 注册认证相关路由
func registerAuthRoutes() {
	Groups.Auth.POST("/login", handler.Login)
	Groups.Auth.POST("/register", handler.Register)
}

// registerAPIRoutes 注册需要认证的API路由
func registerAPIRoutes() {
	// 示例路由，取消注释即可启用
	// Groups.API.GET("/user/profile", handler.GetUserProfile)

	// 刷新令牌路由
	Groups.API.POST("/refresh-token", middleware.RefreshToken)
}

// registerAdminRoutes 注册管理员路由
func registerAdminRoutes() {
	// 示例路由，取消注释即可启用
	// Groups.Admin.GET("/users", handler.GetAllUsers)
}
