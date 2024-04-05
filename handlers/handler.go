// Package handlers - handler.go
package handlers

import (
	"APIGateway/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	NewsServiceURL     = os.Getenv("NEWS_SERVICE_URL")
	CommentsServiceURL = os.Getenv("COMMENT_SERVICE_URL")
	requestTimeout     = 2 * time.Second
)

var CensorshipServiceURL = "http://localhost:8080/censor"

// decodeJSON декодирует JSON из тела запроса.
func decodeJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}

// createCommentInService создает комментарий в другом сервисе.
func createCommentInService(input *repository.Comment) (*repository.Comment, error) {
	client := &http.Client{}
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", CommentsServiceURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create comment, status code: %d", resp.StatusCode)
	}

	var comment repository.Comment
	if err := decodeJSON(resp.Body, &comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

// NewHandler создает и возвращает новую структуру Handler
func NewHandler(repo repository.RepositoryInterface) *Handler {
	return &Handler{Repo: repo}
}

// AddComment обрабатывает HTTP POST запросы и добавлет комментарий.
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var comment repository.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment.CreatedAt = time.Now()

	// Call the censorship service and check the comment
	censorResult, err := h.CensorComment(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !censorResult {
		http.Error(w, "Forbidden content in comment", http.StatusForbidden)
		return
	}

	if err := h.Repo.Save(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CensorComment отправляет запрос для проверки содержимого комментария.
func (h *Handler) CensorComment(comment *repository.Comment) (bool, error) {
	client := &http.Client{}
	reqBody, err := json.Marshal(comment)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", CensorshipServiceURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// GetComments обрабатывает HTTP GET запросы и возвращает комментарии.
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(r.URL.Query().Get("news_id"))
	if err != nil {
		http.Error(w, "'news_id' parameter is required", http.StatusBadRequest)
	}

	comments, err := h.Repo.GetCommentsByNewsID(newsID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}
