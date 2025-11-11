package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/n-korel/CSVviewer-app/internal/storage"
)

func TestUploadCSV(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := NewHandler(store)

	// Создаем CSV данные
	csvContent := "Name,Age,City\nJohn,30,New York\nJane,25,London\n"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.UploadCSV(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["row_count"].(float64) != 2 {
		t.Errorf("Expected 2 rows, got %v", response["row_count"])
	}
}

func TestGetData(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := NewHandler(store)

	// Предзагружаем данные
	csvContent := "Name,Age\nJohn,30\nJane,25\n"
	uploadCSV(handler, csvContent)

	req := httptest.NewRequest("GET", "/api/data?page=1&per_page=10", nil)
	w := httptest.NewRecorder()

	handler.GetData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["total"].(float64) != 2 {
		t.Errorf("Expected total 2, got %v", response["total"])
	}
}

func TestSearchData(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := NewHandler(store)

	// Предзагружаем данные
	csvContent := "Name,Age\nJohn,30\nJane,25\n"
	uploadCSV(handler, csvContent)

	req := httptest.NewRequest("GET", "/api/search?q=John&page=1&per_page=10", nil)
	w := httptest.NewRecorder()

	handler.SearchData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["total"].(float64) != 1 {
		t.Errorf("Expected 1 search result, got %v", response["total"])
	}
}

func TestClearData(t *testing.T) {
	store := storage.NewMemoryStorage()
	handler := NewHandler(store)

	// Предзагружаем данные
	csvContent := "Name,Age\nJohn,30\n"
	uploadCSV(handler, csvContent)

	req := httptest.NewRequest("DELETE", "/api/clear", nil)
	w := httptest.NewRecorder()

	handler.ClearData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем что данные очищены
	data, _ := store.GetAll()
	if len(data.Rows) != 0 {
		t.Errorf("Expected 0 rows after clear, got %d", len(data.Rows))
	}
}

// Функция для загрузки CSV в тестах
func uploadCSV(handler *Handler, csvContent string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.UploadCSV(w, req)
}