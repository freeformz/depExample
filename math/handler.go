package math

import (
	"io"
	"net/http"
	"strconv"

	"fmt"

	"github.com/gorilla/mux"
)

const (
	errorHeader = "Error"
)

// Handler for http requests
type Handler struct {
	Calculator *Calc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	header := w.Header()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	o, ok := vars["op"]
	if !ok {
		header.Add(errorHeader, "op missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	op := operation(o)

	key, ok := vars["key"]
	if !ok {
		header.Add(errorHeader, "key missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var val, sum int
	var err error
	switch op {
	case add, subtract:
		vs, ok := vars["value"]
		if !ok {
			header.Add(errorHeader, "value missing")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		val, err = strconv.Atoi(vs)
		if err != nil {
			header.Add(errorHeader, "error converting value: "+err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch op {
		case add:
			sum, err = h.Calculator.Add(val, key)
		case subtract:
			sum, err = h.Calculator.Subtract(val, key)
		}
	case value:
		sum, err = h.Calculator.Value(key)
	default:
		header.Add(errorHeader, "unknown op: "+string(op))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		header.Add(errorHeader, "processing error: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("%s %s %d %d", op, key, val, sum))
}
