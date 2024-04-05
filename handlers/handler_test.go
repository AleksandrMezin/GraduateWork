// handlers_test.go
package handlers

import (
	"APIGateway/database"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockRepository struct {
	SaveFunc                func(repository.Comment) error
	GetCommentsByNewsIDFunc func(nid int) ([]repository.Comment, error)
	SetupDatabaseFunc       func() error
}

func (m *MockRepository) Save(c repository.Comment) error {
	return m.SaveFunc(c)
}

func (m *MockRepository) GetCommentsByNewsID(nid int) ([]repository.Comment, error) {
	return m.GetCommentsByNewsIDFunc(nid)
}

func (m *MockRepository) SetupDatabase() error {
	return m.SetupDatabaseFunc()
}

func TestAddComment(t *testing.T) {
	reqBody := bytes.NewBuffer([]byte(`{"author":"John","text":"Test comment","news_id":1,"parent_id":null,"created_at":"` + time.Now().Format(time.RFC3339Nano) + `"}`))
	req, _ := http.NewRequest("POST", "/addComment", reqBody)
	rr := httptest.NewRecorder()

	h := &Handler{
		Repo: &MockRepository{
			SaveFunc: func(c repository.Comment) error {
				return nil
			},
			SetupDatabaseFunc: func() error {
				return nil
			},
		},
	}
	h.AddComment(rr, req)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("Expected Status %v but got %v", http.StatusCreated, rr.Result().StatusCode)
	}
}

func TestGetComments(t *testing.T) {
	req, _ := http.NewRequest("GET", "/getComment?news_id=1", nil)
	rr := httptest.NewRecorder()

	h := &Handler{
		Repo: &MockRepository{
			GetCommentsByNewsIDFunc: func(nid int) ([]repository.Comment, error) {
				return []repository.Comment{
					{Author: "Test Author", Text: "Test Text", NewsID: 1, CreatedAt: time.Now()},
				}, nil
			},
			SetupDatabaseFunc: func() error {
				return nil
			},
		},
	}
	h.GetComments(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected Status %v, but got %v", http.StatusOK, rr.Result().StatusCode)
	}
}
