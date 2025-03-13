package admin

import (
	"fmt"
	"log"

	"github.com/VanVodkaer/LawConnect-API/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdminIfNotExists 如果不存在管理员用户则创建一个
func CreateAdminIfNotExists() {
	var count int
	// 检查是否已有管理员用户
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = ?", db.RoleAdmin).Scan(&count)
	if err != nil {
		log.Printf("检查管理员是否存在失败: %v", err)
		return
	}

	// 如果已有管理员用户则返回
	if count > 0 {
		fmt.Println("管理员用户已存在，无需创建")
		return
	}

	// 设置管理员默认信息
	username := "admin"
	email := "admin@example.com"
	password := "admin123"

	// 检查是否已有同名用户
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if count > 0 {
		// 如果有同名普通用户，将其升级为管理员
		_, err := db.DB.Exec("UPDATE users SET role = ? WHERE username = ?", db.RoleAdmin, username)
		if err != nil {
			log.Printf("更新用户为管理员失败: %v", err)
		} else {
			fmt.Println("已将用户 'admin' 升级为管理员")
		}
		return
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("加密密码失败: %v", err)
		return
	}

	// 创建管理员用户
	_, err = db.DB.Exec("INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)",
		username, email, string(hashedPassword), db.RoleAdmin)
	if err != nil {
		log.Printf("创建管理员失败: %v", err)
	} else {
		fmt.Println("管理员账户已创建！")
	}
}
