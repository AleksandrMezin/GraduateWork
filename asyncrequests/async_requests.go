// Package asyncrequests/async_requests.go
package asyncrequests

import (
	"APIGateway/handlers"
	"sync"
)

// NewsWithError Добавляем дополнительное поле Pagination в структуру
type NewsWithError struct {
	News       []handlers.News
	Pagination *handlers.Pagination
	Err        error
}

func ExecuteAsyncHTTPGets(urls []string, title string, page int, pageSize int) ([]handlers.News, *handlers.Pagination, error) {
	var wg sync.WaitGroup
	newsChannel := make(chan NewsWithError, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			news, pagination, err := handlers.GetNewsFromService(url, title, page, pageSize)
			newsChannel <- NewsWithError{news, pagination, err}
		}(url)
	}

	wg.Wait()
	close(newsChannel)

	var allNews []handlers.News
	var allPagination *handlers.Pagination
	for nw := range newsChannel {
		if nw.Err != nil {
			return nil, allPagination, nw.Err
		}
		allNews = append(allNews, nw.News...)
		// Здесь мы обрабатываем сумму страниц от всех пагинаций
		allPagination.TotalPages += nw.Pagination.TotalPages
	}
	return allNews, allPagination, nil
}
