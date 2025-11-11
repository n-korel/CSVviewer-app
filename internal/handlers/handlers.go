package handlers

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/n-korel/CSVviewer-app/internal/models"
	"github.com/n-korel/CSVviewer-app/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
	return &Handler{storage: s}
}

// Обработка загрузки CSV файла
func (h *Handler) UploadCSV(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 10 << 20 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        respondError(w, "file too large", http.StatusRequestEntityTooLarge)
    return
    }

	file, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		respondError(w, "Failed to read CSV headers", http.StatusBadRequest)
		return
	}

	var rows []map[string]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // пропускаем ошибочные строки
		}

		row := make(map[string]interface{})
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		rows = append(rows, row)
	}

	csvData := &models.CSVData{
		Headers: headers,
		Rows:    rows,
		Total:   len(rows),
	}

	if err := h.storage.Store(csvData); err != nil {
		respondError(w, "Failed to store data", http.StatusInternalServerError)
		return
	}

	// Возвращаем образец данных (первые 5 строк)
	sampleSize := 5
	if len(rows) < sampleSize {
		sampleSize = len(rows)
	}

	response := models.UploadResponse{
		Message:    "File uploaded successfully",
		RowCount:   len(rows),
		Headers:    headers,
		SampleData: rows[:sampleSize],
	}

	respondJSON(w, response, http.StatusOK)
}

// Возвращает данные с пагинацией
func (h *Handler) GetData(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 50
	}

	data, err := h.storage.GetPaginated(page, perPage)
	if err != nil {
		respondError(w, "Failed to get data", http.StatusInternalServerError)
		return
	}

	respondJSON(w, data, http.StatusOK)
}

// Ищет данные по запросу
func (h *Handler) SearchData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 50
	}

	data, err := h.storage.Search(query, page, perPage)
	if err != nil {
		respondError(w, "Failed to search data", http.StatusInternalServerError)
		return
	}

	respondJSON(w, data, http.StatusOK)
}

// Очищает все данные
func (h *Handler) ClearData(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Clear(); err != nil {
		respondError(w, "Failed to clear data", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Data cleared successfully"}, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "encoding error", http.StatusInternalServerError)
    }
}

func respondError(w http.ResponseWriter, message string, status int) {
	respondJSON(w, models.ErrorResponse{Error: message}, status)
}