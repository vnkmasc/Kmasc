package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/v2/metadata"
)

type CertificateOnChain struct {
	CertHash            string `json:"cert_hash"`
	HashFile            string `json:"hash_file"`
	UniversitySignature string `json:"university_signature"`
	DateOfIssuing       string `json:"date_of_issuing"`
	SerialNumber        string `json:"serial_number"`
	RegNo               string `json:"registration_number"`
	CertID              string `json:"cert_id"`
	Version             int    `json:"version"`
	UpdatedDate         string `json:"updated_date"`
}

type CertificateTransactionContext struct {
	contractapi.TransactionContext
}

type CertificateContract struct {
	contractapi.Contract
	info metadata.InfoMetadata
}

func (c *CertificateContract) GetName() string {
	return "CertificateContract"
}

func (c *CertificateContract) GetInfo() metadata.InfoMetadata {
	return c.info
}

func (c *CertificateContract) GetTransactionContextHandler() contractapi.SettableTransactionContextInterface {
	return new(CertificateTransactionContext)
}

func (c *CertificateContract) GetEvaluateTransactions() []string {
	return []string{"ReadCertificate", "CertificateExists"}
}

func (c *CertificateContract) IssueCertificate(ctx *CertificateTransactionContext, certificateJSON string) (string, error) {
	var certificate CertificateOnChain
	if err := json.Unmarshal([]byte(certificateJSON), &certificate); err != nil {
		return "", fmt.Errorf("failed to unmarshal certificate JSON: %v", err)
	}

	exists, err := c.CertificateExists(ctx, certificate.CertID)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("certificate %s already exists", certificate.CertID)
	}

	bytes, err := json.Marshal(certificate)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(certificate.CertID, bytes)
	if err != nil {
		return "", err
	}

	txID := ctx.GetStub().GetTxID()

	eventPayload := map[string]string{
		"certID": certificate.CertID,
		"txID":   txID,
	}
	eventBytes, _ := json.Marshal(eventPayload)
	ctx.GetStub().SetEvent("CertificateIssued", eventBytes)

	return txID, nil
}

func (c *CertificateContract) ReadCertificate(ctx *CertificateTransactionContext, certID string) (*CertificateOnChain, error) {
	bytes, err := ctx.GetStub().GetState(certID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("certificate %s does not exist", certID)
	}

	var certificate CertificateOnChain
	if err := json.Unmarshal(bytes, &certificate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal certificate: %v", err)
	}
	return &certificate, nil
}

func (c *CertificateContract) CertificateExists(ctx *CertificateTransactionContext, certID string) (bool, error) {
	bytes, err := ctx.GetStub().GetState(certID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return bytes != nil, nil
}
func (c *CertificateContract) UpdateCertificate(ctx *CertificateTransactionContext, certificateJSON string) error {
	var updated CertificateOnChain
	if err := json.Unmarshal([]byte(certificateJSON), &updated); err != nil {
		return fmt.Errorf("failed to unmarshal certificate: %v", err)
	}

	exists, err := c.CertificateExists(ctx, updated.CertID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("certificate %s does not exist", updated.CertID)
	}

	updated.Version++
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	updated.UpdatedDate = timestamp.AsTime().Format("2006-01-02T15:04:05Z07:00")

	bytes, err := json.Marshal(updated)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(updated.CertID, bytes)
}

func main() {
	certificateContract := &CertificateContract{
		info: metadata.InfoMetadata{
			Title:   "Certificate Contract",
			Version: "2.0.0",
		},
	}

	chaincode, err := contractapi.NewChaincode(certificateContract)
	if err != nil {
		fmt.Printf("Error creating certificate chaincode: %v\n", err)
		panic(err)
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting certificate chaincode: %v\n", err)
		panic(err)
	}
}
