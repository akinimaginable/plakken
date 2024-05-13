package health

import (
	"io"
	"log"
	"net/http"

	"git.gnous.eu/gnouseu/plakken/internal/database"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB *redis.Client
}

func (config Config) Health(w http.ResponseWriter, _ *http.Request) {
	err := database.Ping(config.DB)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := io.WriteString(w, "Redis connection has failed")
		if err != nil {
			log.Println(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, "ok")
	if err != nil {
		log.Println(err)
	}
}
