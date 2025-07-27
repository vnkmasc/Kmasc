package ca

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func GenerateRootCA(caKey, caCert string) error {
	cmd := exec.Command("openssl", "ecparam", "-name", "prime256v1",
		"-genkey", "-noout", "-out", caKey)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("openssl", "req", "-new", "-x509", "-sha256",
		"-key", caKey, "-out", caCert, "-days", "3650",
		"-subj", "/CN=CustomCA/O=Org1")
	return cmd.Run()
}

func GenerateUserKeyCSR(keyFile, csrFile, cn, ou string) error {
	cmd := exec.Command("openssl", "ecparam", "-name", "prime256v1",
		"-genkey", "-noout", "-out", keyFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	subj := fmt.Sprintf("/CN=%s", cn)
	if ou != "" {
		subj += fmt.Sprintf("/OU=%s", ou)
	}
	cmd = exec.Command("openssl", "req", "-new", "-key", keyFile,
		"-out", csrFile, "-subj", subj)
	return cmd.Run()
}

func SignCSR(csr, caKey, caCert, outCert string) error {
	cmd := exec.Command("openssl", "x509", "-req", "-in", csr,
		"-CA", caCert, "-CAkey", caKey, "-CAcreateserial",
		"-out", outCert, "-days", "365", "-sha256")
	return cmd.Run()
}

func GenerateUserKeyCSRWithSAN(keyFile, csrFile, cn, ou string, hosts []string) error {
	if err := exec.Command("openssl", "ecparam", "-name", "prime256v1",
		"-genkey", "-noout", "-out", keyFile).Run(); err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}
	log.Println("GenerateUserKeyCSRWithSAN", hosts)
	tmpConf, err := os.CreateTemp("", "csr_san_*.conf")
	if err != nil {
		return err
	}
	defer os.Remove(tmpConf.Name())

	altNames := ""
	for i, h := range hosts {
		altNames += fmt.Sprintf("DNS.%d = %s\n", i+1, h)
	}

	ouLine := ""
	if ou != "" {
		ouLine = fmt.Sprintf("OU=%s\n", ou)
	}

	confContent := fmt.Sprintf(`
[req]
distinguished_name=req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
CN=%s
%s

[v3_req]
subjectAltName = @alt_names

[alt_names]
%s
`, cn, ouLine, altNames)

	if err := os.WriteFile(tmpConf.Name(), []byte(confContent), 0644); err != nil {
		return err
	}

	cmd := exec.Command("openssl", "req", "-new", "-key", keyFile,
		"-out", csrFile, "-config", tmpConf.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to generate csr: %v, output: %s", err, out)
	}

	return nil
}

func SignCSRWithSAN(csr, caKey, caCert, outCert string, hosts []string) error {
	conf := csr + ".conf"
	altNames := ""
	for i, h := range hosts {
		altNames += fmt.Sprintf("DNS.%d = %s\n", i+1, h)
	}

	log.Println("SignCSRWithSAN", hosts)
	confContent := fmt.Sprintf("[v3_req]\nsubjectAltName = @alt_names\n[alt_names]\n%s", altNames)
	os.WriteFile(conf, []byte(confContent), 0644)

	cmd := exec.Command("openssl", "x509", "-req", "-in", csr,
		"-CA", caCert, "-CAkey", caKey, "-CAcreateserial",
		"-out", outCert, "-days", "365", "-sha256",
		"-extfile", conf, "-extensions", "v3_req")
	return cmd.Run()
}
