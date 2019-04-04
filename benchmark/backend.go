package main

import (
	"log"
	"net/http"
)

var ok = []byte("ok")

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(ok)
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("listen 2000")
	http.ListenAndServe("127.0.0.1:2000", nil)
}
