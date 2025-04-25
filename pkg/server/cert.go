package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

func CheckCert(certFile, keyFile string) error {
	// Check if the certificate and key files exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Printf("Certificate file '%s' does not exist; creating self-signed cert", certFile)
		return CreateSelfSignedCert(certFile, keyFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Printf("Key file '%s' does not exist; creating self-signed cert", certFile)
		return CreateSelfSignedCert(certFile, keyFile)
	}

	// Load the certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load certificate and key: %v", err)
	}

	// Check if the certificate is unexpired
	if time.Now().After(cert.Leaf.NotAfter) {
		log.Printf("Certificate has expired: %s", cert.Leaf.NotAfter)
		if cert.Leaf.Subject.String() == cert.Leaf.Issuer.String() {
			log.Printf("Certificate is self-signed; creating a new self-signed cert")
			return CreateSelfSignedCert(certFile, keyFile)
		} else {
			return fmt.Errorf("certificate is not self-signed and has expired: %s", cert.Leaf.NotAfter)
		}
	}
	return nil
}

func CreateSelfSignedCert(certFile, keyFile string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ECDSA key: %v", err)
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number for self-signed cert: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %v", err)
	}
	dnsNames := []string{hostname}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return fmt.Errorf("failed to create self-signed certificate: %v", err)
	}
	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", certFile, err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("failed to write data to %s: %v", certFile, err)
	}
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %v", keyFile, err)
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write data to %s: %v", keyFile, err)
	}
	log.Printf("Created self-signed certificate: %s", certFile)
	log.Printf("Created private key: %s", keyFile)
	return nil
}
