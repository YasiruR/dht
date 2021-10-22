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
	"strings"
)

// error definitions
const (
	errParamKey   = `invalid path parameter (key)`
	errParamHost  = `invalid path parameter (host)`
	errParamPrime = `invalid query parameter (nprime)`
	errEncode     = `encoding error response failed`
	errBucket     = `checking bucket failed`
	errJoin       = `initiating join failed`
	errStore      = `requested key does not exist in store`
	errReadBody   = `failed to read request body`
	errUnmarshall = `unmarshalling request body failed`
)

func InitServer(ctx context.Context) {
	r := mux.NewRouter()
	// http endpoints for store operations
	r.HandleFunc(`/storage/{key}`, retrieveVal).Methods(http.MethodGet)
	r.HandleFunc(`/storage/{key}`, storeVal).Methods(http.MethodPut)

	// http endpoints for join and leave
	r.HandleFunc(`/join`, join).Methods(http.MethodPost)
	r.HandleFunc(`/internal/join/{host}`, internalJoin).Methods(http.MethodPost)
	r.HandleFunc(`/internal/update-successor/{host}`, updateSuccessor).Methods(http.MethodPost)

	// http endpoints for debugging the network
	r.HandleFunc(`/neighbors`, getNeighbours).Methods(http.MethodGet)

	logger.Log.InfoContext(ctx, `initializing http server`)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Config.Port), r))
}

func getNeighbours(w http.ResponseWriter, r *http.Request) {
	res := []string{neighbors.predHostname, neighbors.sucHostname}

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
		_, err := w.Write([]byte(`invalid storeUrl with param key`))
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
		val, statusCode, err = neighbors.proceedGetKey(key, r)
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
		_, err := w.Write([]byte(`invalid storeUrl with param key`))
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
		status, err := neighbors.proceedStoreKey(key, r)
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

/* dynamic join and leave functionality */

func join(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	logger.Log.DebugContext(ctx, `request received for join`)

	networkHost := r.URL.Query().Get(`nprime`)
	if networkHost == `` {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error(errParamPrime, r.URL.String())
		_, err := w.Write([]byte(`query param nprime is missing`))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	pred, suc, err := neighbors.initJoin(networkHost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.ErrorContext(ctx, err, errJoin, networkHost)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	fmt.Println("PREDD : ", pred, len(pred))
	fmt.Println("SUCC : ", suc, len(suc))

	neighbors.updatePredecessor(pred)
	neighbors.updateSuccessor(suc)

	w.WriteHeader(http.StatusOK)
	return
}

func internalJoin(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	logger.Log.DebugContext(ctx, `request received for internal join`)

	// fetching the hostname
	hostname, ok := mux.Vars(r)[`host`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error(errParamKey, r.URL.String())
		_, err := w.Write([]byte(`invalid internal join url with param hostname`))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	includes, err := node.checkKey(hostname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.ErrorContext(ctx, err, errBucket, hostname)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	if !includes {
		resData, err := neighbors.proceedJoin(hostname, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.ErrorContext(ctx, err, errBucket, hostname)
			_, err = w.Write([]byte(errBucket))
			if err != nil {
				logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resData)
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var newNeighbors string
	// join to the first node
	if node.single {
		newNeighbors = strings.Join([]string{hostname, hostname}, ",")
		neighbors.updateSuccessor(hostname)
		neighbors.updatePredecessor(hostname)
	} else {
		// informs ex-predecessor to update its successor
		err = neighbors.notifyPredecessor(hostname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.ErrorContext(ctx, err, errBucket, hostname)
			_, err = w.Write([]byte(errBucket))
			if err != nil {
				logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
			}
			return
		}
		// setting ex-predecessor as predecessor of the new node
		exPred := neighbors.predHostname
		// adds new node as predecessor
		neighbors.updatePredecessor(hostname)

		// returns successor and predecessor of new node
		newNeighbors = strings.Join([]string{exPred, node.hostname}, ",")
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newNeighbors)
	if err != nil {
		logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}

	return
}

func updateSuccessor(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	logger.Log.DebugContext(ctx, `request received for internal update successor`)

	// fetching hostname
	hostname, ok := mux.Vars(r)[`host`]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error(errParamHost, r.URL.String())
		_, err := w.Write([]byte(`invalid internal update successor url with param hostname`))
		if err != nil {
			logger.Log.ErrorContext(ctx, err, errEncode, r.URL.String())
		}
		return
	}

	neighbors.updateSuccessor(hostname)
	w.WriteHeader(http.StatusOK)
}
