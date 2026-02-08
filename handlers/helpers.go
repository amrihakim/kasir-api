package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func extractIDFromPath(r *http.Request, prefix string) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	idStr = strings.Trim(idStr, "/") // hilangkan trailing slash jika ada
	if idStr == "" {
		return 0, fmt.Errorf("missing id in path")
	}

	if idx := strings.Index(idStr, "/"); idx != -1 {
		idStr = idStr[:idx]
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}
	return id, nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
