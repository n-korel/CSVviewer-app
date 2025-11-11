package storage

import (
	"strings"
	"sync"

	"github.com/n-korel/CSVviewer-app/internal/models"
)

// Работа с данными
type Storage interface {
	Store(data *models.CSVData) error
	GetAll() (*models.CSVData, error)
	GetPaginated(page, perPage int) (*models.PaginatedResponse, error)
	Search(query string, page, perPage int) (*models.PaginatedResponse, error)
	Clear() error
}

// Хранит данные в памяти
type MemoryStorage struct {
	mu   sync.RWMutex
	data *models.CSVData
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: &models.CSVData{
			Headers: []string{},
			Rows:    []map[string]interface{}{},
			Total:   0,
		},
	}
}

// Сохраняет данные
func (s *MemoryStorage) Store(data *models.CSVData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = data
	return nil
}

// Возвращает все данные
func (s *MemoryStorage) GetAll() (*models.CSVData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data, nil
}

// Возвращает данные с пагинацией
func (s *MemoryStorage) GetPaginated(page, perPage int) (*models.PaginatedResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.data.Rows)
	start := (page - 1) * perPage
	end := start + perPage

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &models.PaginatedResponse{
		Data:    s.data.Rows[start:end],
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Headers: s.data.Headers,
	}, nil
}

// Ищет данные по запросу
func (s *MemoryStorage) Search(query string, page, perPage int) (*models.PaginatedResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	var filtered []map[string]interface{}

	for _, row := range s.data.Rows {
		for _, value := range row {
			if str, ok := value.(string); ok {
				if strings.Contains(strings.ToLower(str), query) {
					filtered = append(filtered, row)
					break
				}
			}
		}
	}

	total := len(filtered)
	start := (page - 1) * perPage
	end := start + perPage

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &models.PaginatedResponse{
		Data:    filtered[start:end],
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Headers: s.data.Headers,
	}, nil
}

// Очищает все данные
func (s *MemoryStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = &models.CSVData{
		Headers: []string{},
		Rows:    []map[string]interface{}{},
		Total:   0,
	}
	return nil
}