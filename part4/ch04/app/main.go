package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)

		if err != nil {
			http.Error(w, "Error json deserialize request body", http.StatusInternalServerError)
		}

		log.Printf("body : %v", data)
		fmt.Fprintf(w, "POST %v", data)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe(":8080", nil)
}
