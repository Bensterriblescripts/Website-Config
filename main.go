package main

import (
	"config/fn"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"
)

func main() {
	ldapURL := "ldaps://192.168.5.104:1636"
	username := "uid=admin"
	password := "Worldwarcraft3#%#%"
	fn.CheckPanicErr(bindUser(ldapURL, username, password))

	Run()
}

func Run() {
	fn.LoadFilesIntoCache()
	http.Handle("/", http.HandlerFunc(handlerHome))
	http.Handle("/files", http.HandlerFunc(handlerConfig))
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

	// Parse
	body := fn.PrintErr(io.ReadAll(r.Body))
	defer r.Body.Close()
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Return Content
	if configFile, exists := data["config"]; exists {
		fmt.Printf("Requested Configuration: %v\n", configFile)

		if content, fileExists := fn.GetCachedFile("configs/" + configFile.(string)); fileExists {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", configFile))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(content)
			return
		} else {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	} else if configFile, exists := data["executable"]; exists {
		fmt.Printf("Requested Executable: %v\n", configFile)

		if content, fileExists := fn.GetCachedFile("executables/" + configFile.(string)); fileExists {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", configFile))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(content)
			return
		} else {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	} else if configFile, exists := data["adminscript"]; exists {
		fmt.Printf("Requested Script: %v\n", configFile)

		if content, fileExists := fn.GetCachedFile("scripts/admin/" + configFile.(string)); fileExists {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", configFile))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(content)
			return
		}
	} else if configFile, exists := data["userscript"]; exists {
		fmt.Printf("Requested Script: %v\n", configFile)

		if content, fileExists := fn.GetCachedFile("scripts/user/" + configFile.(string)); fileExists {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", configFile))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(content)
			return
		}
	} else {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
}

/* Auth */
func bindUser(ldapURL, username, password string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn := fn.PanicErr(ldap.DialTLS("tcp", "192.168.5.104:1636", tlsConfig))
	defer conn.Close()

	fn.CheckPanicErr(conn.Bind(username, password))
	fmt.Println("LDAP Server Connected Successfully")
	return nil
}
