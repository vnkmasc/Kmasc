package api

import (
	"encoding/json"
	"log"
	"net/http"
	"openssl-ca/internal/ca"
	"openssl-ca/internal/config"
	"openssl-ca/internal/storage"
	"os"
	"path/filepath"
	"strings"
)

var store = storage.NewMemoryStore()

type RegisterRequest struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	OU     string `json:"ou"`
}

type EnrollRequest struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Hosts  string `json:"hosts"`
	OU     string `json:"ou"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	store.Register(req.ID, req.Secret)
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func EnrollHandler(w http.ResponseWriter, r *http.Request) {
	var req EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	if !store.Verify(req.ID, req.Secret) {
		http.Error(w, "unauthorized", 401)
		return
	}
	hosts := []string{}
	if req.Hosts != "" {
		hosts = strings.Split(req.Hosts, ",")
	} else {
		hosts = []string{"localhost"}
	}

	os.MkdirAll(filepath.Join(config.DataPath, "msp"), 0755)
	keyFile := filepath.Join(config.DataPath, "msp", req.ID+"-key.pem")
	csrFile := filepath.Join(config.DataPath, "msp", req.ID+".csr")
	certFile := filepath.Join(config.DataPath, "msp", req.ID+"-cert.pem")

	ca.GenerateUserKeyCSR(keyFile, csrFile, req.ID, req.OU)
	ca.SignCSRWithSAN(csrFile, filepath.Join(config.DataPath, "ca-key.pem"), filepath.Join(config.DataPath, "ca-cert.pem"), certFile, hosts)

	caBytes, _ := os.ReadFile(filepath.Join(config.DataPath, "ca-cert.pem"))
	certBytes, _ := os.ReadFile(certFile)
	keyBytes, _ := os.ReadFile(keyFile)

	json.NewEncoder(w).Encode(map[string]string{
		"caCert":   string(caBytes),
		"signCert": string(certBytes),
		"key":      string(keyBytes),
	})
}

func InitCAHandler(w http.ResponseWriter, r *http.Request) {
	certPath := filepath.Join(config.DataPath, "ca-cert.pem")
	keyPath := filepath.Join(config.DataPath, "ca-key.pem")

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		if err := ca.GenerateRootCA(keyPath, certPath); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	caCert, err := os.ReadFile(certPath)
	if err != nil {
		http.Error(w, "cannot read ca cert", 500)
		return
	}
	caKey, err := os.ReadFile(keyPath)
	if err != nil {
		http.Error(w, "cannot read ca key", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"caCert": string(caCert),
		"caKey":  string(caKey),
	})
}

func TLSEnrollHandler(w http.ResponseWriter, r *http.Request) {
	var req EnrollRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if !store.Verify(req.ID, req.Secret) {
		http.Error(w, "unauthorized", 401)
		return
	}

	tlsKey := filepath.Join(config.DataPath, "msp", req.ID+"-tls.key")
	tlsCSR := filepath.Join(config.DataPath, "msp", req.ID+"-tls.csr")
	tlsCert := filepath.Join(config.DataPath, "msp", req.ID+"-tls.crt")

	hosts := []string{}
	if req.Hosts != "" {
		hosts = strings.Split(req.Hosts, ",")
	} else {
		hosts = []string{"localhost"}
	}
	log.Println("======", req)
	ca.GenerateUserKeyCSRWithSAN(tlsKey, tlsCSR, req.ID, req.OU, hosts)
	ca.SignCSRWithSAN(tlsCSR, filepath.Join(config.DataPath, "ca-key.pem"), filepath.Join(config.DataPath, "ca-cert.pem"), tlsCert, hosts)

	caBytes, _ := os.ReadFile(filepath.Join(config.DataPath, "ca-cert.pem"))
	tlsCertBytes, _ := os.ReadFile(tlsCert)
	tlsKeyBytes, _ := os.ReadFile(tlsKey)

	json.NewEncoder(w).Encode(map[string]string{
		"caCert":  string(caBytes),
		"tlsCert": string(tlsCertBytes),
		"tlsKey":  string(tlsKeyBytes),
	})
}
