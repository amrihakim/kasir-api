package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func extractIDFromPath(r *http.Request, prefix string) (int, error) {
    idStr := strings.TrimPrefix(r.URL.Path, prefix)
    idStr = strings.Trim(idStr, "/")               // hilangkan trailing slash jika ada
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