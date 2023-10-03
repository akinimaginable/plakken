package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var currentConfig config

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.ReplaceAll(r.URL.Path, "/raw", "")
	switch r.Method {
	case "GET":
		if path == "/" {
			http.ServeFile(w, r, "./static/index.html")

		} else if strings.HasPrefix(path, "/static/") {
			fs := http.FileServer(http.Dir("./static"))
			http.Handle("/static/", http.StripPrefix("/static/", fs))
		} else {
			if urlExist(path) {
				pasteContent := getContent(path)
				fmt.Println(pasteContent)
				if strings.HasSuffix("/raw", path) {
					io.WriteString(w, pasteContent)
				} else {
					io.WriteString(w, pasteContent)
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
	currentConfig = setConfig()
	listen := currentConfig.host + ":" + currentConfig.port
	http.HandleFunc("/", handleRequest)

	err := http.ListenAndServe(listen, nil)
	if err != nil {
		fmt.Println(err)
	}
}
