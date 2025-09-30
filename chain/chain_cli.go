package chain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ChainCLI struct {
	manager *ChainManager
}

func NewChainCLI() *ChainCLI {
	return &ChainCLI{
		manager: NewChainManager(),
	}
}

func (cli *ChainCLI) Run() {
	for {
		fmt.Println("\n=== Request Chain Manager ===")
		fmt.Println("1. Create Chain")
		fmt.Println("2. List Chains")
		fmt.Println("3. View Chain")
		fmt.Println("4. Add Step")
		fmt.Println("5. Update Step")
		fmt.Println("6. Remove Step")
		fmt.Println("7. Execute Chain")
		fmt.Println("8. Export Chain")
		fmt.Println("9. Import Chain")
		fmt.Println("10. Set Global Variable")
		fmt.Println("11. List Global Variables")
		fmt.Println("12. Clear Global Variables")
		fmt.Println("0. Exit")
		
		choice := cli.getInput("Select option: ")
		
		switch choice {
		case "1":
			cli.createChain()
		case "2":
			cli.listChains()
		case "3":
			cli.viewChain()
		case "4":
			cli.addStep()
		case "5":
			cli.updateStep()
		case "6":
			cli.removeStep()
		case "7":
			cli.executeChain()
		case "8":
			cli.exportChain()
		case "9":
			cli.importChain()
		case "10":
			cli.setGlobalVariable()
		case "11":
			cli.listGlobalVariables()
		case "12":
			cli.clearGlobalVariables()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func (cli *ChainCLI) createChain() {
	name := cli.getInput("Chain name: ")
	description := cli.getInput("Chain description: ")
	
	chain := cli.manager.CreateChain(name, description)
	fmt.Printf("✅ Chain created with ID: %s\n", chain.ID)
}

func (cli *ChainCLI) listChains() {
	chains := cli.manager.ListChains()
	
	if len(chains) == 0 {
		fmt.Println("No chains found.")
		return
	}
	
	fmt.Println("\n📋 Available Chains:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-20s %-30s %-20s %-10s\n", "ID", "Name", "Description", "Steps")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, chain := range chains {
		fmt.Printf("%-20s %-30s %-20s %-10d\n", 
			chain.ID[:8], 
			chain.Name, 
			chain.Description, 
			len(chain.Steps))
	}
}

func (cli *ChainCLI) viewChain() {
	chainID := cli.getInput("Chain ID: ")
	
	chain, err := cli.manager.GetChain(chainID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("\n🔗 Chain: %s\n", chain.Name)
	fmt.Printf("Description: %s\n", chain.Description)
	fmt.Printf("Created: %s\n", chain.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", chain.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Steps: %d\n", len(chain.Steps))
	
	if len(chain.Steps) > 0 {
		fmt.Println("\n📝 Steps:")
		for i, step := range chain.Steps {
			status := "✅"
			if !step.Enabled {
				status = "❌"
			}
			fmt.Printf("  %d. %s %s (%s %s)\n", 
				i+1, status, step.Name, step.Method, step.URL)
		}
	}
}

func (cli *ChainCLI) addStep() {
	chainID := cli.getInput("Chain ID: ")
	
	chain, err := cli.manager.GetChain(chainID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("Adding step to chain: %s\n", chain.Name)
	
	step := ChainStep{
		Name:   cli.getInput("Step name: "),
		Method: cli.getInput("HTTP method (GET/POST/PUT/DELETE): "),
		URL:    cli.getInput("URL: "),
		Headers: make(map[string]string),
		Body:    cli.getInput("Request body (optional): "),
		Timeout: cli.getDuration("Timeout (e.g., 30s): "),
		Retries: cli.getInt("Retries: "),
		Delay:   cli.getDuration("Delay between retries (e.g., 1s): "),
		Wait:    cli.getDuration("Wait after step (e.g., 2s): "),
		Enabled: true,
	}
	
	// Add headers
	fmt.Println("Add headers (press Enter with empty key to finish):")
	for {
		key := cli.getInput("Header key: ")
		if key == "" {
			break
		}
		value := cli.getInput("Header value: ")
		step.Headers[key] = value
	}
	
	// Add extract rules
	fmt.Println("Add extract rules (press Enter with empty name to finish):")
	for {
		name := cli.getInput("Extract rule name: ")
		if name == "" {
			break
		}
		
		source := cli.getInput("Source (header/body/status): ")
		pattern := cli.getInput("Pattern (for header source, use header name): ")
		required := cli.getBool("Required (y/n): ")
		
		step.Extract = append(step.Extract, ExtractRule{
			Name:     name,
			Source:   source,
			Pattern:  pattern,
			Required: required,
		})
	}
	
	// Add conditions
	fmt.Println("Add conditions (press Enter with empty field to finish):")
	for {
		field := cli.getInput("Condition field (status_code/success/duration): ")
		if field == "" {
			break
		}
		
		operator := cli.getInput("Operator (equals/not_equals/contains/not_contains/greater_than/less_than): ")
		value := cli.getInput("Expected value: ")
		action := cli.getInput("Action (continue/skip/stop): ")
		
		step.Conditions = append(step.Conditions, Condition{
			Field:    field,
			Operator: operator,
			Value:    value,
			Action:   action,
		})
	}
	
	err = cli.manager.AddStep(chainID, step)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Println("✅ Step added successfully!")
}

func (cli *ChainCLI) updateStep() {
	chainID := cli.getInput("Chain ID: ")
	stepID := cli.getInput("Step ID: ")
	
	updates := make(map[string]interface{})
	
	fmt.Println("Enter new values (press Enter to skip):")
	
	if name := cli.getInput("New name: "); name != "" {
		updates["name"] = name
	}
	
	if method := cli.getInput("New method: "); method != "" {
		updates["method"] = method
	}
	
	if url := cli.getInput("New URL: "); url != "" {
		updates["url"] = url
	}
	
	if body := cli.getInput("New body: "); body != "" {
		updates["body"] = body
	}
	
	if enabled := cli.getInput("Enabled (y/n): "); enabled != "" {
		updates["enabled"] = enabled == "y" || enabled == "Y"
	}
	
	err := cli.manager.UpdateStep(chainID, stepID, updates)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Println("✅ Step updated successfully!")
}

func (cli *ChainCLI) removeStep() {
	chainID := cli.getInput("Chain ID: ")
	stepID := cli.getInput("Step ID: ")
	
	err := cli.manager.RemoveStep(chainID, stepID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Println("✅ Step removed successfully!")
}

func (cli *ChainCLI) executeChain() {
	chainID := cli.getInput("Chain ID: ")
	
	fmt.Println("🚀 Executing chain...")
	
	results, err := cli.manager.ExecuteChain(chainID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("\n📊 Execution Results (%d steps):\n", len(results))
	fmt.Println(strings.Repeat("=", 80))
	
	for i, result := range results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		
		fmt.Printf("\n%d. %s Step: %s\n", i+1, status, result.StepID)
		fmt.Printf("   Status: %d | Duration: %v\n", result.StatusCode, result.Duration)
		
		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}
		
		if len(result.Extracted) > 0 {
			fmt.Println("   Extracted Variables:")
			for k, v := range result.Extracted {
				fmt.Printf("     %s = %s\n", k, v)
			}
		}
	}
}

func (cli *ChainCLI) exportChain() {
	chainID := cli.getInput("Chain ID: ")
	filename := cli.getInput("Export filename: ")
	
	err := cli.manager.ExportChain(chainID, filename)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Chain exported to %s\n", filename)
}

func (cli *ChainCLI) importChain() {
	filename := cli.getInput("Import filename: ")
	
	err := cli.manager.ImportChain(filename)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Chain imported from %s\n", filename)
}

func (cli *ChainCLI) setGlobalVariable() {
	key := cli.getInput("Variable key: ")
	value := cli.getInput("Variable value: ")
	
	cli.manager.SetGlobalVariable(key, value)
	fmt.Printf("✅ Global variable set: %s = %s\n", key, value)
}

func (cli *ChainCLI) listGlobalVariables() {
	variables := cli.manager.variables
	
	if len(variables) == 0 {
		fmt.Println("No global variables set.")
		return
	}
	
	fmt.Println("\n🌍 Global Variables:")
	for k, v := range variables {
		fmt.Printf("  %s = %s\n", k, v)
	}
}

func (cli *ChainCLI) clearGlobalVariables() {
	cli.manager.ClearGlobalVariables()
	fmt.Println("✅ Global variables cleared.")
}

func (cli *ChainCLI) getInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func (cli *ChainCLI) getInt(prompt string) int {
	for {
		input := cli.getInput(prompt)
		if value, err := strconv.Atoi(input); err == nil {
			return value
		}
		fmt.Println("Invalid number. Please try again.")
	}
}

func (cli *ChainCLI) getBool(prompt string) bool {
	for {
		input := cli.getInput(prompt)
		switch strings.ToLower(input) {
		case "y", "yes", "true", "1":
			return true
		case "n", "no", "false", "0":
			return false
		default:
			fmt.Println("Please enter y/n, yes/no, true/false, or 1/0.")
		}
	}
}

func (cli *ChainCLI) getDuration(prompt string) time.Duration {
	for {
		input := cli.getInput(prompt)
		if duration, err := time.ParseDuration(input); err == nil {
			return duration
		}
		fmt.Println("Invalid duration. Please use format like '30s', '1m', '2h'.")
	}
}
