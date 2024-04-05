// repository_test.go
package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func TestModerateComment(t *testing.T) {
	tt := []struct {
		input    Comment
		expected Comment
	}{
		{Comment{Text: "qwerty"}, Comment{Text: "This comment has been moderated"}},
		{Comment{Text: "Йцукен"}, Comment{Text: "This comment has been moderated"}},
		{Comment{Text: "AnalyticMind"}, Comment{Text: "AnalyticMind"}},
	}
	for _, tc := range tt {
		ModerateComment(&tc.input)
		if tc.input.Text != tc.expected.Text {
			t.Errorf("got %q, want %q", tc.input.Text, tc.expected.Text)
		}
	}
}

func TestSave(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	parentID := 2
	comment := Comment{
		Author:    "John",
		Text:      "Hello, world!",
		NewsID:    1,
		ParentID:  &parentID,
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO comments").
		WithArgs(comment.Author, comment.Text, comment.NewsID, *comment.ParentID, comment.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Save(comment); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCommentsByNewsID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	newsID := 1
	rows := sqlmock.NewRows([]string{"id", "author", "text", "news_id", "parent_id", "created_at"}).
		AddRow(1, "John", "Hello, world!", newsID, 2, time.Now())

	mock.ExpectQuery("SELECT id, author, text, news_id, parent_id, created_at FROM comments WHERE news_id = ?").
		WithArgs(newsID).
		WillReturnRows(rows)

	comments, err := repo.GetCommentsByNewsID(newsID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(comments) != 1 {
		t.Errorf("expected one comment, got %d", len(comments))
	}

	if comments[0].Author != "John" || comments[0].Text != "Hello, world!" || comments[0].NewsID != newsID {
		t.Errorf("wrong comment data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
