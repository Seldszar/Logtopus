package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var (
	loggers = make(map[string]*slog.Logger)
)

func getLogger(name string) (*slog.Logger, error) {
	if logger, ok := loggers[name]; ok {
		return logger, nil
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("logs/%s.log", name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, file), nil))
	loggers[name] = logger

	return logger, nil
}

func reply(w http.ResponseWriter, statusCode int, v any) {
	if v != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)

	if v == nil {
		return
	}

	json.NewEncoder(w).Encode(v)
}

func main() {
	addr := flag.String("address", ":3000", "network address to listen")
	mux := http.NewServeMux()

	mux.HandleFunc("POST /{name}", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Level   slog.Level `json:"level"`
			Message string     `json:"message"`
			Detail  any        `json:"detail"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			reply(w, http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})

			return
		}

		logger, err := getLogger(r.PathValue("name"))

		if err != nil {
			reply(w, http.StatusInternalServerError, map[string]any{
				"error": err.Error(),
			})

			return
		}

		logger.Log(r.Context(), slog.Level(data.Level), data.Message, slog.Any("detail", data.Detail))

		reply(w, http.StatusCreated, map[string]any{
			"message": data,
		})
	})

	if err := http.ListenAndServe(*addr, mux); err != nil {
		slog.Error("server failed to listen", slog.Any("error", err))
	}
}
