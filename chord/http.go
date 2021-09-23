package chord

import (
	"context"
	"dht/logger"
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

type response struct {
	Error string `json:"error"`
}

func InitServer(ctx context.Context) {
	r := mux.NewRouter()
	//r.HandleFunc(`/storage/{key}`, func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Println("RECEIVED")
	//}).Methods(http.MethodPost)
	r.HandleFunc(`/storage/{key}`, retrieveVal).Methods(http.MethodGet)
	r.HandleFunc(`/storage/{key}`, storeVal).Methods(http.MethodPost)
	r.HandleFunc(`/neighbors`, getNeighbours).Methods(http.MethodGet)
	logger.Log.InfoContext(ctx, `initializing http server`)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Config.Port), r))
}

func getNeighbours(w http.ResponseWriter, r *http.Request) {
	res := []string{Config.Predecessor + `:` + Config.PredecessorPort, Config.Successor + `:` + Config.SuccessorPort}

	// writing the response
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err, errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func retrieveVal(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug(`request received for get key`)
	var res response
	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = errParamKey
		logger.Log.Error(errParamKey, r.URL.String())

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Log.Error(err, errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = errBucket
		logger.Log.Error(err, errBucket, key)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Log.Error(err, errEncode, r.URL.String())
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
			logger.Log.Error(errStore, key)

			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				logger.Log.Error(err, errEncode, r.URL.String())
			}
			return
		}
	} else {
		val, err = client.proceedGetKey(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res.Error = err.Error()
			logger.Log.Error(err, key)

			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				logger.Log.Error(err, errEncode, r.URL.String())
			}
			return
		}
	}

	// writing the response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		logger.Log.Error(err, errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeVal(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug(`request received for store key`)
	var res response
	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = errParamKey
		logger.Log.Error(errParamKey, r.URL.String())

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Log.Error(err, errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = errBucket
		logger.Log.Error(err, errBucket, key)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Log.Error(err, errEncode, r.URL.String())
		}
		return
	}

	if !includes {
		status, err := client.proceedStoreKey(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res.Error = err.Error()
			logger.Log.Error(err, key)

			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				logger.Log.Error(err, errEncode, r.URL.String())
			}
			return
		}

		w.WriteHeader(status)
		return
	}

	// reading the value from request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res.Error = errReadBody
		logger.Log.Error(err, errReadBody, key)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			logger.Log.Error(err, errEncode, r.URL.String())
		}
		return
	}

	dataStore.set(key, string(data))

	// writing the success header
	w.WriteHeader(http.StatusOK)
}
