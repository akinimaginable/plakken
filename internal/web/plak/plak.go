package plak

import (
	"context"
	"embed"
	"io"
	"log"
	"net/http"
	"time"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
	"git.gnous.eu/gnouseu/plakken/internal/database"
	"git.gnous.eu/gnouseu/plakken/internal/secret"
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

// plak "Object" for plak
type plak struct {
	Key        string
	Content    string
	Expiration time.Duration
	DB         *redis.Client
}

// create a plak
func (plak plak) create() (string, error) {
	dbConf := database.DBConfig{
		DB: plak.DB,
	}

	token, err := secret.GenerateToken()
	if err != nil {
		return "", err
	}

	if dbConf.UrlExist(plak.Key) {
		return "", &createError{message: "key already exist"}
	}

	var hashedSecret string
	hashedSecret, err = secret.Password(token)
	if err != nil {
		return "", err
	}

	dbConf.InsertPaste(plak.Key, plak.Content, hashedSecret, plak.Expiration)

	return token, nil
}

// PostCreate manage POST request for create plak
func (config WebConfig) PostCreate(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filename := r.FormValue("filename")
	var key string
	if len(filename) == 0 {
		key = utils.GenerateUrl(config.UrlLength)
	} else {
		if !utils.ValidKey(filename) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key = filename
	}

	plak := plak{
		Key:     key,
		Content: content,
		DB:      config.DB,
	}

	rawExpiration := r.FormValue("exp")
	expiration, err := utils.ParseExpiration(rawExpiration)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if expiration == 0 {
		plak.Expiration = -1
	} else {
		plak.Expiration = time.Duration(expiration * int(time.Second))
	}

	_, err = plak.create()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+key, http.StatusSeeOther)
}

// CurlCreate PostCreate plak with minimum param, ideal for curl. Force 7 day expiration
func (config WebConfig) CurlCreate(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	content, _ := io.ReadAll(r.Body)
	err := r.Body.Close()
	if err != nil {
		log.Println(err)
	}

	key := utils.GenerateUrl(config.UrlLength)

	plak := plak{
		Key:        key,
		Content:    string(content),
		Expiration: constant.ExpirationCurlCreate,
		DB:         config.DB,
	}

	var token string
	token, err = plak.create()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var baseURL string
	if r.TLS == nil {
		baseURL = "http://" + r.Host + "/" + key
	} else {
		baseURL = "https://" + r.Host + "/" + key
	}

	message := baseURL + "\n" + "Delete with : 'curl -X DELETE " + baseURL + "?secret\\=" + token + "'" + "\n"

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = io.WriteString(w, message)
	if err != nil {
		log.Println(err)
	}
}

// View for plak
func (config WebConfig) View(w http.ResponseWriter, r *http.Request) {
	dbConf := database.DBConfig{
		DB: config.DB,
	}
	var currentPlak plak
	key := r.PathValue("key")

	if dbConf.UrlExist(key) {
		currentPlak = plak{
			Key: key,
			DB:  config.DB,
		}
		currentPlak = currentPlak.getContent()
		if r.PathValue("settings") == "raw" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := io.WriteString(w, currentPlak.Content)
			if err != nil {
				log.Println(err)
			}
		} else {
			t, err := template.ParseFS(config.Templates, "templates/paste.html")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
			err = t.Execute(w, currentPlak)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// DeleteRequest manage plak deletion endpoint
func (config WebConfig) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	dbConf := database.DBConfig{
		DB: config.DB,
	}
	key := r.PathValue("key")
	var valid bool

	if dbConf.UrlExist(key) {
		var err error
		token := r.URL.Query().Get("secret")

		valid, err = dbConf.VerifySecret(key, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if valid {
			plak := plak{
				Key: key,
				DB:  config.DB,
			}
			err := plak.delete()
			if err != nil {
				log.Println(err)
			}

			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

// delete DeleteRequest plak from database
func (plak plak) delete() error {
	err := plak.DB.Del(ctx, plak.Key).Err()
	if err != nil {
		log.Println(err)
		return &deletePlakError{name: plak.Key, err: err}
	}

	return nil
}

// getContent get plak content
func (plak plak) getContent() plak {
	plak.Content = plak.DB.HGet(ctx, plak.Key, "content").Val()
	return plak
}
