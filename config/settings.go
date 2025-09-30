package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Settings struct {
	DefaultURL      string            `json:"default_url"`
	DefaultMethod   string            `json:"default_method"`
	DefaultHeaders  map[string]string `json:"default_headers"`
	DefaultTimeout  time.Duration     `json:"default_timeout"`
	DefaultRetries  int               `json:"default_retries"`
	DefaultDelay    time.Duration     `json:"default_delay"`
	ProxyURL        string            `json:"proxy_url"`
	FollowRedirects bool              `json:"follow_redirects"`
	VerifySSL       bool              `json:"verify_ssl"`
	SaveResponses   bool              `json:"save_responses"`
	ResponseDir     string            `json:"response_dir"`
	LogLevel        string            `json:"log_level"`
	ColorOutput     bool              `json:"color_output"`
	InteractiveMode bool              `json:"interactive_mode"`
	HistorySize     int               `json:"history_size"`
	AutoComplete    bool              `json:"auto_complete"`
	Theme           string            `json:"theme"`
	Aliases         map[string]string `json:"aliases"`
	Environments    map[string]Environment `json:"environments"`
}

type Environment struct {
	Name     string            `json:"name"`
	BaseURL  string            `json:"base_url"`
	Headers  map[string]string `json:"headers"`
	Timeout  time.Duration     `json:"timeout"`
	Retries  int               `json:"retries"`
	Delay    time.Duration     `json:"delay"`
	ProxyURL string            `json:"proxy_url"`
	VerifySSL bool             `json:"verify_ssl"`
}

type ConfigManager struct {
	settings    *Settings
	configFile  string
	configDir   string
}

func NewConfigManager() *ConfigManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".api_test")
	configFile := filepath.Join(configDir, "config.json")
	
	return &ConfigManager{
		configFile: configFile,
		configDir:  configDir,
	}
}

func (cm *ConfigManager) Load() error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(cm.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Load config file if it exists
	if _, err := os.Stat(cm.configFile); os.IsNotExist(err) {
		cm.settings = cm.getDefaultSettings()
		return cm.Save()
	}
	
	data, err := os.ReadFile(cm.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := json.Unmarshal(data, &cm.settings); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return nil
}

func (cm *ConfigManager) Save() error {
	data, err := json.MarshalIndent(cm.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func (cm *ConfigManager) GetSettings() *Settings {
	return cm.settings
}

func (cm *ConfigManager) SetSettings(settings *Settings) {
	cm.settings = settings
}

func (cm *ConfigManager) Get(key string) (interface{}, error) {
	switch key {
	case "default_url":
		return cm.settings.DefaultURL, nil
	case "default_method":
		return cm.settings.DefaultMethod, nil
	case "default_timeout":
		return cm.settings.DefaultTimeout, nil
	case "default_retries":
		return cm.settings.DefaultRetries, nil
	case "default_delay":
		return cm.settings.DefaultDelay, nil
	case "proxy_url":
		return cm.settings.ProxyURL, nil
	case "follow_redirects":
		return cm.settings.FollowRedirects, nil
	case "verify_ssl":
		return cm.settings.VerifySSL, nil
	case "save_responses":
		return cm.settings.SaveResponses, nil
	case "response_dir":
		return cm.settings.ResponseDir, nil
	case "log_level":
		return cm.settings.LogLevel, nil
	case "color_output":
		return cm.settings.ColorOutput, nil
	case "interactive_mode":
		return cm.settings.InteractiveMode, nil
	case "history_size":
		return cm.settings.HistorySize, nil
	case "auto_complete":
		return cm.settings.AutoComplete, nil
	case "theme":
		return cm.settings.Theme, nil
	default:
		return nil, fmt.Errorf("unknown setting: %s", key)
	}
}

func (cm *ConfigManager) Set(key string, value interface{}) error {
	switch key {
	case "default_url":
		if v, ok := value.(string); ok {
			cm.settings.DefaultURL = v
		} else {
			return fmt.Errorf("invalid value type for default_url")
		}
	case "default_method":
		if v, ok := value.(string); ok {
			cm.settings.DefaultMethod = v
		} else {
			return fmt.Errorf("invalid value type for default_method")
		}
	case "default_timeout":
		if v, ok := value.(time.Duration); ok {
			cm.settings.DefaultTimeout = v
		} else {
			return fmt.Errorf("invalid value type for default_timeout")
		}
	case "default_retries":
		if v, ok := value.(int); ok {
			cm.settings.DefaultRetries = v
		} else {
			return fmt.Errorf("invalid value type for default_retries")
		}
	case "default_delay":
		if v, ok := value.(time.Duration); ok {
			cm.settings.DefaultDelay = v
		} else {
			return fmt.Errorf("invalid value type for default_delay")
		}
	case "proxy_url":
		if v, ok := value.(string); ok {
			cm.settings.ProxyURL = v
		} else {
			return fmt.Errorf("invalid value type for proxy_url")
		}
	case "follow_redirects":
		if v, ok := value.(bool); ok {
			cm.settings.FollowRedirects = v
		} else {
			return fmt.Errorf("invalid value type for follow_redirects")
		}
	case "verify_ssl":
		if v, ok := value.(bool); ok {
			cm.settings.VerifySSL = v
		} else {
			return fmt.Errorf("invalid value type for verify_ssl")
		}
	case "save_responses":
		if v, ok := value.(bool); ok {
			cm.settings.SaveResponses = v
		} else {
			return fmt.Errorf("invalid value type for save_responses")
		}
	case "response_dir":
		if v, ok := value.(string); ok {
			cm.settings.ResponseDir = v
		} else {
			return fmt.Errorf("invalid value type for response_dir")
		}
	case "log_level":
		if v, ok := value.(string); ok {
			cm.settings.LogLevel = v
		} else {
			return fmt.Errorf("invalid value type for log_level")
		}
	case "color_output":
		if v, ok := value.(bool); ok {
			cm.settings.ColorOutput = v
		} else {
			return fmt.Errorf("invalid value type for color_output")
		}
	case "interactive_mode":
		if v, ok := value.(bool); ok {
			cm.settings.InteractiveMode = v
		} else {
			return fmt.Errorf("invalid value type for interactive_mode")
		}
	case "history_size":
		if v, ok := value.(int); ok {
			cm.settings.HistorySize = v
		} else {
			return fmt.Errorf("invalid value type for history_size")
		}
	case "auto_complete":
		if v, ok := value.(bool); ok {
			cm.settings.AutoComplete = v
		} else {
			return fmt.Errorf("invalid value type for auto_complete")
		}
	case "theme":
		if v, ok := value.(string); ok {
			cm.settings.Theme = v
		} else {
			return fmt.Errorf("invalid value type for theme")
		}
	default:
		return fmt.Errorf("unknown setting: %s", key)
	}
	
	return nil
}

func (cm *ConfigManager) AddAlias(name, value string) {
	if cm.settings.Aliases == nil {
		cm.settings.Aliases = make(map[string]string)
	}
	cm.settings.Aliases[name] = value
}

func (cm *ConfigManager) RemoveAlias(name string) {
	if cm.settings.Aliases != nil {
		delete(cm.settings.Aliases, name)
	}
}

func (cm *ConfigManager) GetAlias(name string) (string, bool) {
	if cm.settings.Aliases == nil {
		return "", false
	}
	value, exists := cm.settings.Aliases[name]
	return value, exists
}

func (cm *ConfigManager) ListAliases() map[string]string {
	if cm.settings.Aliases == nil {
		return make(map[string]string)
	}
	return cm.settings.Aliases
}

func (cm *ConfigManager) AddEnvironment(name string, env Environment) {
	if cm.settings.Environments == nil {
		cm.settings.Environments = make(map[string]Environment)
	}
	cm.settings.Environments[name] = env
}

func (cm *ConfigManager) RemoveEnvironment(name string) {
	if cm.settings.Environments != nil {
		delete(cm.settings.Environments, name)
	}
}

func (cm *ConfigManager) GetEnvironment(name string) (Environment, bool) {
	if cm.settings.Environments == nil {
		return Environment{}, false
	}
	env, exists := cm.settings.Environments[name]
	return env, exists
}

func (cm *ConfigManager) ListEnvironments() map[string]Environment {
	if cm.settings.Environments == nil {
		return make(map[string]Environment)
	}
	return cm.settings.Environments
}

func (cm *ConfigManager) getDefaultSettings() *Settings {
	return &Settings{
		DefaultURL:      "https://httpbin.org",
		DefaultMethod:   "GET",
		DefaultHeaders:  map[string]string{"User-Agent": "APITestProgram/1.0"},
		DefaultTimeout:  30 * time.Second,
		DefaultRetries:  3,
		DefaultDelay:    1 * time.Second,
		ProxyURL:        "",
		FollowRedirects: true,
		VerifySSL:       true,
		SaveResponses:   false,
		ResponseDir:     "./responses",
		LogLevel:        "INFO",
		ColorOutput:     true,
		InteractiveMode: true,
		HistorySize:     100,
		AutoComplete:    true,
		Theme:           "default",
		Aliases:         make(map[string]string),
		Environments:    make(map[string]Environment),
	}
}

func (cm *ConfigManager) Reset() {
	cm.settings = cm.getDefaultSettings()
}

func (cm *ConfigManager) Export(filename string) error {
	data, err := json.MarshalIndent(cm.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}
	
	return nil
}

func (cm *ConfigManager) Import(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}
	
	if err := json.Unmarshal(data, &cm.settings); err != nil {
		return fmt.Errorf("failed to parse import file: %w", err)
	}
	
	return nil
}
