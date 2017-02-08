package math

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"fmt"
)

const (
	errorHeader = "Error"
)

// Handler for http requests
type Handler struct {
	Calculator *Calc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	p := strings.Split(strings.Trim(r.URL.EscapedPath(), "/"), "/")
	switch len(p) {
	case 0:
		header.Add(errorHeader, "op {value,add,subtract} missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	case 2:
		switch operation(p[0]) {
		case add, subtract:
			header.Add(errorHeader, "missing value")
			w.WriteHeader(http.StatusBadRequest)
			return
		case value:
		default:
			header.Add(errorHeader, "invalid operation")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case 3:
		switch operation(p[0]) {
		case add, subtract:
		case value:
			header.Add(errorHeader, "value requested with invalid key")
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			header.Add(errorHeader, "invalid operation")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		header.Add(errorHeader, "Bad Request, expected /{add,subtract}/{key}/{value} or /value/{key}")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	op := operation(p[0])
	key := p[1]

	var val, sum int
	var err error
	switch op {
	case add, subtract:
		val, err = strconv.Atoi(p[2])
		if err != nil {
			break
		}
		switch op {
		case add:
			sum, err = h.Calculator.Add(val, key)
		case subtract:
			sum, err = h.Calculator.Subtract(val, key)
		}
	case value:
		sum, err = h.Calculator.Value(key)
	}

	if err != nil {
		header.Add(errorHeader, "processing error: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("%s %s %d %d", op, key, val, sum))
}
