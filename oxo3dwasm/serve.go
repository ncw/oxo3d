//+build none

package main

import (
	"fmt"
	"log"
	"net/http"
)

func wasmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	http.ServeFile(w, r, "oxo3d.wasm")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/oxo3d.wasm", wasmHandler)
	fmt.Printf("Serving on http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
