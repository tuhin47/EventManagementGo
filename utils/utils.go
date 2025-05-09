package utils

import (
	"net/http"
	"strconv"
	"strings"
)

// SetJSONHeader sets the Content-Type header to application/json
func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func ExtractIDFromURL(r *http.Request, prefix string) (string, error) {
	id := strings.TrimPrefix(r.URL.Path, prefix)
	if id == "" {
		return "", http.ErrMissingFile
	}
	return id, nil
}

func GetPaginationParams(r *http.Request) (int, int) {

	// Default values
	page, pageSize := 1, 10

	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	offset := (page - 1) * pageSize

	return pageSize, offset
}
