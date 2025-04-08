// go-server/main.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go backend 👋")
	})

	fmt.Println("Go server running on :3000")
	http.ListenAndServe(":6969", nil)
}
