package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var currentConfig Config
var db *redis.Client

type pasteView struct {
	Content string
	Key     string
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	clearPath := strings.ReplaceAll(r.URL.Path, "/raw", "")
	staticResource := "/static/"
	switch r.Method {
	case "GET":
		if path == "/" {
			http.ServeFile(w, r, "./static/index.html")

		} else if strings.HasPrefix(path, staticResource) {
			fs := http.FileServer(http.Dir("./static"))
			http.Handle(staticResource, http.StripPrefix(staticResource, fs))
		} else {
			if UrlExist(clearPath) {
				if strings.HasSuffix(path, "/raw") {
					pasteContent := getContent(clearPath)
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					_, err := io.WriteString(w, pasteContent)
					if err != nil {
						log.Println(err)
					}
				} else {
					pasteContent := getContent(path)
					s := pasteView{Content: pasteContent, Key: strings.TrimPrefix(path, "/")}
					t, err := template.ParseFiles("templates/paste.html")
					if err != nil {
						log.Println(err)
					}
					err = t.Execute(w, s)
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	case "POST":
		if path == "/" {
			secret := GenerateSecret()
			url := "/" + GenerateUrl()
			content := r.FormValue("content")
			rawExpiration := r.FormValue("exp")
			expiration, err := ParseExpiration(rawExpiration)
			if err != nil {
				log.Println(err)
			} else if expiration == 0 {
				insertPaste(url, content, secret, -1)
			} else if expiration == -1 {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				insertPaste(url, content, secret, time.Duration(expiration*int(time.Second)))
			}

			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "DELETE":
		if UrlExist(path) {
			secret := r.URL.Query().Get("secret")
			if VerifySecret(path, secret) {
				DeleteContent(path)
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func main() {
	db = ConnectDB()
	currentConfig = GetConfig()
	listen := currentConfig.host + ":" + currentConfig.port
	http.HandleFunc("/", handleRequest)

	if currentConfig.host == "" {
		fmt.Println("Listening on port " + listen)
	} else {
		fmt.Println("Listening on " + listen)
	}

	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
