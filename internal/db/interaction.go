package db

import (
	"errors"
)

// AddComment 添加评论到文章
func AddComment(articleID int, content string, userID int) (int, error) {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. 插入评论记录
	res, err := tx.Exec(
		"INSERT INTO comments (article_id, content, is_visible, user_id) VALUES (?, ?, ?, ?)",
		articleID, content, 1, userID,
	)
	if err != nil {
		return 0, err
	}

	// 获取新插入评论的ID
	commentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 2. 更新文章的评论计数
	_, err = tx.Exec(
		"UPDATE articles SET comment_count = comment_count + 1 WHERE id = ?",
		articleID,
	)
	if err != nil {
		return 0, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(commentID), nil
}

// LikeArticle 为文章点赞，防止重复点赞
func LikeArticle(articleID int, userID int) error {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 检查文章是否存在
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id = ? AND is_visible = 1)", articleID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("文章不存在或已被删除")
	}

	// 检查是否已经点赞过
	var alreadyLiked bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM article_likes WHERE article_id = ? AND user_id = ?)", articleID, userID).Scan(&alreadyLiked)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return errors.New("您已经点赞过该文章")
	}

	// 记录点赞
	_, err = tx.Exec("INSERT INTO article_likes (article_id, user_id) VALUES (?, ?)", articleID, userID)
	if err != nil {
		return err
	}

	// 更新文章点赞数
	_, err = tx.Exec("UPDATE articles SET likes = likes + 1 WHERE id = ?", articleID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// UnlikeArticle 取消文章点赞
func UnlikeArticle(articleID int, userID int) error {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已点赞
	var alreadyLiked bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM article_likes WHERE article_id = ? AND user_id = ?)", articleID, userID).Scan(&alreadyLiked)
	if err != nil {
		return err
	}
	if !alreadyLiked {
		return errors.New("您尚未点赞该文章")
	}

	// 删除点赞记录
	_, err = tx.Exec("DELETE FROM article_likes WHERE article_id = ? AND user_id = ?", articleID, userID)
	if err != nil {
		return err
	}

	// 更新文章点赞数
	_, err = tx.Exec("UPDATE articles SET likes = likes - 1 WHERE id = ? AND likes > 0", articleID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// LikeComment 为评论点赞，防止重复点赞
func LikeComment(commentID int, userID int) error {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 检查评论是否存在且可见
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM comments WHERE id = ? AND is_visible = 1)", commentID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("评论不存在或未通过审核")
	}

	// 检查是否已经点赞过
	var alreadyLiked bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = ? AND user_id = ?)", commentID, userID).Scan(&alreadyLiked)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return errors.New("您已经点赞过该评论")
	}

	// 记录点赞
	_, err = tx.Exec("INSERT INTO comment_likes (comment_id, user_id) VALUES (?, ?)", commentID, userID)
	if err != nil {
		return err
	}

	// 更新评论点赞数
	_, err = tx.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// UnlikeComment 取消评论点赞
func UnlikeComment(commentID int, userID int) error {
	// 开启事务
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已点赞
	var alreadyLiked bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = ? AND user_id = ?)", commentID, userID).Scan(&alreadyLiked)
	if err != nil {
		return err
	}
	if !alreadyLiked {
		return errors.New("您尚未点赞该评论")
	}

	// 删除点赞记录
	_, err = tx.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
	if err != nil {
		return err
	}

	// 更新评论点赞数
	_, err = tx.Exec("UPDATE comments SET likes = likes - 1 WHERE id = ? AND likes > 0", commentID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// CheckArticleLikeStatus 检查用户是否已点赞文章
func CheckArticleLikeStatus(articleID int, userID int) (bool, error) {
	var liked bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM article_likes WHERE article_id = ? AND user_id = ?)",
		articleID, userID).Scan(&liked)
	return liked, err
}

// CheckCommentLikeStatus 检查用户是否已点赞评论
func CheckCommentLikeStatus(commentID int, userID int) (bool, error) {
	var liked bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = ? AND user_id = ?)",
		commentID, userID).Scan(&liked)
	return liked, err
}
