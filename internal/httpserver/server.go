package httpserver

import (
	"embed"
	"log"
	"net/http"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
	"git.gnous.eu/gnouseu/plakken/internal/web/health"
	"git.gnous.eu/gnouseu/plakken/internal/web/plak"
	"github.com/redis/go-redis/v9"
)

type ServerConfig struct {
	HTTPServer *http.Server
	URLLength  uint8
	DB         *redis.Client
	Static     embed.FS
	Templates  embed.FS
}

func (config ServerConfig) home(w http.ResponseWriter, _ *http.Request) {
	index, err := config.Static.ReadFile("static/index.html")
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(index)
	if err != nil {
		log.Println(err)
	}
}

// Configure HTTP router.
func (config ServerConfig) router() {
	WebConfig := plak.WebConfig{
		DB:        config.DB,
		URLLength: config.URLLength,
		Templates: config.Templates,
	}
	staticFiles := http.FS(config.Static)

	http.HandleFunc("GET /{$}", config.home)
	http.Handle("GET /static/{file}", http.FileServer(staticFiles))
	http.HandleFunc("GET /health/", health.Config{DB: config.DB}.Health)
	http.HandleFunc("GET /{key}/{settings...}", WebConfig.View)
	http.HandleFunc("POST /{$}", WebConfig.CurlCreate)
	http.HandleFunc("POST /create/{$}", WebConfig.PostCreate)
	http.HandleFunc("DELETE /{key}", WebConfig.DeleteRequest)
}

// Config Configure HTTP server.
func Config(listenAddress string) *http.Server {
	server := &http.Server{
		Addr:         listenAddress,
		ReadTimeout:  constant.HTTPTimeout,
		WriteTimeout: constant.HTTPTimeout,
	}

	return server
}

// Server Start HTTP server.
func (config ServerConfig) Server() {
	log.Println("Listening on " + config.HTTPServer.Addr)

	config.router()

	log.Fatal(config.HTTPServer.ListenAndServe())
}
