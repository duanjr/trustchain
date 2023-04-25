package mpt_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	iterations = 1000
)

//type KeyValue struct {
//	Key   string  `json:"key"`
//	Value float64 `json:"value"`
//}

func main() {
	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		kv := KeyValue{
			Key:   fmt.Sprintf("%d&%d", i, i),
			Value: float64(i%(2*iterations))/iterations - 1,
		}
		jsonData, _ := json.Marshal(kv)
		resp, err := http.Post("http://localhost:8080/insert", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error inserting key-value pair: %v\n", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error inserting key-value pair: insert request failed with status: %s\n", resp.Status)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time for %d insertions: %v\n", iterations, elapsedTime)
	startTime = time.Now()
	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("%d&%d", i, i)
		keyJson := fmt.Sprintf(`{"key": "%s"}`, key)
		resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(keyJson))
		if err != nil {
			fmt.Printf("Error querying key-value pair: %v\n", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error querying key-value pair: query request failed with status: %s\n", resp.Status)
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var kv KeyValue
			json.Unmarshal(body, &kv)
		}
	}
	elapsedTime = time.Since(startTime)
	fmt.Printf("Total time for %d queries: %v\n", iterations, elapsedTime)
}
