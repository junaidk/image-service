package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type envelope map[string]interface{}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		slog.Error("error writing response", "error", err.Error())
		return err
	}

	return nil
}
