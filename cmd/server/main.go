package main

import (
	"fmt"
	"log"

	"github.com/VanVodkaer/LawConnect-API/internal/db"
	"github.com/VanVodkaer/LawConnect-API/internal/router"
	"github.com/VanVodkaer/LawConnect-API/utils/admin"
	"github.com/VanVodkaer/LawConnect-API/utils/config"
)

func main() {
	// 初始化配置
	initConfig()

	// 初始化数据库
	initDB()

	// 创建管理员账户（如果不存在）
	admin.CreateAdminIfNotExists()

	// 初始化并启动服务器
	initServer()
}

// initConfig 加载配置文件
func initConfig() {
	config.LoadConfig("config.yaml")
}

// initDB 初始化数据库连接
func initDB() {
	var database = config.GlobalConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		database.User, database.Password, database.Host, database.Port, database.DBName)

	db.InitDB(dsn)
}

// initServer 初始化并启动服务器
func initServer() {
	var server = config.GlobalConfig.Server

	// 设置路由
	r := router.SetupRouter()

	// 启动服务
	log.Printf("服务器正在启动，监听地址: %s:%d", server.Host, server.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", server.Host, server.Port)); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
