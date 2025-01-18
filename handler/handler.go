package handler

import (
	"encoding/json"
	"github.com/gosimple/slug"
	"net/http"
	"regexp"
	"spt-git-prod.lb.local/gbusuioc/gobootcamp/cache"
	"spt-git-prod.lb.local/gbusuioc/gobootcamp/logger"
)

type Handler struct {
	cache  cacheInterface
	logger logger.Logger
}

type HomeHandler struct{}

// cacheInterface represents my local cache
type cacheInterface interface {
	Add(key string, value int) error
	Get(key string) (int, error)
	List() (map[string]int, error)
	Update(name string, value int) error
	Remove(name string) error
}

func NewHandler(c cacheInterface, l logger.Logger) *Handler {
	return &Handler{
		cache:  c,
		logger: l,
	}
}

var (
	CacheRe       = regexp.MustCompile(`^/cache/*$`)
	CacheReWithID = regexp.MustCompile(`^/cache/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheRe := CacheRe.MatchString(r.URL.Path)
	cacheReWithId := CacheReWithID.MatchString(r.URL.Path)

	switch {
	case r.Method == http.MethodPost && cacheRe:
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && cacheRe:
		h.List(w, r)
		return
	case r.Method == http.MethodGet && cacheReWithId:
		h.Get(w, r)
		return
	case r.Method == http.MethodPut && cacheReWithId:
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && cacheReWithId:
		h.Delete(w, r)
		return
	default:
		return
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page for bet365 Go Training Exercise"))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.logger

	// cache value that will be populated from JSON payload
	var cacheItem cache.Item

	if err := json.NewDecoder(r.Body).Decode(&cacheItem); err != nil {
		InternalServerErrorHandler(w, r, l)
		return
	}

	// convert the key into URL friendly string
	resourceID := slug.Make(cacheItem.Key)

	//call the cache and add the key - value pair
	if err := h.cache.Add(resourceID, cacheItem.Value); err != nil {
		InternalServerErrorHandler(w, r, l)
		return
	}

	// set the status code to 200
	w.WriteHeader(http.StatusOK)
	l.Info("Create StatusOK")
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	l := h.logger
	cacheList, err := h.cache.List()

	jsonBytes, err := json.Marshal(cacheList)
	if err != nil {
		InternalServerErrorHandler(w, r, l)
		return
	}

	w.WriteHeader(http.StatusOK)
	l.Info("List StatusOK")
	w.Write(jsonBytes)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	l := h.logger
	matches := CacheReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, l)
		return
	}

	cacheItem, err := h.cache.Get(matches[1])
	if err != nil {
		if err == cache.NotFoundErr {
			NotFoundHandler(w, r, l)
			return
		}

		InternalServerErrorHandler(w, r, l)
		return
	}

	jsonBytes, err := json.Marshal(cacheItem)
	if err != nil {
		InternalServerErrorHandler(w, r, l)
		return
	}

	w.WriteHeader(http.StatusOK)
	l.Info("Get StatusOK")
	w.Write(jsonBytes)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	l := h.logger
	matches := CacheReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, l)
		return
	}

	// cacheItem object that will be populated from JSON payload
	var cacheItem cache.Item
	if err := json.NewDecoder(r.Body).Decode(&cacheItem); err != nil {
		InternalServerErrorHandler(w, r, l)
		return
	}

	if err := h.cache.Update(matches[1], cacheItem.Value); err != nil {
		if err == cache.NotFoundErr {
			NotFoundHandler(w, r, l)
			return
		}
		InternalServerErrorHandler(w, r, l)
		return
	}
	l.Info("Update StatusOK")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	l := h.logger
	matches := CacheReWithID.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, l)
		return
	}

	if err := h.cache.Remove(matches[1]); err != nil {
		InternalServerErrorHandler(w, r, l)
		return
	}
	l.Info("Delete StatusOK")
	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request, l logger.Logger) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
	l.Error("500 Internal Server Error")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request, l logger.Logger) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
	l.Error("404 Not Found")
}
