package main

import (
	"math/rand"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, req *http.Request) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if num := r1.Intn(100); num%5 == 0 {
		w.Write([]byte("index page"))
	} else {
		http.Error(w, "error page", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe(":8080", nil)
}
