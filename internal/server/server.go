package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"go-bug-analysis/internal/analysis"
	"go-bug-analysis/internal/csvparse"
)

// AppState holds the application-level shared state (parsed CSV data).
type AppState struct {
	mu     sync.RWMutex
	Result *csvparse.ParseResult
}

// SetResult stores a parse result (thread-safe).
func (s *AppState) SetResult(r *csvparse.ParseResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Result = r
}

// GetResult retrieves the current parse result (thread-safe).
func (s *AppState) GetResult() *csvparse.ParseResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Result
}

// New creates an HTTP handler with static file serving and API routes.
// staticFS should be the embedded filesystem with the "static" prefix already stripped.
func New(staticFS fs.FS, state *AppState) http.Handler {
	mux := http.NewServeMux()

	// Static file server — serves HTML/CSS/JS/ECharts from embedded FS
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/", fileServer)

	// API routes
	mux.HandleFunc("POST /api/upload", handleUpload(state))
	mux.HandleFunc("GET /api/data", handleData(state))
	mux.HandleFunc("GET /api/analysis", handleAnalysis(state))

	return mux
}

// uploadResponse is the JSON structure returned by the upload endpoint.
type uploadResponse struct {
	Success bool            `json:"success"`
	Summary *uploadSummary  `json:"summary,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type uploadSummary struct {
	TotalBugs int              `json:"totalBugs"`
	Columns   []string         `json:"columns"`
	Warnings  []string         `json:"warnings"`
	SampleBug *sampleBugInfo   `json:"sampleBug,omitempty"`
}

type sampleBugInfo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Creator  string `json:"creator"`
}

// handleUpload processes multipart CSV file uploads.
func handleUpload(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Limit upload size to 100MB
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			writeJSON(w, http.StatusBadRequest, uploadResponse{
				Error: "文件解析失败，请确认文件大小不超过 100MB",
			})
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, uploadResponse{
				Error: "请选择要上传的文件",
			})
			return
		}
		defer file.Close()

		// Validate file extension
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".csv" {
			writeJSON(w, http.StatusBadRequest, uploadResponse{
				Error: "请上传 CSV 文件",
			})
			return
		}

		// Parse CSV
		result, err := csvparse.Parse(file)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, uploadResponse{
				Error: err.Error(),
			})
			return
		}

		// Store result in application state
		state.SetResult(result)

		// Build response summary
		summary := &uploadSummary{
			TotalBugs: result.TotalRows,
			Columns:   result.Columns,
			Warnings:  result.Warnings,
		}

		// Add sample bug if available
		if len(result.Bugs) > 0 {
			first := result.Bugs[0]
			summary.SampleBug = &sampleBugInfo{
				ID:       first.ID,
				Title:    first.Title,
				Status:   first.Status,
				Severity: first.Severity,
				Creator:  first.Creator,
			}
		}

		writeJSON(w, http.StatusOK, uploadResponse{
			Success: true,
			Summary: summary,
		})
	}
}

// handleData returns the full parsed data (for Phase 2 analysis modules).
func handleData(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		result := state.GetResult()
		if result == nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "请先上传 CSV 文件",
			})
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}

// handleAnalysis runs analysis on the stored bug data and returns the results.
func handleAnalysis(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		result := state.GetResult()
		if result == nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "请先上传 CSV 文件",
			})
			return
		}

		analysisResult := analysis.Analyze(result.Bugs)
		writeJSON(w, http.StatusOK, analysisResult)
	}
}

// writeJSON encodes v as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
