package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Mock store implementing URLStore ---

type mockStore struct {
	createURLFn         func(string) (*URL, error)
	getURLByShortCodeFn func(string) (*URL, error)
	updateURLFn         func(string, string) (*URL, error)
	deleteURLFn         func(string) (bool, error)
	getURLStatsFn       func(string) (*URLStats, error)
	getAllURLsFn        func() ([]URL, error)
	createUserFn        func(string, string) (int, error)
	getUserByUsernameFn func(string) (int, string, error)
	recordAccessFn      func(string) error
}

func (m *mockStore) createURL(url string) (*URL, error)       { return m.createURLFn(url) }
func (m *mockStore) getURLByShortCode(c string) (*URL, error) { return m.getURLByShortCodeFn(c) }
func (m *mockStore) updateURL(c, u string) (*URL, error)      { return m.updateURLFn(c, u) }
func (m *mockStore) deleteURL(c string) (bool, error)         { return m.deleteURLFn(c) }
func (m *mockStore) getURLStats(c string) (*URLStats, error)  { return m.getURLStatsFn(c) }
func (m *mockStore) getAllURLs() ([]URL, error)               { return m.getAllURLsFn() }
func (m *mockStore) createUser(u, p string) (int, error)      { return m.createUserFn(u, p) }
func (m *mockStore) getUserByUsername(u string) (int, string, error) {
	return m.getUserByUsernameFn(u)
}
func (m *mockStore) recordAccess(c string) error { return m.recordAccessFn(c) }

type mockCache struct{}

func (m *mockCache) getCachedURL(shortCode string) (string, error) { return "", errors.New("miss") }
func (m *mockCache) setCachedURL(shortCode, url string) error      { return nil }
func (m *mockCache) deleteCachedURL(shortCode string) error        { return nil }

// --- Test setup helpers ---

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/shorten", shortenURL)
	r.GET("/shorten", listURLs)
	r.GET("/shorten/:code", getShortURL)
	r.PUT("/shorten/:code", updateShortURL)
	r.DELETE("/shorten/:code", deleteShortURL)
	r.GET("/shorten/:code/stats", getURLStatsHandler)
	r.GET("/:code", redirectURL)
	return r
}

func doRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHandlers(t *testing.T) {
	sampleURL := &URL{
		ID:        1,
		URL:       "https://example.com",
		ShortCode: "abc123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cache = &mockCache{}

	// ---------- shortenURL ----------

	t.Run("shortenURL - valid request returns 201", func(t *testing.T) {
		store = &mockStore{
			createURLFn: func(url string) (*URL, error) { return sampleURL, nil },
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/shorten", map[string]string{"url": "https://example.com"})

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("shortenURL - invalid URL returns 400", func(t *testing.T) {
		store = &mockStore{}
		r := setupRouter()
		w := doRequest(r, "POST", "/shorten", map[string]string{"url": "not-a-url"})

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("shortenURL - db error returns 500", func(t *testing.T) {
		store = &mockStore{
			createURLFn: func(url string) (*URL, error) { return nil, errors.New("db error") },
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/shorten", map[string]string{"url": "https://example.com"})

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})

	// ---------- getShortURL ----------

	t.Run("getShortURL - found returns 200", func(t *testing.T) {
		store = &mockStore{
			getURLByShortCodeFn: func(code string) (*URL, error) { return sampleURL, nil },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/shorten/abc123", nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("getShortURL - not found returns 404", func(t *testing.T) {
		store = &mockStore{
			getURLByShortCodeFn: func(code string) (*URL, error) { return nil, nil },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/shorten/notfound", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// ---------- updateShortURL ----------

	t.Run("updateShortURL - valid update returns 200", func(t *testing.T) {
		store = &mockStore{
			updateURLFn: func(code, url string) (*URL, error) { return sampleURL, nil },
		}
		r := setupRouter()
		w := doRequest(r, "PUT", "/shorten/abc123", map[string]string{"url": "https://updated.com"})

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("updateShortURL - invalid body returns 400", func(t *testing.T) {
		store = &mockStore{}
		r := setupRouter()
		w := doRequest(r, "PUT", "/shorten/abc123", map[string]string{"url": "bad-url"})

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("updateShortURL - not found returns 404", func(t *testing.T) {
		store = &mockStore{
			updateURLFn: func(code, url string) (*URL, error) { return nil, nil },
		}
		r := setupRouter()
		w := doRequest(r, "PUT", "/shorten/notfound", map[string]string{"url": "https://example.com"})

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// ---------- deleteShortURL ----------

	t.Run("deleteShortURL - success returns 204", func(t *testing.T) {
		store = &mockStore{
			deleteURLFn: func(code string) (bool, error) { return true, nil },
		}
		r := setupRouter()
		w := doRequest(r, "DELETE", "/shorten/abc123", nil)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("deleteShortURL - not found returns 404", func(t *testing.T) {
		store = &mockStore{
			deleteURLFn: func(code string) (bool, error) { return false, nil },
		}
		r := setupRouter()
		w := doRequest(r, "DELETE", "/shorten/notfound", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// ---------- getURLStatsHandler ----------

	t.Run("getURLStatsHandler - found returns 200", func(t *testing.T) {
		store = &mockStore{
			getURLStatsFn: func(code string) (*URLStats, error) {
				return &URLStats{ID: 1, URL: "https://example.com", ShortCode: "abc123", AccessCount: 5}, nil
			},
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/shorten/abc123/stats", nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("getURLStatsHandler - not found returns 404", func(t *testing.T) {
		store = &mockStore{
			getURLStatsFn: func(code string) (*URLStats, error) { return nil, nil },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/shorten/notfound/stats", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	// ---------- listURLs ----------

	t.Run("listURLs - returns 200 with list", func(t *testing.T) {
		store = &mockStore{
			getAllURLsFn: func() ([]URL, error) { return []URL{*sampleURL}, nil },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/shorten", nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("listURLs - db error returns 500", func(t *testing.T) {
		store = &mockStore{
			getAllURLsFn: func() ([]URL, error) { return nil, errors.New("db error") },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/shorten", nil)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})

	// ---------- register ----------

	t.Run("register - valid request returns 201", func(t *testing.T) {
		store = &mockStore{
			createUserFn: func(username, hash string) (int, error) { return 1, nil },
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/register", map[string]string{"username": "ashish", "password": "password123"})

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("register - missing fields returns 400", func(t *testing.T) {
		store = &mockStore{}
		r := setupRouter()
		w := doRequest(r, "POST", "/register", map[string]string{"username": "ashish"})

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("register - username exists returns 400", func(t *testing.T) {
		store = &mockStore{
			createUserFn: func(username, hash string) (int, error) { return 0, errors.New("duplicate") },
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/register", map[string]string{"username": "ashish", "password": "password123"})

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	// ---------- login ----------

	t.Run("login - valid credentials returns 200", func(t *testing.T) {
		hash, _ := hashPassword("password123")
		store = &mockStore{
			getUserByUsernameFn: func(username string) (int, string, error) { return 1, hash, nil },
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/login", map[string]string{"username": "ashish", "password": "password123"})

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("login - wrong password returns 401", func(t *testing.T) {
		hash, _ := hashPassword("correctpassword")
		store = &mockStore{
			getUserByUsernameFn: func(username string) (int, string, error) { return 1, hash, nil },
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/login", map[string]string{"username": "ashish", "password": "wrongpassword"})

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("login - user not found returns 401", func(t *testing.T) {
		store = &mockStore{
			getUserByUsernameFn: func(username string) (int, string, error) {
				return 0, "", errors.New("not found")
			},
		}
		r := setupRouter()
		w := doRequest(r, "POST", "/login", map[string]string{"username": "ghost", "password": "password123"})

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	// ---------- redirectURL ----------

	t.Run("redirectURL - found redirects 302", func(t *testing.T) {
		store = &mockStore{
			getURLByShortCodeFn: func(code string) (*URL, error) { return sampleURL, nil },
			recordAccessFn:      func(code string) error { return nil },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/abc123", nil)

		if w.Code != http.StatusFound {
			t.Errorf("expected 302, got %d", w.Code)
		}
	})

	t.Run("redirectURL - not found returns 404", func(t *testing.T) {
		store = &mockStore{
			getURLByShortCodeFn: func(code string) (*URL, error) { return nil, nil },
		}
		r := setupRouter()
		w := doRequest(r, "GET", "/notfound", nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}
