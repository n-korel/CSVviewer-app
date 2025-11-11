package models

type CSVData struct {
	Headers []string                 `json:"headers"`
	Rows    []map[string]interface{} `json:"rows"`
	Total   int                      `json:"total"`
}

type PaginatedResponse struct {
	Data    []map[string]interface{} `json:"data"`
	Total   int                      `json:"total"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
	Headers []string                 `json:"headers"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UploadResponse struct {
	Message    string                   `json:"message"`
	RowCount   int                      `json:"row_count"`
	Headers    []string                 `json:"headers"`
	SampleData []map[string]interface{} `json:"sample_data"`
}