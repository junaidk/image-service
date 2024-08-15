package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) linkRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/create", s.createLinkHandler)
	return r
}

func (s *Server) createLinkHandler(w http.ResponseWriter, r *http.Request) {

	expirationTimeStr := r.URL.Query().Get("expiration_time")
	expirationTime, err := time.ParseDuration(expirationTimeStr)
	if err != nil {
		s.badRequestResponse(w, r, fmt.Errorf("invalid expiration time"))
		return
	}

	expToken, _ := s.tokenManger.Create(expirationTime, "")

	uploadLink := fmt.Sprintf("http://%s/v1/image/upload/%s", s.Addr, expToken)

	resp := envelope{
		"link": uploadLink,
	}
	err = s.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}
}
