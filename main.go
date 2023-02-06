package main

import (
	"compress/gzip"
	"crypto/rand"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
)

type RandServer struct {
	mu *sync.Mutex
}

func (r *RandServer) randHandler(w http.ResponseWriter, req *http.Request) {
	// Add a lock to prevent the server consuming too many resources
	r.mu.Lock()
	defer r.mu.Unlock()

	// This server just generate some useless garbage
	w.Header().Set("Content-Type", "application/octet-stream")

	// Generate 200 MB random bytes
	buf := make([]byte, 200*1024*1024)
	defer debug.FreeOSMemory()
	rand.Read(buf)

	// Just consume some CPU resource
	gw := gzip.NewWriter(w)
	defer gw.Close()
	gw.Write(buf)
	gw.Flush()
}

func main() {
	r := RandServer{mu: &sync.Mutex{}}
	http.HandleFunc("/", r.randHandler)
	log.Fatal(http.ListenAndServe(":43002", nil))
}
