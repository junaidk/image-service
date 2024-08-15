package http

import "net/http"

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := envelope{
		"status": "available",
	}
	err := s.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}
}
