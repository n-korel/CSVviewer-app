package storage

import (
	"testing"

	"github.com/n-korel/CSVviewer-app/internal/models"
)

func TestMemoryStorage_Store(t *testing.T) {
	store := NewMemoryStorage()
	data := &models.CSVData{
		Headers: []string{"Name", "Age"},
		Rows: []map[string]interface{}{
			{"Name": "John", "Age": "30"},
			{"Name": "Jane", "Age": "25"},
		},
		Total: 2,
	}

	err := store.Store(data)
	if err != nil {
		t.Errorf("Store() error = %v", err)
	}

	retrieved, err := store.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
	}

	if len(retrieved.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(retrieved.Rows))
	}
}

func TestMemoryStorage_GetPaginated(t *testing.T) {
	store := NewMemoryStorage()
	data := &models.CSVData{
		Headers: []string{"ID"},
		Rows:    make([]map[string]interface{}, 100),
		Total:   100,
	}

	for i := 0; i < 100; i++ {
		data.Rows[i] = map[string]interface{}{"ID": i}
	}

	store.Store(data)

	// Тест первой страницы
	result, err := store.GetPaginated(1, 10)
	if err != nil {
		t.Errorf("GetPaginated() error = %v", err)
	}

	if len(result.Data) != 10 {
		t.Errorf("Expected 10 items, got %d", len(result.Data))
	}

	if result.Total != 100 {
		t.Errorf("Expected total 100, got %d", result.Total)
	}

	// Тест последней страницы
	result, err = store.GetPaginated(10, 10)
	if err != nil {
		t.Errorf("GetPaginated() error = %v", err)
	}

	if len(result.Data) != 10 {
		t.Errorf("Expected 10 items on last page, got %d", len(result.Data))
	}
}

func TestMemoryStorage_Search(t *testing.T) {
	store := NewMemoryStorage()
	data := &models.CSVData{
		Headers: []string{"Name", "City"},
		Rows: []map[string]interface{}{
			{"Name": "John", "City": "New York"},
			{"Name": "Jane", "City": "London"},
			{"Name": "Bob", "City": "New York"},
		},
		Total: 3,
	}

	store.Store(data)

	// Поиск по "New York"
	result, err := store.Search("New York", 1, 10)
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 results for 'New York', got %d", len(result.Data))
	}

	// Поиск по "Jane"
	result, err = store.Search("Jane", 1, 10)
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 result for 'Jane', got %d", len(result.Data))
	}

	// Поиск несуществующего
	result, err = store.Search("NonExistent", 1, 10)
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}

	if len(result.Data) != 0 {
		t.Errorf("Expected 0 results for 'NonExistent', got %d", len(result.Data))
	}
}

func TestMemoryStorage_Clear(t *testing.T) {
	store := NewMemoryStorage()
	data := &models.CSVData{
		Headers: []string{"Name"},
		Rows:    []map[string]interface{}{{"Name": "John"}},
		Total:   1,
	}

	store.Store(data)
	store.Clear()

	retrieved, _ := store.GetAll()
	if len(retrieved.Rows) != 0 {
		t.Errorf("Expected 0 rows after clear, got %d", len(retrieved.Rows))
	}
}