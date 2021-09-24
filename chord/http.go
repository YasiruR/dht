package chord

import (
	"context"
	"dht/logger"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	traceableContext "github.com/tryfix/traceable-context"
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

//
//type response struct {
//	Error string `json:"error"`
//}

func InitServer(ctx context.Context) {
	r := mux.NewRouter()
	r.HandleFunc(`/storage/{key}`, retrieveVal).Methods(http.MethodGet)
	r.HandleFunc(`/storage/{key}`, storeVal).Methods(http.MethodPut)
	r.HandleFunc(`/neighbors`, getNeighbours).Methods(http.MethodGet)
	logger.Log.InfoContext(ctx, `initializing http server`)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Config.Port), r))
}

func getNeighbours(w http.ResponseWriter, r *http.Request) {
	var res []string
	if Config.FingerTableEnabled {
		res = Config.Nodes
	} else {
		res = []string{Config.Predecessor + `:` + Config.PredecessorPort, Config.Successor + `:` + Config.SuccessorPort}
	}

	// writing the response
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err, errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func retrieveVal(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	logger.Log.DebugContext(ctx, `request received for get key`)

	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.ErrorContext(ctx, errParamKey, r.URL.String())
		_, err := w.Write([]byte(`invalid url with param key`))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.ErrorContext(ctx, err, errBucket, key)
		_, err = w.Write([]byte(errBucket))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	var val string
	var statusCode int
	if includes {
		// fetching value from store
		val, ok = dataStore.get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			logger.Log.TraceContext(ctx, errStore, key)
			_, err = w.Write([]byte(fmt.Sprintf(`No object with key '%s' on this node`, key)))
			if err != nil {
				logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
			}
			return
		}
		statusCode = http.StatusOK
	} else {
		val, statusCode, err = client.proceedGetKey(key, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.ErrorContext(ctx, err, key)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
			}
			return
		}
	}

	// writing the response
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(val))
	if err != nil {
		logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeVal(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	logger.Log.DebugContext(ctx, `request received for store key`)

	// fetching the key
	key, ok := mux.Vars(r)[`key`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error(errParamKey, r.URL.String())
		_, err := w.Write([]byte(`invalid url with param key`))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.ErrorContext(ctx, err, errBucket, key)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	if !includes {
		status, err := client.proceedStoreKey(key, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.ErrorContext(ctx, err, errBucket, key)
			_, err = w.Write([]byte(errBucket))
			if err != nil {
				logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
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
		logger.Log.ErrorContext(ctx, err, errReadBody, key)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	dataStore.set(key, string(data))

	// writing the success header
	w.WriteHeader(http.StatusOK)
}
