package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	errParamKey   = `fetching key param failed`
	invalidKey    = `requested key does not exist`
	errReadBody   = `failed to read request body`
	errEncode     = `encoding error response failed`
	errUnmarshall = `unmarshalling request body failed`
)

// todo check
type res struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func initServer() {
	r := mux.NewRouter()
	r.HandleFunc(`/storage/key`, retrieveVal).Methods(http.MethodGet)
	r.HandleFunc(`/storage/key`, storeVal).Methods(http.MethodPut)
	// todo neighbours

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(config.Port), r))
}

func retrieveVal(w http.ResponseWriter, r *http.Request) {
	var response res
	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = errParamKey
		log.Error(errParamKey, r.URL.String())

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	// todo check the corresponding node of this key

	// fetching value from store
	val, ok := dataStore.get(key)
	if !ok {
		log.Error(invalidKey, key)
		w.WriteHeader(http.StatusBadRequest)

		response.Error = invalidKey
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	// writing the response
	w.WriteHeader(http.StatusOK)
	response.Value = val
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeVal(w http.ResponseWriter, r *http.Request) {
	var response res
	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = errParamKey
		log.Error(errParamKey, r.URL.String())

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	// todo check the corresponding node of this key

	// reading the value from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Error = errReadBody
		log.Error(errReadBody, key)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	var val string
	err = json.Unmarshal(body, &val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = errUnmarshall
		log.Error(errUnmarshall, key)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	// storing the value
	dataStore.set(key, val)

	// writing the success header
	w.WriteHeader(http.StatusOK)
}
