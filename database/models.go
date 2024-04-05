// Package repository - models.go
package repository

import (
	"net/http"
	"time"
)

// Result структура для хранения результата вычислений.
type Result struct {
	Data     interface{}
	Response *http.Response
	Error    error
}

// Comment структура для представления комментария.
type Comment struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	NewsID    int       `json:"news_id"`
	ParentID  *int      `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewsShortDetailed представляет основную информацию о новости.
type NewsShortDetailed struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// NewsFullDetailed струкутра для полной информации о новости и связанных комментариях.
type NewsFullDetailed struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}
