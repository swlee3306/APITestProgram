package logging

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type SecureLogger struct {
	sensitiveFields []string
	maskPatterns    []*regexp.Regexp
	enabled         bool
}

type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Request     *RequestLog            `json:"request,omitempty"`
	Response    *ResponseLog           `json:"response,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

type RequestLog struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

type ResponseLog struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body,omitempty"`
	Size       int64             `json:"size"`
}

func NewSecureLogger() *SecureLogger {
	sl := &SecureLogger{
		sensitiveFields: []string{
			"password", "passwd", "pwd",
			"token", "access_token", "refresh_token",
			"secret", "key", "api_key",
			"authorization", "auth",
			"credit_card", "card_number",
			"ssn", "social_security",
			"email", "phone", "mobile",
		},
		enabled: true,
	}
	
	// Compile regex patterns for sensitive data
	sl.compileMaskPatterns()
	
	return sl
}

func (sl *SecureLogger) compileMaskPatterns() {
	patterns := []string{
		`(?i)(password|passwd|pwd)\s*[:=]\s*["']?([^"'\s]+)["']?`,
		`(?i)(token|access_token|refresh_token)\s*[:=]\s*["']?([^"'\s]+)["']?`,
		`(?i)(secret|key|api_key)\s*[:=]\s*["']?([^"'\s]+)["']?`,
		`(?i)(authorization|auth)\s*[:=]\s*["']?([^"'\s]+)["']?`,
		`(?i)(credit_card|card_number)\s*[:=]\s*["']?([^"'\s]+)["']?`,
		`(?i)(ssn|social_security)\s*[:=]\s*["']?([^"'\s]+)["']?`,
		`(?i)(email|phone|mobile)\s*[:=]\s*["']?([^"'\s]+)["']?`,
	}
	
	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern)
		if err == nil {
			sl.maskPatterns = append(sl.maskPatterns, regex)
		}
	}
}

func (sl *SecureLogger) AddSensitiveField(field string) {
	sl.sensitiveFields = append(sl.sensitiveFields, strings.ToLower(field))
}

func (sl *SecureLogger) RemoveSensitiveField(field string) {
	for i, f := range sl.sensitiveFields {
		if strings.ToLower(f) == strings.ToLower(field) {
			sl.sensitiveFields = append(sl.sensitiveFields[:i], sl.sensitiveFields[i+1:]...)
			break
		}
	}
}

func (sl *SecureLogger) Enable() {
	sl.enabled = true
}

func (sl *SecureLogger) Disable() {
	sl.enabled = false
}

func (sl *SecureLogger) IsEnabled() bool {
	return sl.enabled
}

func (sl *SecureLogger) LogRequest(req *http.Request, body []byte) *LogEntry {
	if !sl.enabled {
		return nil
	}
	
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "HTTP Request",
		Request: &RequestLog{
			Method:  req.Method,
			URL:     req.URL.String(),
			Headers: sl.sanitizeHeaders(req.Header),
			Body:    sl.sanitizeBody(string(body)),
		},
		UserAgent: req.UserAgent(),
		IPAddress: getClientIP(req),
	}
	
	return entry
}

func (sl *SecureLogger) LogResponse(resp *http.Response, body []byte) *LogEntry {
	if !sl.enabled {
		return nil
	}
	
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "HTTP Response",
		Response: &ResponseLog{
			StatusCode: resp.StatusCode,
			Headers:    sl.sanitizeHeaders(resp.Header),
			Body:      sl.sanitizeBody(string(body)),
			Size:      int64(len(body)),
		},
	}
	
	return entry
}

func (sl *SecureLogger) LogError(err error, req *http.Request) *LogEntry {
	if !sl.enabled {
		return nil
	}
	
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Message:   fmt.Sprintf("Request failed: %v", err),
		Request: &RequestLog{
			Method:  req.Method,
			URL:     req.URL.String(),
			Headers: sl.sanitizeHeaders(req.Header),
		},
		UserAgent: req.UserAgent(),
		IPAddress: getClientIP(req),
	}
	
	return entry
}

func (sl *SecureLogger) sanitizeHeaders(headers http.Header) map[string]string {
	sanitized := make(map[string]string)
	
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		
		// Check if this is a sensitive header
		if sl.isSensitiveField(lowerKey) {
			sanitized[key] = sl.maskValue(strings.Join(values, ", "))
		} else {
			sanitized[key] = strings.Join(values, ", ")
		}
	}
	
	return sanitized
}

func (sl *SecureLogger) sanitizeBody(body string) string {
	if body == "" {
		return body
	}
	
	// Apply regex patterns to mask sensitive data
	sanitized := body
	for _, pattern := range sl.maskPatterns {
		sanitized = pattern.ReplaceAllStringFunc(sanitized, func(match string) string {
			// Extract the field name and value
			parts := strings.Split(match, ":")
			if len(parts) >= 2 {
				field := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(strings.Join(parts[1:], ":"))
				return fmt.Sprintf("%s: %s", field, sl.maskValue(value))
			}
			return match
		})
	}
	
	return sanitized
}

func (sl *SecureLogger) isSensitiveField(field string) bool {
	for _, sensitive := range sl.sensitiveFields {
		if strings.Contains(strings.ToLower(field), strings.ToLower(sensitive)) {
			return true
		}
	}
	return false
}

func (sl *SecureLogger) maskValue(value string) string {
	if value == "" {
		return value
	}
	
	// Remove quotes and whitespace
	value = strings.Trim(value, `"' `)
	
	// If value is too short, mask it completely
	if len(value) <= 4 {
		return "***"
	}
	
	// For longer values, show first 2 and last 2 characters
	if len(value) <= 8 {
		return value[:2] + "***" + value[len(value)-2:]
	}
	
	// For very long values, show first 4 and last 4 characters
	return value[:4] + "***" + value[len(value)-4:]
}

func (sl *SecureLogger) hashSensitiveValue(value string) string {
	if value == "" {
		return value
	}
	
	hash := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(hash[:8])
}

func (sl *SecureLogger) LogWithFields(level, message string, fields map[string]interface{}) *LogEntry {
	if !sl.enabled {
		return nil
	}
	
	// Sanitize fields
	sanitizedFields := make(map[string]interface{})
	for key, value := range fields {
		if sl.isSensitiveField(key) {
			sanitizedFields[key] = sl.maskValue(fmt.Sprintf("%v", value))
		} else {
			sanitizedFields[key] = value
		}
	}
	
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    sanitizedFields,
	}
	
	return entry
}

func getClientIP(req *http.Request) string {
	// Check X-Forwarded-For header
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Check X-Real-IP header
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	ip := req.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

func (sl *SecureLogger) SetSensitiveFields(fields []string) {
	sl.sensitiveFields = fields
	sl.compileMaskPatterns()
}

func (sl *SecureLogger) GetSensitiveFields() []string {
	return sl.sensitiveFields
}
