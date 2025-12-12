package response

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func RenderJson(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	b, err := json.Marshal(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		slog.ErrorContext(ctx, "response.RenderJson: http.ResponseWriter.Write:", "err", err)
	}
}
