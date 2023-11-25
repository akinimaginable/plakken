package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var currentConfig config
var db *redis.Client

type pasteView struct {
	Content string
	Key     string
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	clearPath := strings.ReplaceAll(r.URL.Path, "/raw", "")
	switch r.Method {
	case "GET":
		if path == "/" {
			http.ServeFile(w, r, "./static/index.html")

		} else if strings.HasPrefix(path, "/static/") {
			fs := http.FileServer(http.Dir("./static"))
			http.Handle("/static/", http.StripPrefix("/static/", fs))
		} else {
			if urlExist(clearPath) {
				if strings.HasSuffix(path, "/raw") {
					pasteContent := getContent(clearPath)
					w.Header().Set("Content-Type", "text/plain")
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
			secret := generateSecret()
			url := "/" + generateUrl()
			content := r.FormValue("content")
			insertPaste(url, content, secret, -1)
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "DELETE":
		if strings.HasPrefix(path, "/delete") {
			urlItem := strings.Split(path, "/")
			if urlExist("/" + urlItem[2]) {
				secret := r.URL.Query().Get("secret")
				if verifySecret("/"+urlItem[2], secret) {
					deleteContent("/" + urlItem[2])
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	db = connectDB()
	currentConfig = getConfig()
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
