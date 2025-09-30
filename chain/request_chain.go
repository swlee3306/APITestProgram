package chain

import (
	"APITestProgram/http"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ChainStep struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Timeout     time.Duration     `json:"timeout"`
	Retries     int               `json:"retries"`
	Delay       time.Duration     `json:"delay"`
	Variables   map[string]string `json:"variables"`
	Conditions  []Condition       `json:"conditions"`
	Extract     []ExtractRule     `json:"extract"`
	Wait        time.Duration     `json:"wait"`
	Enabled     bool              `json:"enabled"`
}

type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Action   string `json:"action"` // "continue", "skip", "stop"
}

type ExtractRule struct {
	Name     string `json:"name"`
	Source   string `json:"source"` // "header", "body", "status"
	Pattern  string `json:"pattern"`
	Required bool   `json:"required"`
}

type ChainResult struct {
	StepID     string                 `json:"step_id"`
	Success    bool                   `json:"success"`
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       string                 `json:"body"`
	Duration   time.Duration          `json:"duration"`
	Error      string                 `json:"error,omitempty"`
	Variables  map[string]string      `json:"variables"`
	Extracted  map[string]string      `json:"extracted"`
}

type RequestChain struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Steps       []ChainStep `json:"steps"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ChainManager struct {
	chains    []RequestChain
	client    *http.HTTPClient
	variables map[string]string
}

func NewChainManager() *ChainManager {
	return &ChainManager{
		chains:    make([]RequestChain, 0),
		client:    http.NewHTTPClient(),
		variables: make(map[string]string),
	}
}

func (cm *ChainManager) CreateChain(name, description string) *RequestChain {
	chain := RequestChain{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Steps:       make([]ChainStep, 0),
		Variables:   make(map[string]string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	cm.chains = append(cm.chains, chain)
	return &chain
}

func (cm *ChainManager) GetChain(id string) (*RequestChain, error) {
	for _, chain := range cm.chains {
		if chain.ID == id {
			return &chain, nil
		}
	}
	return nil, fmt.Errorf("chain not found: %s", id)
}

func (cm *ChainManager) ListChains() []RequestChain {
	return cm.chains
}

func (cm *ChainManager) DeleteChain(id string) error {
	for i, chain := range cm.chains {
		if chain.ID == id {
			cm.chains = append(cm.chains[:i], cm.chains[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("chain not found: %s", id)
}

func (cm *ChainManager) AddStep(chainID string, step ChainStep) error {
	for i, chain := range cm.chains {
		if chain.ID == chainID {
			step.ID = generateID()
			step.Enabled = true
			cm.chains[i].Steps = append(cm.chains[i].Steps, step)
			cm.chains[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("chain not found: %s", chainID)
}

func (cm *ChainManager) RemoveStep(chainID, stepID string) error {
	for i, chain := range cm.chains {
		if chain.ID == chainID {
			for j, step := range chain.Steps {
				if step.ID == stepID {
					cm.chains[i].Steps = append(cm.chains[i].Steps[:j], cm.chains[i].Steps[j+1:]...)
					cm.chains[i].UpdatedAt = time.Now()
					return nil
				}
			}
			return fmt.Errorf("step not found: %s", stepID)
		}
	}
	return fmt.Errorf("chain not found: %s", chainID)
}

func (cm *ChainManager) UpdateStep(chainID, stepID string, updates map[string]interface{}) error {
	for i, chain := range cm.chains {
		if chain.ID == chainID {
			for j, step := range chain.Steps {
				if step.ID == stepID {
					// Update fields
					for key, value := range updates {
						switch key {
						case "name":
							if v, ok := value.(string); ok {
								cm.chains[i].Steps[j].Name = v
							}
						case "method":
							if v, ok := value.(string); ok {
								cm.chains[i].Steps[j].Method = v
							}
						case "url":
							if v, ok := value.(string); ok {
								cm.chains[i].Steps[j].URL = v
							}
						case "headers":
							if v, ok := value.(map[string]string); ok {
								cm.chains[i].Steps[j].Headers = v
							}
						case "body":
							if v, ok := value.(string); ok {
								cm.chains[i].Steps[j].Body = v
							}
						case "enabled":
							if v, ok := value.(bool); ok {
								cm.chains[i].Steps[j].Enabled = v
							}
						}
					}
					cm.chains[i].UpdatedAt = time.Now()
					return nil
				}
			}
			return fmt.Errorf("step not found: %s", stepID)
		}
	}
	return fmt.Errorf("chain not found: %s", chainID)
}

func (cm *ChainManager) ExecuteChain(chainID string) ([]ChainResult, error) {
	chain, err := cm.GetChain(chainID)
	if err != nil {
		return nil, err
	}
	
	var results []ChainResult
	variables := make(map[string]string)
	
	// Copy chain variables
	for k, v := range chain.Variables {
		variables[k] = v
	}
	
	// Copy global variables
	for k, v := range cm.variables {
		variables[k] = v
	}
	
	for _, step := range chain.Steps {
		if !step.Enabled {
			continue
		}
		
		// Replace variables in step
		step = cm.replaceVariables(step, variables)
		
		// Execute step
		result, err := cm.executeStep(step, variables)
		if err != nil {
			result.Error = err.Error()
			result.Success = false
		}
		
		results = append(results, result)
		
		// Update variables with extracted values
		for k, v := range result.Extracted {
			variables[k] = v
		}
		
		// Check conditions
		if shouldStop := cm.checkConditions(step.Conditions, result); shouldStop {
			break
		}
		
		// Wait if specified
		if step.Wait > 0 {
			time.Sleep(step.Wait)
		}
	}
	
	return results, nil
}

func (cm *ChainManager) executeStep(step ChainStep, variables map[string]string) (ChainResult, error) {
	start := time.Now()
	
	// Create request options
	options := http.RequestOptions{
		Method:  http.HTTPMethod(step.Method),
		URL:     step.URL,
		Headers: step.Headers,
		Body:    step.Body,
		Timeout: step.Timeout,
		Retries: step.Retries,
		Delay:   step.Delay,
	}
	
	// Execute request
	resp, err := cm.client.Request(options)
	if err != nil {
		return ChainResult{
			StepID:   step.ID,
			Success:  false,
			Error:    err.Error(),
			Duration: time.Since(start),
		}, err
	}
	
	// Extract data
	extracted := make(map[string]string)
	for _, rule := range step.Extract {
		value, err := cm.extractValue(rule, resp)
		if err != nil {
			if rule.Required {
				return ChainResult{
					StepID:   step.ID,
					Success:  false,
					Error:    fmt.Sprintf("failed to extract required field %s: %v", rule.Name, err),
					Duration: time.Since(start),
				}, err
			}
			continue
		}
		extracted[rule.Name] = value
	}
	
	return ChainResult{
		StepID:     step.ID,
		Success:    resp.Success,
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       resp.Body,
		Duration:   time.Since(start),
		Extracted:  extracted,
	}, nil
}

func (cm *ChainManager) replaceVariables(step ChainStep, variables map[string]string) ChainStep {
	// Replace variables in URL
	step.URL = cm.replaceString(step.URL, variables)
	
	// Replace variables in headers
	for k, v := range step.Headers {
		step.Headers[k] = cm.replaceString(v, variables)
	}
	
	// Replace variables in body
	step.Body = cm.replaceString(step.Body, variables)
	
	return step
}

func (cm *ChainManager) replaceString(str string, variables map[string]string) string {
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		str = strings.ReplaceAll(str, placeholder, value)
	}
	return str
}

func (cm *ChainManager) extractValue(rule ExtractRule, resp *http.Response) (string, error) {
	var source string
	
	switch rule.Source {
	case "header":
		if value, exists := resp.Headers[rule.Pattern]; exists {
			source = value
		} else {
			return "", fmt.Errorf("header not found: %s", rule.Pattern)
		}
	case "body":
		source = resp.Body
	case "status":
		source = fmt.Sprintf("%d", resp.StatusCode)
	default:
		return "", fmt.Errorf("unknown source: %s", rule.Source)
	}
	
	// Simple pattern matching (could be enhanced with regex)
	if rule.Pattern != "" && rule.Pattern != rule.Source {
		// For now, just return the source value
		// In a real implementation, you'd use regex or JSON path
		return source, nil
	}
	
	return source, nil
}

func (cm *ChainManager) checkConditions(conditions []Condition, result ChainResult) bool {
	for _, condition := range conditions {
		var value string
		
		switch condition.Field {
		case "status_code":
			value = fmt.Sprintf("%d", result.StatusCode)
		case "success":
			value = fmt.Sprintf("%t", result.Success)
		case "duration":
			value = fmt.Sprintf("%v", result.Duration)
		default:
			continue
		}
		
		if !cm.evaluateCondition(value, condition.Operator, condition.Value) {
			if condition.Action == "stop" {
				return true
			}
		}
	}
	
	return false
}

func (cm *ChainManager) evaluateCondition(value, operator, expected string) bool {
	switch operator {
	case "equals":
		return value == expected
	case "not_equals":
		return value != expected
	case "contains":
		return strings.Contains(value, expected)
	case "not_contains":
		return !strings.Contains(value, expected)
	case "greater_than":
		// For numeric comparison, convert to float64
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			if e, err := strconv.ParseFloat(expected, 64); err == nil {
				return v > e
			}
		}
		return value > expected
	case "less_than":
		// For numeric comparison, convert to float64
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			if e, err := strconv.ParseFloat(expected, 64); err == nil {
				return v < e
			}
		}
		return value < expected
	default:
		return false
	}
}

func (cm *ChainManager) SetGlobalVariable(key, value string) {
	cm.variables[key] = value
}

func (cm *ChainManager) GetGlobalVariable(key string) (string, bool) {
	value, exists := cm.variables[key]
	return value, exists
}

func (cm *ChainManager) ClearGlobalVariables() {
	cm.variables = make(map[string]string)
}

func (cm *ChainManager) ExportChain(chainID, filename string) error {
	chain, err := cm.GetChain(chainID)
	if err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chain: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write chain file: %w", err)
	}
	
	return nil
}

func (cm *ChainManager) ImportChain(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read chain file: %w", err)
	}
	
	var chain RequestChain
	if err := json.Unmarshal(data, &chain); err != nil {
		return fmt.Errorf("failed to parse chain file: %w", err)
	}
	
	chain.ID = generateID()
	chain.UpdatedAt = time.Now()
	
	cm.chains = append(cm.chains, chain)
	return nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
