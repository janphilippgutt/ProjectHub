package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Zag! Welcome to your project.")
}

func main() {

	http.HandleFunc("/", helloHandler)

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
