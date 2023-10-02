package main

import (
	"fmt"
	"net/http"
)

func main() {
	currentConfig := setConfig()
	listen := currentConfig.host + ":" + currentConfig.port
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, you're at %s", r.URL.Path)
		if err != nil {
			return
		}
	})

	err := http.ListenAndServe(listen, nil)
	if err != nil {
		return
	}
}
