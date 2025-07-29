package main

import (
	"log"
	"net/http"
	"openssl-ca/internal/api"
	"openssl-ca/internal/ca"
	"openssl-ca/internal/config"
	"os"
	"path"
)

func main() {

	if _, err := os.Stat(path.Join(config.DataPath, "ca-cert.pem")); os.IsNotExist(err) {
		log.Println("Generating root CA...")
		if err := ca.GenerateRootCA(path.Join(config.DataPath, "ca-key.pem"), path.Join(config.DataPath, "ca-cert.pem")); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/tlsenroll", api.TLSEnrollHandler)
	http.HandleFunc("/initca", api.InitCAHandler)
	http.HandleFunc("/register", api.RegisterHandler)
	http.HandleFunc("/enroll", api.EnrollHandler)
	log.Println("Custom CA running on ", config.Binding)
	log.Fatal(http.ListenAndServe(config.Binding, nil))
}

//curl -X POST http://localhost:7054/register -d '{"id":"user1","secret":"abc123"}'
//curl -X POST http://localhost:7054/enroll -d '{"id":"user1","secret":"abc123"}'
