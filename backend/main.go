package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var payloads map[int][]byte
var once sync.Once

func main() {
	generatePayloads()

	http.HandleFunc("/publiccache/", func(w http.ResponseWriter, r *http.Request) {
		PublicCacheHandler(w, r)
	})
	
   
    http.HandleFunc("/nocache", func(w http.ResponseWriter, r *http.Request) {
        NoCacheHandler(w, r)
    })

    http.HandleFunc("/privatecache", func(w http.ResponseWriter, r *http.Request) {
        PrivateCacheHandler(w, r)
    })
    http.HandleFunc("/getresponse", func(w http.ResponseWriter, r *http.Request) {
        getresponseWithoutHeaders(w, r)
    })

	// Start server
	fmt.Println("Server is listening on port 8083...")
	http.ListenAndServe(":8083", nil)
}

// func PublicCacheHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.URL.Path[len("/publiccache/"):])
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	payload, ok := payloads[id]
// 	if !ok {
// 		http.Error(w, "ID not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Cache-Control", "public, max-age=180")
// 	w.Write(payload)
// 	additionalContent := []byte(" cached as public for ID " + strconv.Itoa(id))
// 	w.Write(additionalContent)
// }

func PublicCacheHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-current-user")
	if username == "" {
		http.Error(w, "User information missing", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/publiccache/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	payload, ok := payloads[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

    w.Header().Set("X-Custom-Header", "CustomHeaderValue")
	w.Header().Set("Cache-Control", "public, max-age=600")
	w.Write(payload)
	additionalContent := []byte(" cached as public for ID " + strconv.Itoa(id) + " for user " + username)
	w.Write(additionalContent)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("AUTH_TOKEN")
		if token == "" {
			http.Error(w, "Authorization token missing", http.StatusInternalServerError)
			return
		}
		r.Header.Set("Authorization", "Bearer "+token)
		next.ServeHTTP(w, r)
	}
}


func PrivateCacheHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "private, max-age=3600")
   // w.Write(payload)
    // sleepBeforeRespond()
}

func getresponseWithoutHeaders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Caching is applied through UI"))

}

func NoCacheHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-store")
   // w.Write(payload)
    // sleepBeforeRespond()
}

func generatePayloads() {
	once.Do(func() {
		payloads = make(map[int][]byte)
		for i := 1; i <= 50; i++ {
			payload := make([]byte, 102400)
			for j := range payload {
				payload[j] = 'x'
			}
			payloads[i] = payload
		}
	})
}
