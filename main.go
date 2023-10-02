package main

import (
	"fmt"
	"net/http"
)

func main() {
	i := 0
	i = 5

	var k int8 = 0

	fmt.Println(i)
	fmt.Println(k)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, you're at %s", r.URL.Path)
		if err != nil {
			return
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
