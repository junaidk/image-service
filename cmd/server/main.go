package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/junaidk/image-service/http"

	"github.com/junaidk/image-service/postgress"
)

type config struct {
	db struct {
		dsn string
	}
	http struct {
		addr          string
		port          string
		staticToken   string
		signingSecret string
		imageDir      string
	}
}

type application struct {
	config     config
	DB         *postgress.DB
	HTTPServer *http.Server
}

func main() {

	var app = application{}

	app.ParseFlags()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := app.Run(); err != nil {
		slog.Error("error starting", "msg", err)
		app.Close()
	}

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.Close()
		cancel()
	}()

	<-ctx.Done()
}

func (a *application) ParseFlags() error {
	var cfg config
	flag.StringVar(&cfg.http.port, "port", lookupEnvOrString("HTTP_PORT", "8080"), "API server port")
	flag.StringVar(&cfg.http.addr, "addr", lookupEnvOrString("HTTP_ADDR", "localhost:8080"), "API server address exposed to API user")
	flag.StringVar(&cfg.http.staticToken, "token", lookupEnvOrString("HTTP_TOKEN", "secret"), "API server token for static auth")
	flag.StringVar(&cfg.http.signingSecret, "secret", lookupEnvOrString("HTTP_SECRET", "my-secret"), "Secret used to sign upload urls")
	flag.StringVar(&cfg.http.imageDir, "image-dir", lookupEnvOrString("HTTP_IMAGE_DIR", "/home/junaid/documents/gowork/samples/image-api/image-data"), "Dir to store images")
	flag.StringVar(&cfg.db.dsn, "dsn", lookupEnvOrString("DB_DSN", "postgres://app-user:secret@localhost:5432/app_db?sslmode=disable"), "DSN")

	a.config = cfg

	return nil
}

func (a *application) Run() error {

	db := postgress.NewDB(a.config.db.dsn)
	imageService := postgress.NewImageService(db)
	statisticsService := postgress.NewStatisticsService(db)

	err := db.Open()
	if err != nil {
		panic(err)
	}

	server := http.NewServer()

	server.Port = a.config.http.port
	server.ImageDir = a.config.http.imageDir
	server.Addr = a.config.http.addr
	server.StaticToken = a.config.http.staticToken
	server.ImageService = imageService
	server.StatisticsService = statisticsService

	return server.Open()
}

func (a *application) Close() error {
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
