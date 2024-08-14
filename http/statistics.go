package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) statisticsRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.getStatisticsHandler)
	return r
}

func (s *Server) getStatisticsHandler(w http.ResponseWriter, r *http.Request) {

	stats, err := s.StatisticsService.GetStatistics(r.Context())
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

	resp := envelope{
		"data": stats,
	}
	err = s.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

}
