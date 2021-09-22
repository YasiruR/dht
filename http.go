package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

// error definitions
const (
	errParamKey = `invalid path parameter (key)`
	errEncode   = `encoding error response failed`
	errBucket   = `checking bucket failed`

	errStore      = `requested key does not exist in store`
	errReadBody   = `failed to read request body`
	errUnmarshall = `unmarshalling request body failed`
)

// todo check
type response struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func initServer() {
	r := mux.NewRouter()
	r.HandleFunc(`/storage/{key}`, retrieveVal).Methods(http.MethodGet)
	r.HandleFunc(`/storage/{key}`, storeVal).Methods(http.MethodPut)
	// todo neighbours

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), r))
}

func retrieveVal(w http.ResponseWriter, r *http.Request) {
	var res response
	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = errParamKey
		log.Error(errParamKey, r.URL.String())

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = errBucket
		log.Error(err.Error(), errBucket, key)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	var val string
	if includes {
		// fetching value from store
		val, ok = dataStore.get(key)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			res.Error = errStore
			log.Error(errStore, key)

			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				log.Error(errEncode, r.URL.String())
			}
			return
		}
	} else {
		val, err = client.proceedGetKey(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res.Error = err.Error()
			log.Error(err, key)

			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				log.Error(errEncode, r.URL.String())
			}
			return
		}
	}

	// writing the response
	w.WriteHeader(http.StatusOK)
	res.Value = val
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeVal(w http.ResponseWriter, r *http.Request) {
	var res response
	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = errParamKey
		log.Error(errParamKey, r.URL.String())

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = errBucket
		log.Error(err.Error(), errBucket, key)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	if !includes {
		status, err := client.proceedStoreKey(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res.Error = err.Error()
			log.Error(err, key)

			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				log.Error(errEncode, r.URL.String())
			}
			return
		}

		w.WriteHeader(status)
		return
	}

	// reading the value from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = errReadBody
		log.Error(errReadBody, key)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(errEncode, r.URL.String())
		}
		return
	}

	var val string
	err = json.Unmarshal(body, &val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = errUnmarshall
		log.Error(errUnmarshall, key)

		err = json.NewEncoder(w).Encode(res)
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
