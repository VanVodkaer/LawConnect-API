package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VanVodkaer/LawConnect-API/internal/db"
	"github.com/VanVodkaer/LawConnect-API/utils/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token  string  `json:"token"`
	User   db.User `json:"user"`
	Expire int64   `json:"expire"`
}

// Claims 自定义JWT声明结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// Login 处理用户登录请求
func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求格式不正确",
		})
		return
	}

	// 根据用户名查找用户
	user, err := db.GetUserByUsername(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码不正确",
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(loginReq.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码不正确",
		})
		return
	}

	// 生成JWT令牌
	expireTime := time.Now().Add(time.Duration(config.GlobalConfig.JWT.Expire) * time.Hour)
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败",
		})
		return
	}

	// 返回令牌和用户信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": LoginResponse{
			Token:  tokenString,
			User:   *user,
			Expire: expireTime.Unix(),
		},
	})
}

// Register 处理用户注册请求
func Register(c *gin.Context) {
	var registerReq struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求格式不正确",
		})
		return
	}

	// 创建新用户
	user := &db.User{
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Password: registerReq.Password,
		Role:     db.RoleUser, // 默认为普通用户
	}

	// 保存到数据库
	if err := db.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
