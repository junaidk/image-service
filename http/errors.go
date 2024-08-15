package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (s *Server) logError(r *http.Request, err error) {
	slog.Error(err.Error(), "attr",
		map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
		})
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an error
	// then log it, and fall back to sending the client an empty response with a 500 Internal
	// Server Error status code
	err := s.writeJSON(w, status, env, nil)
	if err != nil {
		s.logError(r, err)
		w.WriteHeader(500)
	}
}

func (s *Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	s.errorResponse(w, r, 500, message)
}

func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	s.errorResponse(w, r, http.StatusNotFound, message)
}

func (s *Server) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	s.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (s *Server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (s *Server) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	s.errorResponse(w, r, http.StatusUnauthorized, message)
}
