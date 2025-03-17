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
