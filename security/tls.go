package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TLSConfig struct {
	VerifyCertificates bool
	CertFile          string
	KeyFile           string
	CAFile            string
	InsecureSkipVerify bool
}

type TLSManager struct {
	config *TLSConfig
	client *http.Client
}

func NewTLSManager(config *TLSConfig) *TLSManager {
	tm := &TLSManager{
		config: config,
	}
	
	tm.setupClient()
	return tm
}

func (tm *TLSManager) setupClient() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: tm.config.InsecureSkipVerify,
	}
	
	// Load custom CA certificate if provided
	if tm.config.CAFile != "" {
		caCert, err := ioutil.ReadFile(tm.config.CAFile)
		if err != nil {
			fmt.Printf("Warning: Failed to load CA certificate: %v\n", err)
		} else {
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				fmt.Printf("Warning: Failed to parse CA certificate\n")
			} else {
				tlsConfig.RootCAs = caCertPool
			}
		}
	}
	
	// Load client certificate if provided
	if tm.config.CertFile != "" && tm.config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(tm.config.CertFile, tm.config.KeyFile)
		if err != nil {
			fmt.Printf("Warning: Failed to load client certificate: %v\n", err)
		} else {
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}
	
	// Set up HTTP client with custom transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	
	tm.client = &http.Client{
		Transport: transport,
	}
}

func (tm *TLSManager) GetClient() *http.Client {
	return tm.client
}

func (tm *TLSManager) ValidateCertificate(hostname string) error {
	if !tm.config.VerifyCertificates {
		return nil
	}
	
	// Create a temporary client for certificate validation
	tempClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	
	// Make a test request to validate the certificate
	resp, err := tempClient.Get("https://" + hostname)
	if err != nil {
		return fmt.Errorf("certificate validation failed: %w", err)
	}
	defer resp.Body.Close()
	
	return nil
}

func (tm *TLSManager) SetInsecureSkipVerify(skip bool) {
	tm.config.InsecureSkipVerify = skip
	tm.setupClient()
}

func (tm *TLSManager) LoadCustomCA(caFile string) error {
	if _, err := os.Stat(caFile); os.IsNotExist(err) {
		return fmt.Errorf("CA file does not exist: %s", caFile)
	}
	
	tm.config.CAFile = caFile
	tm.setupClient()
	return nil
}

func (tm *TLSManager) LoadClientCertificate(certFile, keyFile string) error {
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return fmt.Errorf("certificate file does not exist: %s", certFile)
	}
	
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf("key file does not exist: %s", keyFile)
	}
	
	tm.config.CertFile = certFile
	tm.config.KeyFile = keyFile
	tm.setupClient()
	return nil
}

func (tm *TLSManager) GetConfig() *TLSConfig {
	return tm.config
}

func (tm *TLSManager) SetVerifyCertificates(verify bool) {
	tm.config.VerifyCertificates = verify
	if !verify {
		tm.config.InsecureSkipVerify = true
	} else {
		tm.config.InsecureSkipVerify = false
	}
	tm.setupClient()
}
