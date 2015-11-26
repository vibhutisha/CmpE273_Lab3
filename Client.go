package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/julienschmidt/httprouter"

	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type KeyValueId struct {
	Key   int
	Value string
}

var HashingMapping = make(map[string]string)
var Values = make(map[int]string)
var ListOfServers = []string{"3000", "3001", "3002"}
var SortedHashingMappingKeys []string

func calculateHashValue(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}

func getServerValue(n int) string {
	readValue := 0
	readCache := 0
	index := 0
	for index < len(SortedHashingMappingKeys) {
		if readValue != 1 {
			if calculateHashValue(strconv.Itoa(n)) == SortedHashingMappingKeys[index] {
				readValue = 1
			}
		} else if readValue == 1 {
			if createStrings(HashingMapping[SortedHashingMappingKeys[index]], ListOfServers) {
				readCache = 1
				break
			}
		}
		if index == len(SortedHashingMappingKeys)-1 && readCache == 0 {
			index = 0
		} else {
			index += 1
		}
	}
	return HashingMapping[SortedHashingMappingKeys[index]]
}

func putDataTo(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	k, _ := strconv.Atoi(p.ByName("key"))
	value := p.ByName("value")
	url := "http://localhost:" + getServerValue(k) + "/keys/" + strconv.Itoa(k) + "/" + value
	req, err := http.NewRequest("PUT", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var values interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &values)
	var m = values.(interface{}).(float64)
	fmt.Fprint(rw, m)
}

func getRequest(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	k, _ := strconv.Atoi(p.ByName("key"))
	resp, err := http.Get("http://localhost:" + getServerValue(k) + "/keys/" + strconv.Itoa(k))
	if err == nil {
		var values interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &values)
		var m = values.(map[string]interface{})
		keyValueId := new(KeyValueId)
		keyValueId.Key = int(m["Key"].(float64))
		keyValueId.Value = m["Value"].(string)
		outgoingJSON, err := json.Marshal(keyValueId)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprint(rw, string(outgoingJSON))
	} else {
		fmt.Println(err)
	}
}

func createStrings(str1 string, list []string) bool {
	for _, v := range list {
		if v == str1 {
			return true
		}
	}
	return false
}

func main() {
	Values[1] = "z"
	Values[2] = "y"
	Values[3] = "x"
	Values[4] = "q"
	Values[5] = "w"
	Values[6] = "v"
	Values[7] = "u"
	Values[8] = "t"
	Values[9] = "s"
	Values[10] = "r"

	for _, each := range ListOfServers {
		HashingMapping[calculateHashValue(each)] = each
	}

	for k, _ := range Values {
		HashingMapping[calculateHashValue(strconv.Itoa(k))] = strconv.Itoa(k)
	}

	for k, _ := range HashingMapping {
		SortedHashingMappingKeys = append(SortedHashingMappingKeys, k)
	}

	sort.Strings(SortedHashingMappingKeys)
	mux := httprouter.New()
	mux.PUT("/keys/:key/:value", putDataTo)
	mux.GET("/keys/:key", getRequest)
	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}
	server.ListenAndServe()
}
