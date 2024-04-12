// censorship_service/main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Comment структура комментария.
type Comment struct {
	Text string `json:"text"`
}

var forbiddenWords = []string{"qwerty", "йцукен", "zxvbnm"} // Specific words not allowed

func main() {
	http.HandleFunc("/censor", censorHandler)
	http.ListenAndServe(":8080", nil)
}

// censorHandler обработчик HTTP-запросов, проверяет комментарии на наличие запрещенных слов.
func censorHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var comment Comment
	err = json.Unmarshal(body, &comment)
	if err != nil {
		http.Error(w, "Invalid comment format", http.StatusBadRequest)
		return
	}

	lowercaseComment := strings.ToLower(comment.Text)
	for _, word := range forbiddenWords {
		if strings.Contains(lowercaseComment, word) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
