package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CertificateConfig struct {
	Environment string            `json:"environment"`
	Certificates map[string]CertSettings `json:"certificates"`
}

type CertSettings struct {
	CAFile            string `json:"ca_file"`
	CertFile          string `json:"cert_file"`
	KeyFile           string `json:"key_file"`
	InsecureSkipVerify bool  `json:"insecure_skip_verify"`
	VerifyCertificates bool  `json:"verify_certificates"`
}

type CertificateManager struct {
	config     *CertificateConfig
	configFile string
}

func NewCertificateManager(configFile string) *CertificateManager {
	return &CertificateManager{
		configFile: configFile,
	}
}

func (cm *CertificateManager) LoadConfig() error {
	if _, err := os.Stat(cm.configFile); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		cm.createDefaultConfig()
		return nil
	}
	
	data, err := os.ReadFile(cm.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := json.Unmarshal(data, &cm.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return nil
}

func (cm *CertificateManager) createDefaultConfig() {
	cm.config = &CertificateConfig{
		Environment: "development",
		Certificates: map[string]CertSettings{
			"development": {
				InsecureSkipVerify: true,
				VerifyCertificates: false,
			},
			"staging": {
				InsecureSkipVerify: false,
				VerifyCertificates: true,
				CAFile:            "./certs/staging/ca.crt",
			},
			"production": {
				InsecureSkipVerify: false,
				VerifyCertificates: true,
				CAFile:            "./certs/production/ca.crt",
				CertFile:          "./certs/production/client.crt",
				KeyFile:           "./certs/production/client.key",
			},
		},
	}
	
	cm.SaveConfig()
}

func (cm *CertificateManager) SaveConfig() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(cm.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	if err := os.WriteFile(cm.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func (cm *CertificateManager) GetEnvironment() string {
	return cm.config.Environment
}

func (cm *CertificateManager) SetEnvironment(env string) {
	cm.config.Environment = env
}

func (cm *CertificateManager) GetCertSettings(env string) (*CertSettings, error) {
	settings, exists := cm.config.Certificates[env]
	if !exists {
		return nil, fmt.Errorf("no certificate settings found for environment: %s", env)
	}
	return &settings, nil
}

func (cm *CertificateManager) GetCurrentCertSettings() (*CertSettings, error) {
	return cm.GetCertSettings(cm.config.Environment)
}

func (cm *CertificateManager) SetCertSettings(env string, settings CertSettings) {
	cm.config.Certificates[env] = settings
}

func (cm *CertificateManager) AddEnvironment(env string, settings CertSettings) {
	cm.config.Certificates[env] = settings
}

func (cm *CertificateManager) RemoveEnvironment(env string) {
	delete(cm.config.Certificates, env)
}

func (cm *CertificateManager) ListEnvironments() []string {
	var environments []string
	for env := range cm.config.Certificates {
		environments = append(environments, env)
	}
	return environments
}

func (cm *CertificateManager) ValidateCertFiles(env string) error {
	settings, err := cm.GetCertSettings(env)
	if err != nil {
		return err
	}
	
	// Check CA file
	if settings.CAFile != "" {
		if _, err := os.Stat(settings.CAFile); os.IsNotExist(err) {
			return fmt.Errorf("CA file does not exist: %s", settings.CAFile)
		}
	}
	
	// Check client certificate files
	if settings.CertFile != "" {
		if _, err := os.Stat(settings.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("certificate file does not exist: %s", settings.CertFile)
		}
	}
	
	if settings.KeyFile != "" {
		if _, err := os.Stat(settings.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("key file does not exist: %s", settings.KeyFile)
		}
	}
	
	return nil
}

func (cm *CertificateManager) GetConfig() *CertificateConfig {
	return cm.config
}
