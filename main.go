package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var currentConfig config

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch r.Method {
	case "GET":
		if path == "/" {
			http.ServeFile(w, r, "./static/index.html")
		} else if strings.HasPrefix(path, "/static/") {
			fs := http.FileServer(http.Dir("./static"))
			http.Handle("/static/", http.StripPrefix("/static/", fs))
		} else {
			if urlExist(path) {
				_, err := io.WriteString(w, "This plak exists")
				if err != nil {
					return
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
	}
}

func main() {
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
