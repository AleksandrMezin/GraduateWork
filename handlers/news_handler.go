// Package handlers - news_handler.go
package handlers

import (
	"APIGateway/asyncrequests"
	repository "APIGateway/database"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	BBCAPI       = "https://www.bbc.co.uk/news/"
	NYTAPI       = "https://developer.nytimes.com/apis"
	requestIDKey = "requestID"
)

type Handler struct {
	Repo repository.RepositoryInterface
}

type News struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published string `json:"published"`
	Content   string `json:"content"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

type NewsResponse struct {
	RequestID  string      `json:"requestId"`
	News       []News      `json:"news"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// makeGETRequest делает GET-запрос и возвращает тело ответа.
func makeGETRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// respondWithError отправляет HTTP-ответ с кодом ошибки и сообщением.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

var GetNews = GetNewsFromService

func GetNewsFromService(serviceURL string, title string, page int, pageSize int) ([]News, *Pagination, error) {
	client := &http.Client{}

	URL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, nil, err
	}

	// append title, page and pageSize query parameters to the URL
	query := URL.Query()
	query.Add("title", title)
	query.Add("page", strconv.Itoa(page))
	query.Add("pageSize", strconv.Itoa(pageSize))
	URL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var data map[string]json.RawMessage
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, nil, err
	}

	var news []News
	err = json.Unmarshal(data["news"], &news)
	if err != nil {
		return nil, nil, err
	}

	var pagination Pagination
	err = json.Unmarshal(data["pagination"], &pagination)
	if err != nil {
		return nil, nil, err
	}

	return news, &pagination, nil
}

var log = logrus.New()

// NewsHandler обрабатывает запросы и возвращает новости.
func (h *Handler) NewsHandler(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDKey).(string)
	if !ok {
		log.Error("requestID missing from context")
		respondWithError(w, http.StatusBadRequest, "requestID missing from context")
		return
	}

	searchQuery := r.URL.Query().Get("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid or missing page number")
		return
	}

	pageSize := 10 // define the appropriate page size

	serviceURLs := []string{BBCAPI, NYTAPI}

	news, pagination, err := asyncrequests.ExecuteAsyncHTTPGets(serviceURLs, searchQuery, page, pageSize)
	if err != nil {
		log.Error("Failed to get news from services", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get news from services")
		return
	}

	response := NewsResponse{
		RequestID:  requestID,
		News:       news,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, &response)
}

// NewsDetailHandler обрабатывает запросы и возвращает детали новости.
func (h *Handler) NewsDetailHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid or missing page number")
		return
	}

	pageSize := 10 // define the appropriate page size

	details, pagination, err := GetNewsFromService(BBCAPI, searchQuery, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get news details"))
		return
	}

	response := NewsResponse{
		RequestID:  h.getContextRequestID(r),
		News:       details,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, &response)
}

func (h *Handler) NewsFilterHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid or missing page number")
		return
	}

	pageSize := 10 // define the appropriate page size

	filterServiceURL := "https://news-filter-service.com"
	filtered, pagination, err := GetNewsFromService(filterServiceURL, searchQuery, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get filtered news"))
		return
	}

	response := NewsResponse{
		RequestID:  h.getContextRequestID(r),
		News:       filtered,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, &response)
}

// getContextRequestID извлекает и возвращает идентификатор запроса из контекста.
func (h *Handler) getContextRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

// sendJSONResponse отправляет ответ в формате JSON.
func (h *Handler) sendJSONResponse(w http.ResponseWriter, response *NewsResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ForwardNewsRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=us&apiKey=818fab7db99641809f117b29d2ffcfe8") // замените на URL вашего сервиса новостей.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (h *Handler) ForwardNewsRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=us&apiKey=818fab7db99641809f117b29d2ffcfe8") // Замените на URL вашего API новостей.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
