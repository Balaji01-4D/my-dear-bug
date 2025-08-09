package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	confpkg "github.com/Balaji01-4D/my-dear-bug/internals/confession"
	tagpkg "github.com/Balaji01-4D/my-dear-bug/internals/tag"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var reqCounter uint64

func nextRemoteAddr() string {
	id := atomic.AddUint64(&reqCounter, 1)
	return fmt.Sprintf("127.0.0.%d:12345", (id%200)+1)
}

func setupRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Ensure admin creds are set BEFORE route registration (captured by middleware)
	_ = os.Setenv("ADMIN_USERNAME", "admin")
	_ = os.Setenv("ADMIN_PASSWORD", "secret")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Migrate required models (including tags for many2many join table)
	if err := db.AutoMigrate(&confpkg.Confession{}, &tagpkg.Tag{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	r := gin.New()
	confpkg.RegisterRoutes(r, db)
	return r, db
}

func doJSONRequest(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	// use a unique client IP for each request to avoid rate limiter collisions
	req.RemoteAddr = nextRemoteAddr()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestListConfessions_Empty(t *testing.T) {
	r, _ := setupRouter(t)

	w := doJSONRequest(r, http.MethodGet, "/confessions?offset=0&limit=10", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not a JSON array: %v", err)
	}
	if len(resp) != 0 {
		t.Fatalf("expected empty list, got %d items", len(resp))
	}
}

func TestGetConfession_NotFound(t *testing.T) {
	r, _ := setupRouter(t)

	w := doJSONRequest(r, http.MethodGet, "/confessions/999999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if resp["error"] != "not found" {
		t.Fatalf("expected error 'not found', got %v", resp["error"])
	}
}

func TestCreateConfession_ValidationError(t *testing.T) {
	r, _ := setupRouter(t)

	// Missing required fields like title/description/language
	payload := map[string]any{
		"title": "shrt", // too short (min 5) but present; also omit description and language to trigger 400
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] == nil {
		t.Fatalf("expected validation error message")
	}
}

func TestCreateAndGetConfession_Success(t *testing.T) {
	r, _ := setupRouter(t)

	payload := map[string]any{
		"title":       "Deadlock Detected",
		"description": "Two goroutines wait on each other causing deadlock",
		"snippet":     "mutex1.Lock(); mutex2.Lock();",
		"language":    "go",
		"tags":        []string{"concurrency", "Concurrency", "mutex"}, // dedupe case-insensitive
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var created struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
		Tags  []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"tags"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected created ID > 0")
	}
	if created.Title != "Deadlock Detected" {
		t.Fatalf("unexpected title: %s", created.Title)
	}
	if len(created.Tags) != 2 { // concurrency (deduped) + mutex
		t.Fatalf("expected 2 unique tags, got %d", len(created.Tags))
	}

	// Fetch by ID
	w = doJSONRequest(r, http.MethodGet,
		"/confessions/"+jsonNumber(created.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on GET by id, got %d", w.Code)
	}
}

func TestDeleteConfession_Unauthorized(t *testing.T) {
	r, _ := setupRouter(t)

	w := doJSONRequest(r, http.MethodDelete, "/confessions/1", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestDeleteConfession_Authorized(t *testing.T) {
	r, _ := setupRouter(t)

	// First create a confession to delete
	payload := map[string]any{
		"title":       "Race Condition",
		"description": "Non-deterministic writes",
		"language":    "go",
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusOK {
		t.Fatalf("create failed: status %d body=%s", w.Code, w.Body.String())
	}
	var created struct {
		ID uint `json:"id"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.ID == 0 {
		t.Fatalf("expected created ID > 0")
	}

	// Issue DELETE with Basic Auth
	req := httptest.NewRequest(http.MethodDelete, "/confessions/"+jsonNumber(created.ID), nil)
	req.RemoteAddr = nextRemoteAddr()
	auth := base64.StdEncoding.EncodeToString([]byte("admin:secret"))
	req.Header.Set("Authorization", "Basic "+auth)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200 on delete, got %d (body=%s)", resp.Code, resp.Body.String())
	}
	var delResp map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &delResp)
	if delResp["message"] != "successfully deleted" {
		t.Fatalf("unexpected delete message: %v", delResp["message"])
	}

	// Verify it is gone
	w = doJSONRequest(r, http.MethodGet, "/confessions/"+jsonNumber(created.ID), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}

func TestGetByLanguage_Filtering(t *testing.T) {
	r, _ := setupRouter(t)

	// Create two confessions in different languages
	_ = doJSONRequest(r, http.MethodPost, "/confessions", map[string]any{
		"title":       "panic in go",
		"description": "index out of range",
		"language":    "go",
	})
	_ = doJSONRequest(r, http.MethodPost, "/confessions", map[string]any{
		"title":       "null pointer",
		"description": "undefined is not an object",
		"language":    "javascript",
	})

	w := doJSONRequest(r, http.MethodGet, "/confessions/language/go?limit=10&offset=0", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, c := range list {
		if lang, ok := c["language"].(string); !ok || lang == "" {
			t.Fatalf("missing language in response")
		} else if lang != "go" && lang != "GO" && lang != "Go" {
			t.Fatalf("expected language 'go', got %q", lang)
		}
	}
}

func TestTopConfessions_Empty(t *testing.T) {
	r, _ := setupRouter(t)

	w := doJSONRequest(r, http.MethodGet, "/confessions/top?offset=0&limit=10", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not a JSON array: %v", err)
	}
	if len(resp) != 0 {
		t.Fatalf("expected empty list, got %d items", len(resp))
	}
}

// jsonNumber formats a uint as a JSON-friendly string for URL paths
func jsonNumber(id uint) string {
	b, _ := json.Marshal(id)
	return string(bytes.Trim(b, "\""))
}
