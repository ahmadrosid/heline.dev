package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var allowedMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
}

var allowedHeaders = []string{
	"Content-Type",
	"Authorization",
	"Accept",
}

type CorsConfig struct {
	allowedOrigin string
}

func respondError(w http.ResponseWriter, status int, err ...error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errors := []string{}
	for _, e := range err {
		errors = append(errors, e.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(Map{
		"status": status,
		"errors": errors,
	})
}

func StrListContains(sources []string, target string) bool {
	for _, item := range sources {
		if item == target {
			return true
		}
	}
	return false
}

func wrapCORSHandler(h http.Handler, config *CorsConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		requestMethod := req.Header.Get("Access-Control-Request-Method")

		if origin == "" {
			h.ServeHTTP(w, req)
			return
		}

		if !isValidOrigin(config, origin) {
			respondError(w, http.StatusForbidden, fmt.Errorf("origin not allowed"))
			return
		}

		if req.Method == http.MethodOptions && !StrListContains(allowedMethods, requestMethod) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")

		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", "300")

			return
		}

		h.ServeHTTP(w, req)
	})
}

func isValidOrigin(config *CorsConfig, origin string) bool {
	if len(config.allowedOrigin) == 0 {
		return false
	}

	if len(config.allowedOrigin) == 1 && (config.allowedOrigin == "*") {
		return true
	}

	return config.allowedOrigin == origin
}
