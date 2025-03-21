package db

import (
	"time"
)

// Article 文章数据模型
type Article struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	Likes        int       `json:"likes"`
	CommentCount int       `json:"comment_count"`
	CategoryID   int       `json:"category_id"`
}

// GetArticlesByCategoryAndOrder 根据分类 ID 和排序条件查询文章列表
func GetArticlesByCategoryAndOrder(categoryID int, orderClause string) ([]Article, error) {
	query := "SELECT id, title, content, created_at, likes, comment_count, category_id FROM articles WHERE category_id = ? ORDER BY " + orderClause
	rows, err := DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var art Article
		if err := rows.Scan(&art.ID, &art.Title, &art.Content, &art.CreatedAt, &art.Likes, &art.CommentCount, &art.CategoryID); err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}

	return articles, nil
}

// GetArticleByID 根据文章ID获取单篇文章详情
func GetArticleByID(id int) (*Article, error) {
	query := "SELECT id, title, content, created_at, likes, comment_count, category_id FROM articles WHERE id = ? AND is_visible = 1"
	var article Article
	err := DB.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.CreatedAt,
		&article.Likes,
		&article.CommentCount,
		&article.CategoryID,
	)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// 修改原有的 Comment 结构体，添加 UserID 字段
type Comment struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsVisible int       `json:"is_visible"`
	Likes     int       `json:"likes"`
	UserID    int       `json:"user_id"` // 添加用户ID字段
}

// 需要相应修改 GetCommentsByArticleID 函数
func GetCommentsByArticleID(articleID int) ([]Comment, error) {
	query := "SELECT id, article_id, content, created_at, is_visible, likes, user_id FROM comments WHERE article_id = ? AND is_visible = 1 ORDER BY created_at DESC"
	rows, err := DB.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.ArticleID, &comment.Content, &comment.CreatedAt, &comment.IsVisible, &comment.Likes, &comment.UserID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
