package status

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"git.gnous.eu/gnouseu/plakken/internal/database"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB        *redis.Client
	StartTime time.Time
}

type info struct {
	Uptime    time.Duration `json:"uptime"`
	Version   string        `json:"version"`
	GoVersion string        `json:"goVersion"`
	Source    string        `json:"source"`
}

type health struct {
	Status string `json:"status"`
	DB     string `json:"db"` // TODO: struct with ping duration ?
}

func (config Config) Ready(w http.ResponseWriter, _ *http.Request) {
	err := database.Ping(config.DB)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := io.WriteString(w, "ko")
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

func (config Config) Info(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := info{
		Uptime:    uptime(config.StartTime),
		Version:   "nightly", // TODO
		GoVersion: runtime.Version(),
		Source:    "https://git.gnous.eu/gnouseu/plakken",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func (config Config) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var response health
	if database.Ping(config.DB) != nil {
		response.DB = "ko"
	} else {
		response.DB = "ok"
	}

	if response.DB == "ok" {
		response.Status = "ok"
	} else {
		response.Status = "ko"
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func uptime(startTime time.Time) time.Duration {
	return time.Since(startTime)
}
