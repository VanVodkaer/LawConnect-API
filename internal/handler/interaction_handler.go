package handler

import (
	"net/http"
	"strconv"

	"github.com/VanVodkaer/LawConnect-API/internal/db"
	"github.com/gin-gonic/gin"
)

// AddCommentRequest 添加评论的请求结构
type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// AddComment 添加评论处理程序
func AddComment(c *gin.Context) {
	// 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的文章ID",
		})
		return
	}

	// 验证文章是否存在
	article, err := db.GetArticleByID(articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在或已被删除",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录后才能评论",
		})
		return
	}

	// 解析请求体
	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 内容不能为空
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "评论内容不能为空",
		})
		return
	}

	// 添加评论
	commentID, err := db.AddComment(articleID, req.Content, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "添加评论失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "评论成功",
		"data": gin.H{
			"comment_id": commentID,
			"article_id": article.ID,
		},
	})
}

// LikeArticle 文章点赞处理程序
func LikeArticle(c *gin.Context) {
	// 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的文章ID",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录后才能点赞",
		})
		return
	}

	// 为文章点赞
	err = db.LikeArticle(articleID, userID.(int))
	if err != nil {
		// 如果是重复点赞错误，返回特定响应
		if err.Error() == "您已经点赞过该文章" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "点赞失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "点赞成功",
		"data": gin.H{
			"article_id": articleID,
		},
	})
}

// UnlikeArticle 取消文章点赞处理程序
func UnlikeArticle(c *gin.Context) {
	// 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的文章ID",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录后才能取消点赞",
		})
		return
	}

	// 取消文章点赞
	err = db.UnlikeArticle(articleID, userID.(int))
	if err != nil {
		// 如果是未点赞错误，返回特定响应
		if err.Error() == "您尚未点赞该文章" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "取消点赞失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "取消点赞成功",
		"data": gin.H{
			"article_id": articleID,
		},
	})
}

// LikeComment 评论点赞处理程序
func LikeComment(c *gin.Context) {
	// 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的评论ID",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录后才能点赞",
		})
		return
	}

	// 为评论点赞
	err = db.LikeComment(commentID, userID.(int))
	if err != nil {
		// 如果是重复点赞错误，返回特定响应
		if err.Error() == "您已经点赞过该评论" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "点赞失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "点赞成功",
		"data": gin.H{
			"comment_id": commentID,
		},
	})
}

// UnlikeComment 取消评论点赞处理程序
func UnlikeComment(c *gin.Context) {
	// 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的评论ID",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录后才能取消点赞",
		})
		return
	}

	// 取消评论点赞
	err = db.UnlikeComment(commentID, userID.(int))
	if err != nil {
		// 如果是未点赞错误，返回特定响应
		if err.Error() == "您尚未点赞该评论" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "取消点赞失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "取消点赞成功",
		"data": gin.H{
			"comment_id": commentID,
		},
	})
}

// GetArticleLikeStatus 获取文章点赞状态
func GetArticleLikeStatus(c *gin.Context) {
	// 获取文章ID
	articleIDStr := c.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的文章ID",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录",
		})
		return
	}

	// 检查点赞状态
	liked, err := db.CheckArticleLikeStatus(articleID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取点赞状态失败: " + err.Error(),
		})
		return
	}

	// 返回点赞状态
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"article_id": articleID,
			"liked":      liked,
		},
	})
}

// GetCommentLikeStatus 获取评论点赞状态
func GetCommentLikeStatus(c *gin.Context) {
	// 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的评论ID",
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "需要登录",
		})
		return
	}

	// 检查点赞状态
	liked, err := db.CheckCommentLikeStatus(commentID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取点赞状态失败: " + err.Error(),
		})
		return
	}

	// 返回点赞状态
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"comment_id": commentID,
			"liked":      liked,
		},
	})
}
