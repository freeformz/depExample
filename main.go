package main

import (
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/freeformz/depExample/math"
)

func main() {
	db, err := bolt.Open("maths.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	c := math.Calc{DB: db}
	h := math.Handler{Calculator: &c}

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/math/", http.StripPrefix("/math", &h))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
