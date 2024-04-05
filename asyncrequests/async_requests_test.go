// async_request_test.go
package asyncrequests

import (
	"APIGateway/handlers"
	"testing"
)

func TestExecuteAsyncHTTPGets(t *testing.T) {
	handlers.GetNews = func(url string, title string, page int, pageSize int) ([]handlers.News, *handlers.Pagination, error) {
		news := []handlers.News{{Title: "test title", Author: "test author", Published: "test date"}}
		pagination := &handlers.Pagination{CurrentPage: 1, TotalPages: 1}
		return news, pagination, nil
	}

	urls := []string{"test1", "test2"}
	title := "test title"
	page := 1
	pageSize := 1
	news, pagination, err := ExecuteAsyncHTTPGets(urls, title, page, pageSize)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	if len(news) != len(urls) {
		t.Errorf("Expected %v news but got %v due to length mismatch", len(urls), len(news))
		return
	}

	for _, n := range news {
		if n.Title != title {
			t.Errorf("Expected news title %q but got %q", title, n.Title)
		}
	}

	if pagination.TotalPages != len(urls) {
		t.Errorf("Expected total number of pages to be %v but got %v", len(urls), pagination.TotalPages)
	}
}
