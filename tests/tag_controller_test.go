
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

var reqCounterTag uint64

func nextRemoteAddrTag() string {
	id := atomic.AddUint64(&reqCounterTag, 1)
	return fmt.Sprintf("127.0.0.%d:12345", (id%200)+1)
}

func setupRouterTag(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	_ = os.Setenv("ADMIN_USERNAME", "admin")
	_ = os.Setenv("ADMIN_PASSWORD", "secret")
	_ = os.Setenv("TEST_DATABASE_URL", "postgres://tester:test@123@localhost:5432/test_bug?sslmode=disable")

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
	tag.RegisterRoutes(r, db)
	return r, db
}

func doJSONRequestTag(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = nextRemoteAddrTag()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}


/* ---------- Tag Route Tests ---------- */

func TestTags_ListEmpty(t *testing.T) {
	r, _ := setupRouterTag(t)
	w := doJSONRequestTag(r, http.MethodGet, "/tags", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var tags []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &tags); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected empty list, got %d", len(tags))
	}
}

func TestTags_CreateValidationError(t *testing.T) {
	r, _ := setupRouterTag(t)
	w := doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": ""})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestTags_CreateAndList(t *testing.T) {
	r, _ := setupRouterTag(t)
	w := doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": "concurrency"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on create, got %d (%s)", w.Code, w.Body.String())
	}
	w = doJSONRequestTag(r, http.MethodGet, "/tags", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on list, got %d", w.Code)
	}
	var tags []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &tags)
	found := false
	for _, tg := range tags {
		if tg["name"] == "concurrency" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created tag not found in list")
	}
}

func TestTags_Suggest(t *testing.T) {
	r, _ := setupRouterTag(t)
	_ = doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": "Memory"})
	_ = doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": "memory leak"})
	_ = doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": "performance"})
	w := doJSONRequestTag(r, http.MethodGet, "/tags/suggest?query=mem", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
	var suggestions []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &suggestions); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatalf("expected at least one suggestion")
	}
	for _, s := range suggestions {
		name, _ := s["name"].(string)
		if name == "" {
			t.Fatalf("suggestion missing name")
		}
	}
}

func TestTags_Suggest_QueryTooShort(t *testing.T) {
	r, _ := setupRouterTag(t)
	w := doJSONRequestTag(r, http.MethodGet, "/tags/suggest?query=", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "Query too short" {
		t.Fatalf("unexpected error message: %v", resp["error"])
	}
}

func TestTags_Delete_Unauthorized(t *testing.T) {
	r, _ := setupRouterTag(t)
	_ = doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": "delete-me"})
	// need id
	w := doJSONRequestTag(r, http.MethodGet, "/tags", nil)
	var tags []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &tags)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	id := jsonNumber(uintFromAny(tags[0]["id"]))
	resp := doJSONRequestTag(r, http.MethodDelete, "/tags/"+id, nil)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 unauthorized, got %d", resp.Code)
	}
}

func TestTags_Delete_Authorized(t *testing.T) {
	r, _ := setupRouterTag(t)
	_ = doJSONRequestTag(r, http.MethodPost, "/tags", map[string]any{"name": "refactor"})
	list := doJSONRequestTag(r, http.MethodGet, "/tags", nil)
	var tagsResp []map[string]any
	_ = json.Unmarshal(list.Body.Bytes(), &tagsResp)
	if len(tagsResp) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tagsResp))
	}
	id := jsonNumber(uintFromAny(tagsResp[0]["id"]))
	req := httptest.NewRequest(http.MethodDelete, "/tags/"+id, nil)
	req.RemoteAddr = nextRemoteAddrTag()
	auth := base64.StdEncoding.EncodeToString([]byte("admin:secret"))
	req.Header.Set("Authorization", "Basic "+auth)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", rec.Code, rec.Body.String())
	}
	var delResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &delResp)
	if delResp["message"] != "tag deleted successfully" {
		t.Fatalf("unexpected delete response: %v", delResp)
	}
	after := doJSONRequestTag(r, http.MethodGet, "/tags", nil)
	var remaining []map[string]any
	_ = json.Unmarshal(after.Body.Bytes(), &remaining)
	if len(remaining) != 0 {
		t.Fatalf("expected no tags after delete, got %d", len(remaining))
	}
}

/* ---------- Helpers ---------- 

func jsonNumber(id uint) string {
	b, _ := json.Marshal(id)
	return string(bytes.Trim(b, "\""))
}

*/

func uintFromAny(v any) uint {
	switch n := v.(type) {
	case float64:
		return uint(n)
	case int:
		return uint(n)
	case uint:
		return n
	case json.Number:
		i, _ := n.Int64()
		return uint(i)
	default:
		return 0
	}
}
