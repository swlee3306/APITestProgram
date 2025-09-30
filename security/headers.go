package security

import (
	"fmt"
	"net/http"
	"strings"
)

type SecurityHeaders struct {
	XContentTypeOptions   string
	XFrameOptions        string
	XSSProtection        string
	StrictTransportSecurity string
	ContentSecurityPolicy string
	ReferrerPolicy       string
	PermissionsPolicy    string
}

type SecurityHeaderManager struct {
	headers SecurityHeaders
	enabled bool
}

func NewSecurityHeaderManager() *SecurityHeaderManager {
	return &SecurityHeaderManager{
		headers: SecurityHeaders{
			XContentTypeOptions:   "nosniff",
			XFrameOptions:        "DENY",
			XSSProtection:        "1; mode=block",
			StrictTransportSecurity: "max-age=31536000; includeSubDomains",
			ContentSecurityPolicy: "default-src 'self'",
			ReferrerPolicy:       "strict-origin-when-cross-origin",
			PermissionsPolicy:    "geolocation=(), microphone=(), camera=()",
		},
		enabled: true,
	}
}

func (shm *SecurityHeaderManager) SetHeaders(headers SecurityHeaders) {
	shm.headers = headers
}

func (shm *SecurityHeaderManager) Enable() {
	shm.enabled = true
}

func (shm *SecurityHeaderManager) Disable() {
	shm.enabled = false
}

func (shm *SecurityHeaderManager) IsEnabled() bool {
	return shm.enabled
}

func (shm *SecurityHeaderManager) AddSecurityHeaders(req *http.Request) {
	if !shm.enabled {
		return
	}
	
	// Add security headers to the request
	req.Header.Set("X-Content-Type-Options", shm.headers.XContentTypeOptions)
	req.Header.Set("X-Frame-Options", shm.headers.XFrameOptions)
	req.Header.Set("X-XSS-Protection", shm.headers.XSSProtection)
	req.Header.Set("Strict-Transport-Security", shm.headers.StrictTransportSecurity)
	req.Header.Set("Content-Security-Policy", shm.headers.ContentSecurityPolicy)
	req.Header.Set("Referrer-Policy", shm.headers.ReferrerPolicy)
	req.Header.Set("Permissions-Policy", shm.headers.PermissionsPolicy)
}

func (shm *SecurityHeaderManager) AddSecurityHeadersToResponse(resp *http.Response) {
	if !shm.enabled {
		return
	}
	
	// Add security headers to the response
	resp.Header.Set("X-Content-Type-Options", shm.headers.XContentTypeOptions)
	resp.Header.Set("X-Frame-Options", shm.headers.XFrameOptions)
	resp.Header.Set("X-XSS-Protection", shm.headers.XSSProtection)
	resp.Header.Set("Strict-Transport-Security", shm.headers.StrictTransportSecurity)
	resp.Header.Set("Content-Security-Policy", shm.headers.ContentSecurityPolicy)
	resp.Header.Set("Referrer-Policy", shm.headers.ReferrerPolicy)
	resp.Header.Set("Permissions-Policy", shm.headers.PermissionsPolicy)
}

func (shm *SecurityHeaderManager) SetContentSecurityPolicy(policy string) {
	shm.headers.ContentSecurityPolicy = policy
}

func (shm *SecurityHeaderManager) SetStrictTransportSecurity(maxAge int, includeSubDomains bool) {
	policy := fmt.Sprintf("max-age=%d", maxAge)
	if includeSubDomains {
		policy += "; includeSubDomains"
	}
	shm.headers.StrictTransportSecurity = policy
}

func (shm *SecurityHeaderManager) SetXFrameOptions(option string) {
	validOptions := []string{"DENY", "SAMEORIGIN", "ALLOW-FROM"}
	if contains(validOptions, option) {
		shm.headers.XFrameOptions = option
	}
}

func (shm *SecurityHeaderManager) SetReferrerPolicy(policy string) {
	validPolicies := []string{
		"no-referrer",
		"no-referrer-when-downgrade",
		"origin",
		"origin-when-cross-origin",
		"same-origin",
		"strict-origin",
		"strict-origin-when-cross-origin",
		"unsafe-url",
	}
	if contains(validPolicies, policy) {
		shm.headers.ReferrerPolicy = policy
	}
}

func (shm *SecurityHeaderManager) GetHeaders() SecurityHeaders {
	return shm.headers
}

func (shm *SecurityHeaderManager) ValidateHeaders() []string {
	var warnings []string
	
	// Check for missing headers
	if shm.headers.XContentTypeOptions == "" {
		warnings = append(warnings, "X-Content-Type-Options header is missing")
	}
	
	if shm.headers.XFrameOptions == "" {
		warnings = append(warnings, "X-Frame-Options header is missing")
	}
	
	if shm.headers.XSSProtection == "" {
		warnings = append(warnings, "X-XSS-Protection header is missing")
	}
	
	// Check for weak security settings
	if shm.headers.XFrameOptions == "ALLOWALL" {
		warnings = append(warnings, "X-Frame-Options is set to ALLOWALL, which is insecure")
	}
	
	if strings.Contains(shm.headers.XSSProtection, "0") {
		warnings = append(warnings, "X-XSS-Protection is disabled")
	}
	
	return warnings
}

func (shm *SecurityHeaderManager) GetSecurityScore() int {
	score := 0
	
	// Base score for having security headers enabled
	if shm.enabled {
		score += 20
	}
	
	// Score for individual headers
	if shm.headers.XContentTypeOptions != "" {
		score += 10
	}
	if shm.headers.XFrameOptions != "" {
		score += 10
	}
	if shm.headers.XSSProtection != "" {
		score += 10
	}
	if shm.headers.StrictTransportSecurity != "" {
		score += 15
	}
	if shm.headers.ContentSecurityPolicy != "" {
		score += 15
	}
	if shm.headers.ReferrerPolicy != "" {
		score += 10
	}
	if shm.headers.PermissionsPolicy != "" {
		score += 10
	}
	
	return score
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Predefined security header configurations
func GetDevelopmentHeaders() SecurityHeaders {
	return SecurityHeaders{
		XContentTypeOptions:   "nosniff",
		XFrameOptions:        "SAMEORIGIN",
		XSSProtection:        "1; mode=block",
		StrictTransportSecurity: "", // Not needed for development
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval'",
		ReferrerPolicy:       "strict-origin-when-cross-origin",
		PermissionsPolicy:    "geolocation=(), microphone=(), camera=()",
	}
}

func GetProductionHeaders() SecurityHeaders {
	return SecurityHeaders{
		XContentTypeOptions:   "nosniff",
		XFrameOptions:        "DENY",
		XSSProtection:        "1; mode=block",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains; preload",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'",
		ReferrerPolicy:       "strict-origin-when-cross-origin",
		PermissionsPolicy:    "geolocation=(), microphone=(), camera=()",
	}
}
