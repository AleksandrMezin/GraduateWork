// Package repository - repository.go
package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

// Repository обертка над DB.
type Repository struct {
	db *sql.DB
}

// RepositoryInterface interface declaration
type RepositoryInterface interface {
	Save(c Comment) error
	GetCommentsByNewsID(newsID int) ([]Comment, error)
	SetupDatabase() error
}

// NewRepository создает новый репозиторий.
func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{db: db}
}

// ForbiddenWords запрещённые слова для комментариев.
var ForbiddenWords = []string{"qwerty", "йцукен", "zxvbnm"}

// ModerateComment проверяет, что комментарий не содержит запрещенных слов.
func ModerateComment(comment *Comment) {
	lowercaseBody := strings.ToLower(comment.Text)
	for _, word := range ForbiddenWords {
		match, _ := regexp.MatchString(word, lowercaseBody)
		if match {
			comment.Text = "This comment has been moderated"
			break
		}
	}
}

// Save сохраняет объект комментария в базу данных.
func (r *Repository) Save(comment Comment) error {
	_, err := r.db.Exec(`
		INSERT INTO comments (author, text, news_id, parent_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, comment.Author, comment.Text, comment.NewsID, comment.ParentID, comment.CreatedAt)
	return err
}

// GetCommentsByNewsID извлекает все комментарии, связанные с определенной новостью.
func (r *Repository) GetCommentsByNewsID(newsID int) ([]Comment, error) {
	rows, err := r.db.Query(`
		SELECT id, author, text, news_id, parent_id, created_at 
		FROM comments 
		WHERE news_id = ?
	`, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Author, &comment.Text, &comment.NewsID, &comment.ParentID, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// SetupDatabase создает таблицу комментариев, если она еще не существует.
func (r *Repository) SetupDatabase() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			author VARCHAR(255),
			text TEXT,
			news_id INT,
			parent_id INT,
			created_at DATETIME
		)
	`)
	return err
}
