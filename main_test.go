package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"spt-git-prod.lb.local/gbusuioc/gobootcamp/cache"
	"spt-git-prod.lb.local/gbusuioc/gobootcamp/config"
	"spt-git-prod.lb.local/gbusuioc/gobootcamp/handler"
	"spt-git-prod.lb.local/gbusuioc/gobootcamp/logger"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Errorf("Could not read %v", name)
	}

	return content
}

func TestRecipesHandlerCRUD_Integration(t *testing.T) {

	// Create a MemStore and Recipe Handler
	c := config.NewConfig()
	l := logger.NewLogger(c)
	ch := cache.NewCache()
	h := handler.NewHandler(ch, *l)

	// Testdata
	cacheItem1 := readTestData(t, "cacheItem1.json")
	cacheItem1Reader := bytes.NewReader(cacheItem1)

	cacheItem2 := readTestData(t, "cacheItem2.json")
	cacheItem2Reader := bytes.NewReader(cacheItem2)

	// CREATE - add a new cacheItem to the cache
	req := httptest.NewRequest(http.MethodPost, "/cache", cacheItem1Reader)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ := ch.List()
	assert.Len(t, saved, 1)

	//// GET - find the cacheItem we just added
	req = httptest.NewRequest(http.MethodGet, "/cache/key-one", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.JSONEq(t, "1", string(data))

	// UPDATE - add butter to ham and cheese recipe
	req = httptest.NewRequest(http.MethodPut, "/cache/key-one", cacheItem2Reader)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	// it does not update 1 doesn't become 2
	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	updatedCacheItem, err := ch.Get("key-one")
	assert.NoError(t, err)

	assert.Equal(t, updatedCacheItem, 2)

	//DELETE - remove key1 from cache
	req = httptest.NewRequest(http.MethodDelete, "/cache/key-one", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ = ch.List()
	assert.Len(t, saved, 0)
}
