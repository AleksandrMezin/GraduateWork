// news_service/main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type NewsAPIResponse struct {
	Articles []struct {
		Source struct {
			Name string `json:"name"`
		} `json:"source"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		UrlToImage  string `json:"urlToImage"`
		PublishedAt string `json:"publishedAt"`
		Content     string `json:"content"`
	} `json:"articles"`
}

func handleNews(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=us&apiKey=818fab7db99641809f117b29d2ffcfe8")
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

	var newsAPIResponse NewsAPIResponse

	err = json.Unmarshal(body, &newsAPIResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(newsAPIResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/news", handleNews)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
