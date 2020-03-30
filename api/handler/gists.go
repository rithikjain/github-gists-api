package handler

import (
	"encoding/json"
	"github.com/rithikjain/GistsBackend/api/middleware"
	"github.com/rithikjain/GistsBackend/api/view"
	"github.com/rithikjain/GistsBackend/pkg/gists"
	"net/http"
)

// Protected Request
func viewAllFiles(s gists.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		files, err := s.ViewAllFiles(claims["id"].(float64))
		if err != nil {
			view.Wrap(err, w)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Files Retrieved",
			"status":  http.StatusOK,
			"files":   files,
		})
	})
}

// Protected Request
func createGist(s gists.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		gistFile := &gists.CreateGist{}
		_ = json.NewDecoder(r.Body).Decode(gistFile)

		files, err := s.CreateGist(claims["id"].(float64), gistFile)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		if len(*files) != 0 {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "File Created",
				"status":  http.StatusCreated,
				"files":   files,
			})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "File Creation Failed",
				"status":  http.StatusBadRequest,
			})
		}
	})
}

func MakGistsHandler(r *http.ServeMux, svc gists.Service) {
	r.Handle("/api/gists/view", middleware.Validate(viewAllFiles(svc)))
	r.Handle("/api/gists/create", middleware.Validate(createGist(svc)))
}
