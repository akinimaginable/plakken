package main

import (
	"fmt"
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
			w.WriteHeader(http.StatusNotFound)
		}
	case "POST":
		fmt.Println("Post!")
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
