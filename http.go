package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"net/http"
)

const (
	errParamKey = `fetching key param failed`
	invalidKey  = `requested key does not exist`
)

type res struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func initServer(port int) {
	r := mux.NewRouter()
}

func retrieve(w http.ResponseWriter, r *http.Request) {
	var response res
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(errParamKey, r.URL.String())

		response.Error = errParamKey
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(`encoding error response failed`, r.URL.String())
		}
		return
	}

	val, ok := s.get(key)
	if !ok {
		log.Error(invalidKey, key)
		w.WriteHeader(http.StatusBadRequest)

		response.Error = invalidKey
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(`encoding error response failed`, r.URL.String())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Value = val
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(`encoding error response failed`, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
