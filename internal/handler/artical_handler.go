package handler

import (
	"net/http"
	"strconv"

	"github.com/VanVodkaer/LawConnect-API/internal/db"
	"github.com/gin-gonic/gin"
)

// ----------------------- 法学交流社区 -----------------------

// GetCommunityLatest 获取法学交流社区【最新动态】（按发布时间倒序）
func GetCommunityLatest(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(1, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetCommunityHottest 获取法学交流社区【最热帖子】（按点赞数量倒序）
func GetCommunityHottest(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(1, "likes DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetCommunityHotQA 获取法学交流社区【热门问答】（按评论数量倒序）
func GetCommunityHotQA(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(1, "comment_count DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// ----------------------- 政策推送专区 -----------------------

// GetPolicyLatest 获取政策推送专区【最新政策】（parent category_id=2）
func GetPolicyLatest(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(2, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetPolicyLocal 获取政策推送专区【地方政策】（category_id=4）
func GetPolicyLocal(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(4, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetPolicyInterpretation 获取政策推送专区【政策解读】（category_id=5）
func GetPolicyInterpretation(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(5, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// ----------------------- 线下实践平台 -----------------------

// GetOfflineCooperation 获取线下实践平台【线下联动】（category_id=6）
func GetOfflineCooperation(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(6, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetOfflineOnline 获取线下实践平台【线上活动】（category_id=7）
func GetOfflineOnline(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(7, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetOfflineRegistration 获取线下实践平台【报名中心】（category_id=8）
func GetOfflineRegistration(c *gin.Context) {
	articles, err := db.GetArticlesByCategoryAndOrder(8, "created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": articles})
}

// GetArticleDetail 获取文章详情及评论
func GetArticleDetail(c *gin.Context) {
	// 获取文章ID参数
	id := c.Param("id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的文章ID"})
		return
	}

	// 查询文章详情
	article, err := db.GetArticleByID(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 查询文章评论
	comments, err := db.GetCommentsByArticleID(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询评论失败"})
		return
	}

	// 返回文章详情和评论
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"article":  article,
			"comments": comments,
		},
	})
}
