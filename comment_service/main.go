// comment_service/main.go
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

	addCommentHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.AddComment))
	getCommentsHandler := middleware.LoggingMiddleware(http.HandlerFunc(handler.GetComments))

	http.Handle("/comments/get", getCommentsHandler)
	http.Handle("/comments/add", addCommentHandler)

	log.Println("Comment service started on port 8081")
	http.ListenAndServe(":8081", nil)
}
