package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type InteractiveCLI struct {
	scanner *bufio.Scanner
	history []string
	config  *CLIConfig
}

type CLIConfig struct {
	Prompt      string
	HistoryFile string
	MaxHistory  int
	AutoComplete bool
}

type MenuOption struct {
	Key         string
	Description string
	Action      func() error
}

func NewInteractiveCLI() *InteractiveCLI {
	return &InteractiveCLI{
		scanner: bufio.NewScanner(os.Stdin),
		history: make([]string, 0),
		config: &CLIConfig{
			Prompt:      "APITest> ",
			HistoryFile: ".api_test_history",
			MaxHistory:  100,
			AutoComplete: true,
		},
	}
}

func (cli *InteractiveCLI) Start() error {
	fmt.Println("🚀 API Test Program - Interactive Mode")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")
	
	// Load history
	cli.loadHistory()
	
	for {
		// Display prompt
		fmt.Print(cli.config.Prompt)
		
		// Read input
		if !cli.scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(cli.scanner.Text())
		
		// Handle empty input
		if input == "" {
			continue
		}
		
		// Add to history
		cli.addToHistory(input)
		
		// Handle commands
		if err := cli.handleCommand(input); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
	}
	
	// Save history
	cli.saveHistory()
	
	return nil
}

func (cli *InteractiveCLI) handleCommand(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	
	command := strings.ToLower(parts[0])
	args := parts[1:]
	
	switch command {
	case "help", "h":
		cli.showHelp()
	case "exit", "quit", "q":
		fmt.Println("👋 Goodbye!")
		os.Exit(0)
	case "clear", "cls":
		cli.clearScreen()
	case "history":
		cli.showHistory()
	case "config":
		cli.showConfig()
	case "set":
		return cli.handleSetCommand(args)
	case "get":
		return cli.handleGetCommand(args)
	case "test":
		return cli.handleTestCommand(args)
	case "save":
		return cli.handleSaveCommand(args)
	case "load":
		return cli.handleLoadCommand(args)
	case "list":
		cli.showAvailableCommands()
	default:
		fmt.Printf("❓ Unknown command: %s\n", command)
		fmt.Println("Type 'help' for available commands")
	}
	
	return nil
}

func (cli *InteractiveCLI) showHelp() {
	fmt.Println("\n📚 Available Commands:")
	fmt.Println("  help, h           - Show this help message")
	fmt.Println("  exit, quit, q     - Exit the program")
	fmt.Println("  clear, cls        - Clear the screen")
	fmt.Println("  history           - Show command history")
	fmt.Println("  config            - Show current configuration")
	fmt.Println("  set <key> <value> - Set configuration value")
	fmt.Println("  get <key>         - Get configuration value")
	fmt.Println("  test <url>        - Test API endpoint")
	fmt.Println("  save <filename>   - Save current session")
	fmt.Println("  load <filename>   - Load saved session")
	fmt.Println("  list              - List available commands")
	fmt.Println()
}

func (cli *InteractiveCLI) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func (cli *InteractiveCLI) showHistory() {
	fmt.Println("\n📜 Command History:")
	for i, cmd := range cli.history {
		fmt.Printf("  %d: %s\n", i+1, cmd)
	}
	fmt.Println()
}

func (cli *InteractiveCLI) showConfig() {
	fmt.Println("\n⚙️  Current Configuration:")
	fmt.Printf("  Prompt: %s\n", cli.config.Prompt)
	fmt.Printf("  History File: %s\n", cli.config.HistoryFile)
	fmt.Printf("  Max History: %d\n", cli.config.MaxHistory)
	fmt.Printf("  Auto Complete: %t\n", cli.config.AutoComplete)
	fmt.Println()
}

func (cli *InteractiveCLI) handleSetCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: set <key> <value>")
	}
	
	key := strings.ToLower(args[0])
	value := strings.Join(args[1:], " ")
	
	switch key {
	case "prompt":
		cli.config.Prompt = value + " "
		fmt.Printf("✅ Prompt set to: %s\n", cli.config.Prompt)
	case "maxhistory":
		if max, err := strconv.Atoi(value); err == nil {
			cli.config.MaxHistory = max
			fmt.Printf("✅ Max History set to: %d\n", max)
		} else {
			return fmt.Errorf("invalid number: %s", value)
		}
	case "autocomplete":
		if value == "true" || value == "1" {
			cli.config.AutoComplete = true
		} else if value == "false" || value == "0" {
			cli.config.AutoComplete = false
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
		fmt.Printf("✅ Auto Complete set to: %t\n", cli.config.AutoComplete)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	
	return nil
}

func (cli *InteractiveCLI) handleGetCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: get <key>")
	}
	
	key := strings.ToLower(args[0])
	
	switch key {
	case "prompt":
		fmt.Printf("Prompt: %s\n", cli.config.Prompt)
	case "maxhistory":
		fmt.Printf("Max History: %d\n", cli.config.MaxHistory)
	case "autocomplete":
		fmt.Printf("Auto Complete: %t\n", cli.config.AutoComplete)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	
	return nil
}

func (cli *InteractiveCLI) handleTestCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: test <url>")
	}
	
	url := args[0]
	fmt.Printf("🧪 Testing API endpoint: %s\n", url)
	
	// Simulate API test
	time.Sleep(1 * time.Second)
	fmt.Printf("✅ API test completed for: %s\n", url)
	
	return nil
}

func (cli *InteractiveCLI) handleSaveCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: save <filename>")
	}
	
	filename := args[0]
	fmt.Printf("💾 Saving session to: %s\n", filename)
	
	// Simulate save
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("✅ Session saved to: %s\n", filename)
	
	return nil
}

func (cli *InteractiveCLI) handleLoadCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: load <filename>")
	}
	
	filename := args[0]
	fmt.Printf("📂 Loading session from: %s\n", filename)
	
	// Simulate load
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("✅ Session loaded from: %s\n", filename)
	
	return nil
}

func (cli *InteractiveCLI) showAvailableCommands() {
	fmt.Println("\n🔧 Available Commands:")
	commands := []string{
		"help", "exit", "clear", "history", "config",
		"set", "get", "test", "save", "load", "list",
	}
	
	for i, cmd := range commands {
		if i%4 == 0 {
			fmt.Println()
		}
		fmt.Printf("  %-10s", cmd)
	}
	fmt.Println("\n")
}

func (cli *InteractiveCLI) addToHistory(command string) {
	cli.history = append(cli.history, command)
	
	// Limit history size
	if len(cli.history) > cli.config.MaxHistory {
		cli.history = cli.history[1:]
	}
}

func (cli *InteractiveCLI) loadHistory() {
	file, err := os.Open(cli.config.HistoryFile)
	if err != nil {
		return
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cli.history = append(cli.history, scanner.Text())
	}
}

func (cli *InteractiveCLI) saveHistory() {
	file, err := os.Create(cli.config.HistoryFile)
	if err != nil {
		return
	}
	defer file.Close()
	
	for _, cmd := range cli.history {
		fmt.Fprintln(file, cmd)
	}
}

func (cli *InteractiveCLI) ShowMenu(title string, options []MenuOption) (string, error) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))
	
	for _, option := range options {
		fmt.Printf("  %s - %s\n", option.Key, option.Description)
	}
	
	fmt.Print("\nEnter your choice: ")
	if !cli.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}
	
	choice := strings.TrimSpace(cli.scanner.Text())
	
	// Find and execute the chosen option
	for _, option := range options {
		if strings.ToLower(option.Key) == strings.ToLower(choice) {
			if option.Action != nil {
				if err := option.Action(); err != nil {
					return choice, err
				}
			}
			return choice, nil
		}
	}
	
	return choice, fmt.Errorf("invalid choice: %s", choice)
}

func (cli *InteractiveCLI) Confirm(message string) (bool, error) {
	fmt.Printf("%s (y/N): ", message)
	if !cli.scanner.Scan() {
		return false, fmt.Errorf("failed to read input")
	}
	
	response := strings.ToLower(strings.TrimSpace(cli.scanner.Text()))
	return response == "y" || response == "yes", nil
}

func (cli *InteractiveCLI) GetInput(prompt string) (string, error) {
	fmt.Print(prompt)
	if !cli.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}
	
	return strings.TrimSpace(cli.scanner.Text()), nil
}
