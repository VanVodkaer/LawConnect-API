package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // 密码不通过JSON返回
	Role     int    `json:"role"`
}

// 角色常量
const (
	RoleUser  = 1 // 普通用户
	RoleAdmin = 2 // 管理员
)

// GetUserByEmail 通过邮箱获取用户
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password, IFNULL(role, 1) FROM users WHERE email = ?"
	err := DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password, IFNULL(role, 1) FROM users WHERE username = ?"
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// GetUserByID 通过ID获取用户
func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, password, IFNULL(role, 1) FROM users WHERE id = ?"
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 插入用户数据
	query := "INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)"
	result, err := DB.Exec(query, user.Username, user.Email, string(hashedPassword), user.Role)
	if err != nil {
		return err
	}

	// 获取自增ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

// CheckPassword 验证用户密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// AdminExists 检查是否存在管理员用户
func AdminExists() (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE role = ?"
	err := DB.QueryRow(query, RoleAdmin).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
