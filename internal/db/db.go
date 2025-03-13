package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// 定义全局数据库连接对象
var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) {
	var err error

	// 连接 MySQL 服务器（直接连接数据库，不再创建数据库）
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 测试数据库连接
	if err := DB.Ping(); err != nil {
		log.Fatal("数据库连接不可用:", err)
	}

	log.Println("成功连接到数据库")

	// 运行外部 SQL 脚本（如果需要）
	err = execSchemaSQL("schema.sql")
	if err != nil {
		log.Fatal("执行数据库初始化 SQL 失败:", err)
	}
}

// execSchemaSQL 执行外部 SQL 文件
func execSchemaSQL(filepath string) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("读取 SQL 文件失败: %v", err)
	}

	// 将 SQL 文件内容拆分成多个语句并执行
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("执行 SQL 失败: %v", err)
		}
	}
	log.Println("数据库初始化 SQL 执行完成")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("数据库连接已关闭")
	}
}
