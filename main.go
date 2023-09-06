package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Helloes")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Heelo :D"))
	})
	server := http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("ooopps, %v", err)
	}
}
