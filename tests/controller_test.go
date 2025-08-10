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
	"github.com/Balaji01-4D/my-dear-bug/internals/tag"
	"github.com/Balaji01-4D/my-dear-bug/internals/upvote"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
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

	_ = os.Setenv("ADMIN_USERNAME", "admin")
	_ = os.Setenv("ADMIN_PASSWORD", "secret")
	_ = os.Setenv("TEST_DATABASE_URL", "postgres://tester:test@123@localhost:5432/test_bug?sslmode=disable")
	// Do NOT force a Postgres DSN here; allow sqlite in-memory by default for portability.
	// Users can export TEST_DATABASE_URL externally to run tests against Postgres.

	var db *gorm.DB
	var err error

	if dsn := os.Getenv("TEST_DATABASE_URL"); dsn != "" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open postgres test db: %v", err)
		}
		t.Cleanup(func() {
			db.Exec("TRUNCATE TABLE confessions RESTART IDENTITY CASCADE")
			db.Exec("TRUNCATE TABLE tags RESTART IDENTITY CASCADE")
			db.Exec("TRUNCATE TABLE confession_tags RESTART IDENTITY CASCADE")
			db.Exec("TRUNCATE TABLE upvotes RESTART IDENTITY CASCADE")
		})
	} else {
			// Fallback to in-memory sqlite
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open sqlite test db: %v", err)
		}
	}

	if err := db.AutoMigrate(&confpkg.Confession{}, &tag.Tag{}, &upvote.Upvote{}); err != nil {
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
	req.RemoteAddr = nextRemoteAddr()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createConfession(t *testing.T, r *gin.Engine, title, desc, lang string, tags []string) uint {
	t.Helper()
	payload := map[string]any{
		"title":       title,
		"description": desc,
		"language":    lang,
		"tags":        tags,
		"snippet":     "",
		"isFlagged":   false,
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create confession (%s): status=%d body=%s", title, w.Code, w.Body.String())
	}
	var resp struct {
		ID uint `json:"id"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.ID
}

/* ---------- Tests (deduplicated and consolidated) ---------- */

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
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "not found" {
		t.Fatalf("expected error 'not found', got %v", resp["error"])
	}
}

func TestGetConfession_InvalidID(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/abc", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-numeric id, got %d", w.Code)
	}
	w = doJSONRequest(r, http.MethodGet, "/confessions/-1", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for negative id, got %d", w.Code)
	}
}

func TestCreateConfession_ValidationError(t *testing.T) {
	r, _ := setupRouter(t)
	payload := map[string]any{"title": "shrt"}
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
		"tags":        []string{"concurrency", "Concurrency", "mutex"},
		"isFlagged":   false,
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body=%s", w.Code, w.Body.String())
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
	if len(created.Tags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d", len(created.Tags))
	}
	w = doJSONRequest(r, http.MethodGet, "/confessions/"+jsonNumber(created.ID), nil)
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
	payload := map[string]any{
		"title":       "Race Condition",
		"description": "Non-deterministic writes",
		"language":    "go",
		"tags":        []string{"concurrency", "race condition"},
		"snippet":     "go func() { x++; y++ }(); go func() { x--; y-- }();",
		"isFlagged":   false,
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("create failed: status %d body=%s", w.Code, w.Body.String())
	}
	var created struct{ ID uint `json:"id"` }
	_ = json.Unmarshal(w.Body.Bytes(), &created)
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
	w = doJSONRequest(r, http.MethodGet, "/confessions/"+jsonNumber(created.ID), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}

func TestSearchConfessions_Empty(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/search?q=unavailabe", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Fatalf("expected empty list, got %d items", len(resp))
	}
}

func TestSearchConfessions_Found(t *testing.T) {
	r, _ := setupRouter(t)
	payload := map[string]any{
		"title":       "Memory Leak",
		"description": "Forgetting to free memory in C",
		"language":    "c",
		"tags":        []string{"memory", "leak"},
		"snippet":     "malloc(100); // forgot to free",
		"isFlagged":   false,
	}
	w := doJSONRequest(r, http.MethodPost, "/confessions", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	w = doJSONRequest(r, http.MethodGet, "/confessions/search?q=memory", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) == 0 {
		t.Fatal("expected at least one confession matching 'memory'")
	}
	for _, c := range resp {
		if c["title"] == nil || c["description"] == nil {
			t.Fatal("expected title and description in search results")
		}
	}
}

func TestSearchConfessions_InvalidQuery(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/search?q=", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty query, got %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "query parameters required" {
		t.Fatalf("expected error 'query parameters required', got %v", resp["error"])
	}
}

func TestGetByLanguage_Filtering(t *testing.T) {
	r, _ := setupRouter(t)
	_ = doJSONRequest(r, http.MethodPost, "/confessions", map[string]any{
		"title":       "panic in go",
		"description": "index out of range",
		"language":    "go",
		"tags":        []string{"go", "panic"},
		"snippet":     "panic(\"index out of range\")",
		"isFlagged":   false,
	})
	_ = doJSONRequest(r, http.MethodPost, "/confessions", map[string]any{
		"title":       "null pointer",
		"description": "undefined is not an object",
		"language":    "javascript",
		"tags":        []string{"http", "headers"},
		"snippet":     "console.log(undefined.property)",
		"isFlagged":   false,
	})
	w := doJSONRequest(r, http.MethodGet, "/confessions/language/go?limit=10&offset=0", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &list)
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
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) != 0 {
		t.Fatalf("expected empty list, got %d items", len(resp))
	}
}

func TestRandom_NotFoundThenFound(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/random", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for empty random, got %d", w.Code)
	}
	createConfession(t, r, "One Random", "Just a random confession", "go", []string{"random"})
	w = doJSONRequest(r, http.MethodGet, "/confessions/random", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 after creating confession, got %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["id"] == nil {
		t.Fatalf("expected id in random response")
	}
}

func TestTrendingWeekly_Empty(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/trending/weekly?limit=5&offset=0", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Fatalf("expected empty trending weekly list")
	}
}

func TestTrendingMonthly_Empty(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/trending/monthly?limit=5&offset=0", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Fatalf("expected empty trending monthly list")
	}
}

func TestHallOfFame_Empty(t *testing.T) {
	r, _ := setupRouter(t)
	w := doJSONRequest(r, http.MethodGet, "/confessions/hall-of-fame?limit=5&offset=0", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Fatalf("expected empty hall of fame list")
	}
}

func TestSearchConfessions_ByTagOnly(t *testing.T) {
	r, _ := setupRouter(t)
	createConfession(t, r, "Tag Only", "Uses special tag", "python", []string{"unique_tag"})
	w := doJSONRequest(r, http.MethodGet, "/confessions/search?tag=unique_tag", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) == 0 {
		t.Fatalf("expected at least one result for tag search")
	}
}

func TestSearchConfessions_ByLanguageOnly(t *testing.T) {
	r, _ := setupRouter(t)
	createConfession(t, r, "Lang Only", "Lang search test", "rust", []string{"ownership"})
	w := doJSONRequest(r, http.MethodGet, "/confessions/search?language=rust", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	found := false
	for _, c := range list {
		if c["language"] == "rust" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected at least one rust confession")
	}
}

func TestListConfessions_Pagination(t *testing.T) {
	r, _ := setupRouter(t)
	id1 := createConfession(t, r, "Paginate 1", "pagination test description one", "java", []string{"p1"})
	id2 := createConfession(t, r, "Paginate 2", "pagination test description two", "go", []string{"p2"})
	id3 := createConfession(t, r, "Paginate 3", "pagination test description three", "go", []string{"p3"})
	if id1 == id2 || id2 == id3 {
		t.Fatalf("ids should differ")
	}
	w := doJSONRequest(r, http.MethodGet, "/confessions?limit=2&offset=0", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var firstPage []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &firstPage)
	if len(firstPage) != 2 {
		t.Fatalf("expected 2 items on first page, got %d", len(firstPage))
	}
	w = doJSONRequest(r, http.MethodGet, "/confessions?limit=2&offset=2", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var secondPage []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &secondPage)
	if len(secondPage) != 1 {
		t.Fatalf("expected 1 item on second page, got %d", len(secondPage))
	}
}

/* ---------- Helpers ---------- */

// jsonNumber formats a uint as a JSON-friendly string for URL paths
func jsonNumber(id uint) string {
	b, _ := json.Marshal(id)
	return string(bytes.Trim(b, "\""))
}
