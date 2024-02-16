package plak

import (
	"context"
	"embed"
	"io"
	"log"
	"net/http"
	"time"

	"git.gnous.eu/gnouseu/plakken/internal/database"
	"git.gnous.eu/gnouseu/plakken/internal/utils"
	"github.com/redis/go-redis/v9"

	"html/template"
)

var ctx = context.Background()

type WebConfig struct {
	DB        *redis.Client
	UrlLength uint8
	Templates embed.FS
}

// Plak "Object" for plak
type Plak struct {
	Key        string
	Content    string
	Expiration time.Duration
	DB         *redis.Client
}

// Create manage POST request for create Plak
func (config WebConfig) Create(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	dbConf := database.DBConfig{
		DB: config.DB,
	}

	secret := utils.GenerateSecret()
	key := utils.GenerateUrl(config.UrlLength)
	rawExpiration := r.FormValue("exp")
	expiration, err := utils.ParseExpiration(rawExpiration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else if expiration == 0 {
		dbConf.InsertPaste(key, content, secret, -1)
	} else {
		dbConf.InsertPaste(key, content, secret, time.Duration(expiration*int(time.Second)))
	}

	http.Redirect(w, r, key, http.StatusSeeOther)
}

// View for plak
func (config WebConfig) View(w http.ResponseWriter, r *http.Request) {
	dbConf := database.DBConfig{
		DB: config.DB,
	}
	var plak Plak
	key := r.PathValue("key")

	if dbConf.UrlExist(key) {
		plak = Plak{
			Key: key,
			DB:  config.DB,
		}
		plak = plak.GetContent()
		if r.PathValue("settings") == "raw" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := io.WriteString(w, plak.Content)
			if err != nil {
				log.Println(err)
			}
		} else {
			t, err := template.ParseFS(config.Templates, "templates/paste.html")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
			}
			err = t.Execute(w, plak)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Delete manage plak deletion endpoint
func (config WebConfig) Delete(w http.ResponseWriter, r *http.Request) {
	dbConf := database.DBConfig{
		DB: config.DB,
	}
	key := r.PathValue("key")

	if dbConf.UrlExist(key) {
		secret := r.URL.Query().Get("secret")
		if dbConf.VerifySecret(key, secret) {
			plak := Plak{
				Key: key,
				DB:  config.DB,
			}
			err := plak.deletePlak()
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// deletePlak Delete plak from database
func (plak Plak) deletePlak() error {
	err := plak.DB.Del(ctx, plak.Key).Err()
	if err != nil {
		log.Println(err)
		return &DeletePlakError{Name: plak.Key, Err: err}
	}

	return nil
}

// GetContent get plak content
func (plak Plak) GetContent() Plak {
	plak.Content = plak.DB.HGet(ctx, plak.Key, "content").Val()
	return plak
}
