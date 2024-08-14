package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

// logError method is a generic helper for logging an error message in *Server, as well
// as the requested method and request URL.
func (s *Server) logError(r *http.Request, err error) {
	slog.Error(err.Error(), "attr",
		map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
		})
}

// errorResponse method is a generic helper for sending JSON-formatted error messages to the
// client with a given status code. Note that we're using an interface{} type for the message
// parameter, rather than just a string type, as this gives us more flexibility over the values
// that we can include in the response.
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

// serverErrorResponse method is used when our Server encounters an unexpected problem
// at runtime. it logs the detailed error message, then uses the errorResponse() helper to send a
// 500 Internal Server Error status code and JSON response (containing the generic error message)
// to the client
func (s *Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	s.errorResponse(w, r, 500, message)
}

// notFoundResponse method is used to send a 404 Not Found status code and JSON response to the
// client.
func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	s.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse method is used to send a 405 Method Not Allowed status code and
// JSON response to the client.
func (s *Server) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	s.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse sends JSON-formatted error message with 400 Bad Request status code.
func (s *Server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (s *Server) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	s.errorResponse(w, r, http.StatusUnauthorized, message)
}
