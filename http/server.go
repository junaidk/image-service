package http

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ims "github.com/junaidk/image-service"
	"github.com/junaidk/image-service/internal/token"
)

const AUTH_TOKEN = "secret"

type Server struct {
	ln                net.Listener
	server            *http.Server
	tokenManger       *token.Manager
	Addr              string
	Port              string
	ImageService      ims.ImageService
	StatisticsService ims.StatisticsService
	ImageDir          string
	StaticToken       string
	SigningSecret     string
}

func NewServer() *Server {
	svr := &Server{
		server: &http.Server{},
	}
	svr.tokenManger = token.New(svr.SigningSecret)

	router := chi.NewRouter()
	router.Use(loggingMiddleware(slog.Default()))
	router.Use(middleware.Recoverer)
	router.MethodNotAllowed(svr.methodNotAllowedResponse)

	router.Get("/healthcheck", svr.healthcheckHandler)

	// v1
	v1 := chi.NewRouter()
	router.Mount("/v1", v1)

	v1.Mount("/link", svr.authMiddleware(svr.linkRoutes()))
	v1.Mount("/image", svr.imageRoutes())
	v1.Mount("/statistics", svr.authMiddleware(svr.statisticsRoutes()))

	svr.server.Handler = router

	return svr
}

func (s *Server) Open() (err error) {

	if s.ln, err = net.Listen("tcp", ":"+s.Port); err != nil {
		return err
	}

	slog.Info("Starting server", slog.String("port", s.Port))
	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) Close() (err error) {
	return s.server.Close()
}
