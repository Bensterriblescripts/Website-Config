package main

import (
	"config/fn"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	Run()
}

func Run() {
	fn.LoadFilesIntoCache()
	http.Handle("/", http.HandlerFunc(handlerHome))
	http.Handle("/config", http.HandlerFunc(handlerConfig))
	http.Handle("/static/", http.HandlerFunc(fn.CachedFileHandler))

	fmt.Println("Initiating Development Webserver on 3050...")
	localServer := &http.Server{
		Addr:         "127.0.0.1:3050",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	fn.CheckPanicErr(localServer.ListenAndServe())
}

/* Pages */
func handlerHome(w http.ResponseWriter, r *http.Request) {
	fn.ServePage(w, "index.html")
}

/* POST Pages */
func handlerConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	body := fn.PrintErr(io.ReadAll(r.Body))
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Display the parsed JSON
	fmt.Printf("Received JSON data: %+v\n", data)

	// You can also access specific fields if you know the structure
	if testValue, exists := data["test"]; exists {
		fmt.Printf("Test field value: %v\n", testValue)
	}

	// Send a response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "message": "JSON received"}
	json.NewEncoder(w).Encode(response)
}
