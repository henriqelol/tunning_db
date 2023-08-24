package main

import (
	"fmt"
	"log"
	"net/http"
	"tuning_db/handle"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/tuning", handle.HandleRequest)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
