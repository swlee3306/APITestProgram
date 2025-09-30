# Request Chain Functionality

This package provides a powerful request chaining system for APITestProgram, allowing you to create, manage, and execute sequences of HTTP requests with variable passing, conditional logic, and data extraction.

## Features

- **Chain Management**: Create, update, and delete request chains
- **Step Configuration**: Define HTTP requests with headers, body, timeout, and retry settings
- **Variable Support**: Pass data between steps using variables
- **Data Extraction**: Extract values from responses (headers, body, status)
- **Conditional Logic**: Control flow based on response conditions
- **Export/Import**: Save and load chains as JSON files
- **Global Variables**: Set variables that persist across all chains
- **Interactive CLI**: User-friendly command-line interface

## Quick Start

```go
package main

import (
    "APITestProgram/chain"
    "fmt"
)

func main() {
    // Create a new chain manager
    cm := chain.NewChainManager()
    
    // Create a new chain
    chain := cm.CreateChain("API Test Flow", "Test user registration and login")
    
    // Add a registration step
    registrationStep := chain.ChainStep{
        Name:   "Register User",
        Method: "POST",
        URL:    "https://api.example.com/register",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: `{"username": "testuser", "email": "test@example.com"}`,
        Extract: []chain.ExtractRule{
            {
                Name:     "user_id",
                Source:   "body",
                Pattern:  "id",
                Required: true,
            },
        },
    }
    
    // Add a login step that uses the extracted user ID
    loginStep := chain.ChainStep{
        Name:   "Login User",
        Method: "POST",
        URL:    "https://api.example.com/login",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: `{"user_id": "{{user_id}}", "password": "password123"}`,
        Extract: []chain.ExtractRule{
            {
                Name:     "auth_token",
                Source:   "body",
                Pattern:  "token",
                Required: true,
            },
        },
    }
    
    // Add steps to chain
    cm.AddStep(chain.ID, registrationStep)
    cm.AddStep(chain.ID, loginStep)
    
    // Execute the chain
    results, err := cm.ExecuteChain(chain.ID)
    if err != nil {
        fmt.Printf("Chain execution failed: %v\n", err)
        return
    }
    
    // Print results
    for i, result := range results {
        fmt.Printf("Step %d: %s\n", i+1, result.StepID)
        fmt.Printf("  Status: %d\n", result.StatusCode)
        fmt.Printf("  Success: %t\n", result.Success)
        if len(result.Extracted) > 0 {
            fmt.Println("  Extracted Variables:")
            for k, v := range result.Extracted {
                fmt.Printf("    %s = %s\n", k, v)
            }
        }
    }
}
```

## Chain Structure

### RequestChain
```go
type RequestChain struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Steps       []ChainStep       `json:"steps"`
    Variables   map[string]string `json:"variables"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}
```

### ChainStep
```go
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
```

## Variable System

Variables can be used in URLs, headers, and request bodies using the `{{variable_name}}` syntax.

### Setting Variables
```go
// Set global variables
cm.SetGlobalVariable("base_url", "https://api.example.com")
cm.SetGlobalVariable("api_version", "v1")

// Set chain variables
chain.Variables["user_id"] = "12345"
```

### Using Variables
```go
step := ChainStep{
    URL: "{{base_url}}/{{api_version}}/users/{{user_id}}",
    Headers: map[string]string{
        "Authorization": "Bearer {{auth_token}}",
    },
    Body: `{"user_id": "{{user_id}}", "action": "update"}`,
}
```

## Data Extraction

Extract data from responses using extract rules:

```go
extractRules := []ExtractRule{
    {
        Name:     "user_id",
        Source:   "body",        // "header", "body", or "status"
        Pattern:  "id",          // JSON path or header name
        Required: true,
    },
    {
        Name:     "auth_token",
        Source:   "header",
        Pattern:  "Authorization",
        Required: true,
    },
    {
        Name:     "status",
        Source:   "status",
        Required: true,
    },
}
```

## Conditional Logic

Control chain execution based on response conditions:

```go
conditions := []Condition{
    {
        Field:    "status_code",
        Operator: "equals",
        Value:    "200",
        Action:   "continue",
    },
    {
        Field:    "status_code",
        Operator: "not_equals",
        Value:    "200",
        Action:   "stop",
    },
}
```

### Supported Operators
- `equals`: Value equals expected
- `not_equals`: Value does not equal expected
- `contains`: Value contains expected string
- `not_contains`: Value does not contain expected string
- `greater_than`: Value is greater than expected
- `less_than`: Value is less than expected

### Supported Fields
- `status_code`: HTTP status code
- `success`: Request success status
- `duration`: Request duration

### Actions
- `continue`: Continue to next step
- `skip`: Skip remaining steps
- `stop`: Stop chain execution

## Interactive CLI

Use the interactive CLI for easy chain management:

```go
package main

import "APITestProgram/chain"

func main() {
    cli := chain.NewChainCLI()
    cli.Run()
}
```

### CLI Commands
1. **Create Chain**: Create a new request chain
2. **List Chains**: View all available chains
3. **View Chain**: Display chain details and steps
4. **Add Step**: Add a new step to a chain
5. **Update Step**: Modify an existing step
6. **Remove Step**: Delete a step from a chain
7. **Execute Chain**: Run a chain and view results
8. **Export Chain**: Save chain to JSON file
9. **Import Chain**: Load chain from JSON file
10. **Set Global Variable**: Set a global variable
11. **List Global Variables**: View all global variables
12. **Clear Global Variables**: Remove all global variables

## Export/Import

### Export Chain
```go
err := cm.ExportChain(chainID, "my_chain.json")
if err != nil {
    fmt.Printf("Export failed: %v\n", err)
}
```

### Import Chain
```go
err := cm.ImportChain("my_chain.json")
if err != nil {
    fmt.Printf("Import failed: %v\n", err)
}
```

## Error Handling

The chain system provides comprehensive error handling:

- **Step Errors**: Individual step failures are captured in results
- **Validation Errors**: Invalid configurations are caught early
- **Network Errors**: Connection and timeout issues are handled gracefully
- **Extraction Errors**: Failed data extraction is logged with details

## Best Practices

1. **Use Descriptive Names**: Give chains and steps meaningful names
2. **Set Appropriate Timeouts**: Configure timeouts based on expected response times
3. **Handle Errors Gracefully**: Use conditions to handle different response scenarios
4. **Extract Key Data**: Use extract rules to capture important response data
5. **Use Variables**: Leverage variables for dynamic URLs and data
6. **Test Incrementally**: Test individual steps before running full chains
7. **Document Chains**: Add descriptions to explain chain purpose and flow

## Example Use Cases

### API Testing Workflow
```go
// 1. Register user
// 2. Login with credentials
// 3. Get user profile
// 4. Update profile
// 5. Logout
```

### Data Migration
```go
// 1. Export data from old API
// 2. Transform data format
// 3. Import data to new API
// 4. Verify data integrity
```

### Integration Testing
```go
// 1. Create test data
// 2. Execute business logic
// 3. Verify results
// 4. Clean up test data
```

## Advanced Features

### Custom Extractors
You can extend the extraction system by implementing custom extractors for specific data formats.

### Chain Composition
Chains can be composed by calling other chains as steps, enabling complex workflows.

### Parallel Execution
Future versions will support parallel step execution for improved performance.

### Monitoring and Logging
Comprehensive logging and monitoring capabilities for production use.

## Troubleshooting

### Common Issues

1. **Variable Not Found**: Ensure variables are set before use
2. **Extraction Failed**: Check extract rule patterns and sources
3. **Condition Not Met**: Verify condition operators and values
4. **Timeout Errors**: Increase timeout values for slow endpoints
5. **Network Errors**: Check network connectivity and endpoint availability

### Debug Tips

1. **Enable Verbose Logging**: Use detailed logging to trace execution
2. **Test Individual Steps**: Run steps in isolation to identify issues
3. **Check Variable Values**: Verify variable values at each step
4. **Validate Conditions**: Test conditions with known response data
5. **Review Extract Rules**: Ensure extract patterns match response format
