package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/VanVodkaer/LawConnect-API/internal/db"
	"github.com/VanVodkaer/LawConnect-API/utils/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// 自定义JWT声明结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 检查令牌格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式不正确",
			})
			c.Abort()
			return
		}

		// 解析令牌
		tokenString := parts[1]
		claims := &Claims{}

		// 使用密钥解析令牌
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 确保令牌使用了预期的签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("无效的签名方法")
			}
			return []byte(config.GlobalConfig.JWT.Secret), nil
		})

		// 处理解析错误
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "认证令牌已过期",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "无效的认证令牌",
				})
			}
			c.Abort()
			return
		}

		// 验证令牌是否有效
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌已过期或无效",
			})
			c.Abort()
			return
		}

		// 获取用户信息并存储在上下文中
		user, err := db.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户不存在或已被删除",
			})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中以便后续处理
		c.Set("user", user)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminRequired 验证用户是否为管理员的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "需要认证",
			})
			c.Abort()
			return
		}

		// 检查用户是否为管理员
		u, ok := user.(*db.User)
		if !ok || !u.IsAdmin() {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	// 从上下文中获取用户
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要认证",
		})
		return
	}

	u, ok := user.(*db.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 生成新令牌
	expireTime := time.Now().Add(time.Duration(config.GlobalConfig.JWT.Expire) * time.Hour)
	claims := Claims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   u.Username,
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

	// 返回新令牌
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "刷新令牌成功",
		"data": gin.H{
			"token":  tokenString,
			"expire": expireTime.Unix(),
		},
	})
}
