package main

import (
	"APIGateway/database"
	"APIGateway/handlers"
	"APIGateway/middleware"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:love@tcp(127.0.0.1:3306)/mydatabase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	err = repo.SetupDatabase()
	if err != nil {
		log.Fatal("Cannot setup database:", err)
	}

	handler := handlers.NewHandler(repo)

	if handlers.NewsServiceURL == "" || handlers.CommentsServiceURL == "" {
		log.Fatal("NEWS_SERVICE_URL or COMMENT_SERVICE_URL not set")
	}

	// Заменяем функции обработчика в соответствии с новым типом Handler
	newsHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.NewsHandler))
	newsDetailHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.NewsDetailHandler))
	newsFilterHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.NewsFilterHandler))
	addCommentHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.AddComment))
	getCommentsHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.GetComments))
	http.Handle("/comments/get", getCommentsHandler)
	http.Handle("/news", newsHandler)
	http.Handle("/news/details", newsDetailHandler)
	http.Handle("/news/filter", newsFilterHandler)
	http.Handle("/comments/add", addCommentHandler)

	http.ListenAndServe(":8080", nil)

}
