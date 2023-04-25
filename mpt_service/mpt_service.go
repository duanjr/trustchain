package mpt_service

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/trie"
	"io/ioutil"
	"net/http"
)

// 初始化一个空的Merkle Patricia Trie
var memDB = memorydb.New()
var t, _ = trie.New(common.Hash{}, trie.NewDatabase(memDB))

type KeyValue struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var kv KeyValue
	err = json.Unmarshal(body, &kv)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	_ = t.TryUpdate([]byte(kv.Key), []byte(fmt.Sprintf("%f", kv.Value)))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key-value pair inserted successfully"))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var kv KeyValue
	err = json.Unmarshal(body, &kv)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	value, err := t.TryGet([]byte(kv.Key))
	if err != nil {
		http.Error(w, "Error querying trie", http.StatusInternalServerError)
		return
	}

	if value == nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"key": "%s", "value": %s}`, kv.Key, value)))
}

func main() {
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/query", queryHandler)

	fmt.Println("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
