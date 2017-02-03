package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
