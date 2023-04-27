package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	iterations = 100
)

//type KeyValue struct {
//	Key   string  `json:"key"`
//	Value float64 `json:"value"`
//}

func main() {
	//startTime := time.Now()
	//for i := 0; i < iterations; i++ {
	//	kv := KeyValue{
	//		Key:   fmt.Sprintf("%d&%d", i, i),
	//		Value: float64(i%(2*iterations))/iterations - 1,
	//	}
	//	jsonData, _ := json.Marshal(kv)
	//	resp, err := http.Post("http://localhost:8080/insert", "application/json", bytes.NewBuffer(jsonData))
	//	if err != nil {
	//		fmt.Printf("Error inserting key-value pair: %v\n", err)
	//		continue
	//	}
	//	resp.Body.Close()
	//	if resp.StatusCode != http.StatusOK {
	//		fmt.Printf("Error inserting key-value pair: insert request failed with status: %s\n", resp.Status)
	//	}
	//}
	//
	//elapsedTime := time.Since(startTime)
	//fmt.Printf("Total time for %d insertions: %v\n", iterations, elapsedTime)
	startTime := time.Now()
	for j := 0; j < 10; j++ {
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("192.168.0.%d", i+1)
			keyJson := fmt.Sprintf(`{"addresss": "%s"}`, key)
			resp, err := http.Post("http://localhost:8080/pki/query", "application/json", bytes.NewBufferString(keyJson))
			if err != nil {
				fmt.Printf("Error querying key-value pair: %v\n", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error querying key-value pair: query request failed with status: %s\n", resp.Status)
			} else {
				defer resp.Body.Close()
				//body, _ := ioutil.ReadAll(resp.Body)
				//var kv KeyValue
				//json.Unmarshal(body, &kv)
			}
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time for %d queries: %v\n", iterations, elapsedTime)
}
