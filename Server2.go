package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/julienschmidt/httprouter"

	"strconv"
)

type KeyValueId struct {
	Key   int
	Value string
}

type KeyValueIdArr struct {
	Maps []KeyValueId
}

var data = make(map[int]string)

func updateKeyValueId(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	key, _ := strconv.Atoi(p.ByName("key"))
	value := p.ByName("value")
	data[key] = value
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprint(rw, "200")
}

func getKeyValueId(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	key, _ := strconv.Atoi(p.ByName("key"))
	keyValueId := new(KeyValueId)
	keyValueId.Key = key
	keyValueId.Value = data[key]
	outgoingJSON, err := json.Marshal(keyValueId)
	if err != nil {
		//log.Println(error.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprint(rw, string(outgoingJSON))
}

func getAllKeyValueId(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	keyValueIdArr := new(KeyValueIdArr)
	keyValueIdArr.Maps = []KeyValueId{}
	for k, v := range data {
		keyValueId := new(KeyValueId)
		keyValueId.Key = k
		keyValueId.Value = v
		keyValueIdArr.Maps = append(keyValueIdArr.Maps, *keyValueId)
	}
	outgoingJSON, err := json.Marshal(keyValueIdArr)
	if err != nil {
		//log.Println(error.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprint(rw, string(outgoingJSON))
}

func main() {
	mux := httprouter.New()
	mux.PUT("/keys/:key/:value", updateKeyValueId)
	mux.GET("/keys/:key", getKeyValueId)
	mux.GET("/keys", getAllKeyValueId)
	server := http.Server{
		Addr:    "0.0.0.0:3001",
		Handler: mux,
	}
	server.ListenAndServe()
}
