package main

import (
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/freeformz/depExample/math"
	"github.com/gorilla/mux"
)

func main() {
	db, err := bolt.Open("maths.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	c := math.Calc{DB: db}
	h := math.Handler{Calculator: &c}

	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(".")))
	r.Handle("/math/{op}/{key}", &h)
	r.Handle("/math/{op}/{key}/{value}", &h)
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
