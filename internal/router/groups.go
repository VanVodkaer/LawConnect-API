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
	// 法学交流社区
	Groups.Public.GET("/community/latest", handler.GetCommunityLatest)
	Groups.Public.GET("/community/hottest", handler.GetCommunityHottest)
	Groups.Public.GET("/community/hotqa", handler.GetCommunityHotQA)
	// 政策推送专区
	Groups.Public.GET("/policy/latest", handler.GetPolicyLatest)
	Groups.Public.GET("/policy/local", handler.GetPolicyLocal)
	Groups.Public.GET("/policy/interpretation", handler.GetPolicyInterpretation)
	// 线下实践平台
	Groups.Public.GET("/offline/cooperation", handler.GetOfflineCooperation)
	Groups.Public.GET("/offline/online", handler.GetOfflineOnline)
	Groups.Public.GET("/offline/registration", handler.GetOfflineRegistration)
	// 文章详情路由
	Groups.Public.GET("/article/:id", handler.GetArticleDetail)
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

	// 评论相关路由
	Groups.API.POST("/article/:id/comment", handler.AddComment)

	// 文章点赞相关路由
	Groups.API.POST("/article/:id/like", handler.LikeArticle)         // 点赞
	Groups.API.DELETE("/article/:id/like", handler.UnlikeArticle)     // 取消点赞
	Groups.API.GET("/article/:id/like", handler.GetArticleLikeStatus) // 获取点赞状态

	// 评论点赞相关路由
	Groups.API.POST("/comment/:id/like", handler.LikeComment)         // 点赞
	Groups.API.DELETE("/comment/:id/like", handler.UnlikeComment)     // 取消点赞
	Groups.API.GET("/comment/:id/like", handler.GetCommentLikeStatus) // 获取点赞状态

}

// registerAdminRoutes 注册管理员路由
func registerAdminRoutes() {
	// 示例路由，取消注释即可启用
	// Groups.Admin.GET("/users", handler.GetAllUsers)
}
