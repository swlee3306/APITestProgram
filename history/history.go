package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type HistoryEntry struct {
	ID          string            `json:"id"`
	Timestamp   time.Time         `json:"timestamp"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body,omitempty"`
	StatusCode  int               `json:"status_code,omitempty"`
	Response    string            `json:"response,omitempty"`
	Duration    time.Duration     `json:"duration,omitempty"`
	Success     bool              `json:"success"`
	Tags        []string          `json:"tags,omitempty"`
	Notes       string            `json:"notes,omitempty"`
}

type HistoryManager struct {
	entries    []HistoryEntry
	historyFile string
	maxEntries  int
}

func NewHistoryManager(historyFile string, maxEntries int) *HistoryManager {
	return &HistoryManager{
		entries:    make([]HistoryEntry, 0),
		historyFile: historyFile,
		maxEntries:  maxEntries,
	}
}

func (hm *HistoryManager) Load() error {
	// Create history directory if it doesn't exist
	dir := filepath.Dir(hm.historyFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}
	
	// Load history file if it exists
	if _, err := os.Stat(hm.historyFile); os.IsNotExist(err) {
		return nil
	}
	
	data, err := os.ReadFile(hm.historyFile)
	if err != nil {
		return fmt.Errorf("failed to read history file: %w", err)
	}
	
	if err := json.Unmarshal(data, &hm.entries); err != nil {
		return fmt.Errorf("failed to parse history file: %w", err)
	}
	
	// Sort by timestamp (newest first)
	sort.Slice(hm.entries, func(i, j int) bool {
		return hm.entries[i].Timestamp.After(hm.entries[j].Timestamp)
	})
	
	return nil
}

func (hm *HistoryManager) Save() error {
	data, err := json.MarshalIndent(hm.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}
	
	if err := os.WriteFile(hm.historyFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}
	
	return nil
}

func (hm *HistoryManager) AddEntry(entry HistoryEntry) {
	// Generate ID if not provided
	if entry.ID == "" {
		entry.ID = hm.generateID()
	}
	
	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	// Add to entries
	hm.entries = append(hm.entries, entry)
	
	// Sort by timestamp (newest first)
	sort.Slice(hm.entries, func(i, j int) bool {
		return hm.entries[i].Timestamp.After(hm.entries[j].Timestamp)
	})
	
	// Limit entries
	if len(hm.entries) > hm.maxEntries {
		hm.entries = hm.entries[:hm.maxEntries]
	}
}

func (hm *HistoryManager) GetEntry(id string) (*HistoryEntry, error) {
	for _, entry := range hm.entries {
		if entry.ID == id {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("entry not found: %s", id)
}

func (hm *HistoryManager) GetEntries(limit int) []HistoryEntry {
	if limit <= 0 || limit > len(hm.entries) {
		return hm.entries
	}
	return hm.entries[:limit]
}

func (hm *HistoryManager) SearchEntries(query string) []HistoryEntry {
	var results []HistoryEntry
	query = strings.ToLower(query)
	
	for _, entry := range hm.entries {
		if strings.Contains(strings.ToLower(entry.URL), query) ||
		   strings.Contains(strings.ToLower(entry.Method), query) ||
		   strings.Contains(strings.ToLower(entry.Notes), query) {
			results = append(results, entry)
		}
	}
	
	return results
}

func (hm *HistoryManager) FilterByMethod(method string) []HistoryEntry {
	var results []HistoryEntry
	method = strings.ToUpper(method)
	
	for _, entry := range hm.entries {
		if strings.ToUpper(entry.Method) == method {
			results = append(results, entry)
		}
	}
	
	return results
}

func (hm *HistoryManager) FilterByStatus(success bool) []HistoryEntry {
	var results []HistoryEntry
	
	for _, entry := range hm.entries {
		if entry.Success == success {
			results = append(results, entry)
		}
	}
	
	return results
}

func (hm *HistoryManager) FilterByTag(tag string) []HistoryEntry {
	var results []HistoryEntry
	tag = strings.ToLower(tag)
	
	for _, entry := range hm.entries {
		for _, entryTag := range entry.Tags {
			if strings.ToLower(entryTag) == tag {
				results = append(results, entry)
				break
			}
		}
	}
	
	return results
}

func (hm *HistoryManager) FilterByDateRange(start, end time.Time) []HistoryEntry {
	var results []HistoryEntry
	
	for _, entry := range hm.entries {
		if entry.Timestamp.After(start) && entry.Timestamp.Before(end) {
			results = append(results, entry)
		}
	}
	
	return results
}

func (hm *HistoryManager) UpdateEntry(id string, updates map[string]interface{}) error {
	for i, entry := range hm.entries {
		if entry.ID == id {
			// Update fields
			for key, value := range updates {
				switch key {
				case "tags":
					if tags, ok := value.([]string); ok {
						hm.entries[i].Tags = tags
					}
				case "notes":
					if notes, ok := value.(string); ok {
						hm.entries[i].Notes = notes
					}
				case "success":
					if success, ok := value.(bool); ok {
						hm.entries[i].Success = success
					}
				}
			}
			return nil
		}
	}
	return fmt.Errorf("entry not found: %s", id)
}

func (hm *HistoryManager) DeleteEntry(id string) error {
	for i, entry := range hm.entries {
		if entry.ID == id {
			hm.entries = append(hm.entries[:i], hm.entries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("entry not found: %s", id)
}

func (hm *HistoryManager) Clear() {
	hm.entries = make([]HistoryEntry, 0)
}

func (hm *HistoryManager) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	total := len(hm.entries)
	successful := 0
	failed := 0
	
	methodCounts := make(map[string]int)
	statusCounts := make(map[int]int)
	
	for _, entry := range hm.entries {
		if entry.Success {
			successful++
		} else {
			failed++
		}
		
		methodCounts[entry.Method]++
		statusCounts[entry.StatusCode]++
	}
	
	stats["total"] = total
	stats["successful"] = successful
	stats["failed"] = failed
	stats["success_rate"] = float64(successful) / float64(total) * 100
	stats["method_counts"] = methodCounts
	stats["status_counts"] = statusCounts
	
	return stats
}

func (hm *HistoryManager) Export(filename string) error {
	data, err := json.MarshalIndent(hm.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}
	
	return nil
}

func (hm *HistoryManager) Import(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}
	
	var importedEntries []HistoryEntry
	if err := json.Unmarshal(data, &importedEntries); err != nil {
		return fmt.Errorf("failed to parse import file: %w", err)
	}
	
	// Merge with existing entries
	hm.entries = append(hm.entries, importedEntries...)
	
	// Sort by timestamp
	sort.Slice(hm.entries, func(i, j int) bool {
		return hm.entries[i].Timestamp.After(hm.entries[j].Timestamp)
	})
	
	// Limit entries
	if len(hm.entries) > hm.maxEntries {
		hm.entries = hm.entries[:hm.maxEntries]
	}
	
	return nil
}

func (hm *HistoryManager) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (hm *HistoryManager) GetRecentEntries(count int) []HistoryEntry {
	if count <= 0 || count > len(hm.entries) {
		return hm.entries
	}
	return hm.entries[:count]
}

func (hm *HistoryManager) GetFrequentURLs(limit int) []string {
	urlCounts := make(map[string]int)
	
	for _, entry := range hm.entries {
		urlCounts[entry.URL]++
	}
	
	// Sort by count
	type urlCount struct {
		URL   string
		Count int
	}
	
	var sorted []urlCount
	for url, count := range urlCounts {
		sorted = append(sorted, urlCount{url, count})
	}
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})
	
	var result []string
	for i, item := range sorted {
		if i >= limit {
			break
		}
		result = append(result, item.URL)
	}
	
	return result
}

func (hm *HistoryManager) GetFrequentMethods(limit int) []string {
	methodCounts := make(map[string]int)
	
	for _, entry := range hm.entries {
		methodCounts[entry.Method]++
	}
	
	// Sort by count
	type methodCount struct {
		Method string
		Count  int
	}
	
	var sorted []methodCount
	for method, count := range methodCounts {
		sorted = append(sorted, methodCount{method, count})
	}
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})
	
	var result []string
	for i, item := range sorted {
		if i >= limit {
			break
		}
		result = append(result, item.Method)
	}
	
	return result
}
