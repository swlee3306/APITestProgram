package chain

import (
	"APITestProgram/http"
	"os"
	"testing"
	"time"
)

func TestChainManager(t *testing.T) {
	cm := NewChainManager()
	
	// Test creating a chain
	chain := cm.CreateChain("Test Chain", "A test chain for API testing")
	if chain == nil {
		t.Fatal("Failed to create chain")
	}
	
	if chain.Name != "Test Chain" {
		t.Errorf("Expected chain name 'Test Chain', got '%s'", chain.Name)
	}
	
	// Test adding a step
	step := ChainStep{
		Name:   "First Step",
		Method: "GET",
		URL:    "https://httpbin.org/get",
		Headers: map[string]string{
			"User-Agent": "APITestProgram",
		},
		Timeout: 30 * time.Second,
		Retries: 3,
		Delay:   1 * time.Second,
		Extract: []ExtractRule{
			{
				Name:     "status",
				Source:   "status",
				Required: true,
			},
		},
		Conditions: []Condition{
			{
				Field:    "status_code",
				Operator: "equals",
				Value:    "200",
				Action:   "continue",
			},
		},
	}
	
	err := cm.AddStep(chain.ID, step)
	if err != nil {
		t.Fatalf("Failed to add step: %v", err)
	}
	
	// Test getting chain
	retrievedChain, err := cm.GetChain(chain.ID)
	if err != nil {
		t.Fatalf("Failed to get chain: %v", err)
	}
	
	if len(retrievedChain.Steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(retrievedChain.Steps))
	}
	
	// Test listing chains
	chains := cm.ListChains()
	if len(chains) != 1 {
		t.Errorf("Expected 1 chain, got %d", len(chains))
	}
	
	// Test updating step
	updates := map[string]interface{}{
		"name": "Updated Step",
		"url":  "https://httpbin.org/post",
	}
	
	err = cm.UpdateStep(chain.ID, retrievedChain.Steps[0].ID, updates)
	if err != nil {
		t.Fatalf("Failed to update step: %v", err)
	}
	
	// Test global variables
	cm.SetGlobalVariable("base_url", "https://httpbin.org")
	value, exists := cm.GetGlobalVariable("base_url")
	if !exists || value != "https://httpbin.org" {
		t.Errorf("Expected global variable 'base_url' to be 'https://httpbin.org', got '%s'", value)
	}
	
	// Test export/import
	tempFile := "test_chain.json"
	defer os.Remove(tempFile)
	
	err = cm.ExportChain(chain.ID, tempFile)
	if err != nil {
		t.Fatalf("Failed to export chain: %v", err)
	}
	
	// Create new manager and import
	cm2 := NewChainManager()
	err = cm2.ImportChain(tempFile)
	if err != nil {
		t.Fatalf("Failed to import chain: %v", err)
	}
	
	importedChains := cm2.ListChains()
	if len(importedChains) != 1 {
		t.Errorf("Expected 1 imported chain, got %d", len(importedChains))
	}
	
	// Test deleting chain
	err = cm.DeleteChain(chain.ID)
	if err != nil {
		t.Fatalf("Failed to delete chain: %v", err)
	}
	
	chains = cm.ListChains()
	if len(chains) != 0 {
		t.Errorf("Expected 0 chains after deletion, got %d", len(chains))
	}
}

func TestChainStepValidation(t *testing.T) {
	step := ChainStep{
		Name:   "Test Step",
		Method: "GET",
		URL:    "https://httpbin.org/get",
		Timeout: 30 * time.Second,
		Retries: 3,
		Delay:   1 * time.Second,
	}
	
	// Test step validation
	if step.Name == "" {
		t.Error("Step name should not be empty")
	}
	
	if step.Method == "" {
		t.Error("Step method should not be empty")
	}
	
	if step.URL == "" {
		t.Error("Step URL should not be empty")
	}
	
	if step.Timeout <= 0 {
		t.Error("Step timeout should be positive")
	}
	
	if step.Retries < 0 {
		t.Error("Step retries should be non-negative")
	}
	
	if step.Delay < 0 {
		t.Error("Step delay should be non-negative")
	}
}

func TestConditionEvaluation(t *testing.T) {
	cm := NewChainManager()
	
	tests := []struct {
		value    string
		operator string
		expected string
		result   bool
	}{
		{"200", "equals", "200", true},
		{"200", "equals", "404", false},
		{"200", "not_equals", "404", true},
		{"hello world", "contains", "world", true},
		{"hello world", "contains", "universe", false},
		{"hello world", "not_contains", "universe", true},
		{"10", "greater_than", "5", true},
		{"5", "greater_than", "10", false},
		{"5", "less_than", "10", true},
		{"10", "less_than", "5", false},
	}
	
	for _, test := range tests {
		result := cm.evaluateCondition(test.value, test.operator, test.expected)
		if result != test.result {
			t.Errorf("evaluateCondition(%s, %s, %s) = %v, expected %v", 
				test.value, test.operator, test.expected, result, test.result)
		}
	}
}

func TestVariableReplacement(t *testing.T) {
	cm := NewChainManager()
	
	variables := map[string]string{
		"base_url": "https://httpbin.org",
		"version":  "v1",
		"endpoint": "users",
	}
	
	step := ChainStep{
		URL: "{{base_url}}/{{version}}/{{endpoint}}",
		Headers: map[string]string{
			"Authorization": "Bearer {{token}}",
			"Content-Type":  "application/json",
		},
		Body: `{"user": "{{username}}", "email": "{{email}}"}`,
	}
	
	replaced := cm.replaceVariables(step, variables)
	
	expectedURL := "https://httpbin.org/v1/users"
	if replaced.URL != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, replaced.URL)
	}
	
	expectedAuth := "Bearer {{token}}"
	if replaced.Headers["Authorization"] != expectedAuth {
		t.Errorf("Expected Authorization header '%s', got '%s'", expectedAuth, replaced.Headers["Authorization"])
	}
	
	expectedBody := `{"user": "{{username}}", "email": "{{email}}"}`
	if replaced.Body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, replaced.Body)
	}
}

func TestExtractRules(t *testing.T) {
	cm := NewChainManager()
	
	// Mock response
	resp := &http.Response{
		Success:    true,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Rate-Limit": "100",
		},
		Body: `{"id": 123, "name": "John Doe", "email": "john@example.com"}`,
	}
	
	tests := []struct {
		rule     ExtractRule
		expected string
		hasError bool
	}{
		{
			rule:     ExtractRule{Name: "status", Source: "status", Required: true},
			expected: "200",
			hasError: false,
		},
		{
			rule:     ExtractRule{Name: "content_type", Source: "header", Pattern: "Content-Type", Required: true},
			expected: "application/json",
			hasError: false,
		},
		{
			rule:     ExtractRule{Name: "rate_limit", Source: "header", Pattern: "X-Rate-Limit", Required: true},
			expected: "100",
			hasError: false,
		},
		{
			rule:     ExtractRule{Name: "body", Source: "body", Required: true},
			expected: `{"id": 123, "name": "John Doe", "email": "john@example.com"}`,
			hasError: false,
		},
	}
	
	for _, test := range tests {
		value, err := cm.extractValue(test.rule, resp)
		if test.hasError && err == nil {
			t.Errorf("Expected error for rule %s, got none", test.rule.Name)
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for rule %s: %v", test.rule.Name, err)
		}
		if !test.hasError && value != test.expected {
			t.Errorf("Expected value '%s' for rule %s, got '%s'", test.expected, test.rule.Name, value)
		}
	}
}
