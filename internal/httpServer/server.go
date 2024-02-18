package httpServer

import (
	"log"
	"net/http"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
	"git.gnous.eu/gnouseu/plakken/internal/web/plak"
	"git.gnous.eu/gnouseu/plakken/internal/web/static"
	"github.com/redis/go-redis/v9"
)

type ServerConfig struct {
	HTTPServer *http.Server
	UrlLength  uint8
	DB         *redis.Client
}

// Configure HTTP router
func (config ServerConfig) router() {
	WebConfig := plak.WebConfig{
		DB:        config.DB,
		UrlLength: config.UrlLength,
	}

	http.HandleFunc("GET /{$}", static.Home)
	http.HandleFunc("GET /{key}/{settings...}", WebConfig.View)
	http.HandleFunc("GET /static/{file}", static.ServeStatic)
	http.HandleFunc("POST /{$}", WebConfig.Create)
	http.HandleFunc("DELETE /{key}", WebConfig.Delete)
}

// Config Configure HTTP server
func Config(listenAddress string) *http.Server {
	server := &http.Server{
		Addr:         listenAddress,
		ReadTimeout:  constant.HTTPTimeout,
		WriteTimeout: constant.HTTPTimeout,
	}

	return server
}

// Server Start HTTP server
func (config ServerConfig) Server() {
	log.Println("Listening on " + config.HTTPServer.Addr)

	config.router()

	log.Fatal(config.HTTPServer.ListenAndServe())
}
