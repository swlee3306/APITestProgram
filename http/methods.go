package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
	TRACE   HTTPMethod = "TRACE"
	CONNECT HTTPMethod = "CONNECT"
)

type HTTPClient struct {
	client      *http.Client
	baseURL     string
	headers     map[string]string
	timeout     time.Duration
	retries     int
	delay       time.Duration
	verifySSL   bool
	proxyURL    string
	followRedirects bool
}

type RequestOptions struct {
	Method      HTTPMethod
	URL         string
	Headers     map[string]string
	Body        string
	Timeout     time.Duration
	Retries     int
	Delay       time.Duration
	VerifySSL   bool
	ProxyURL    string
	FollowRedirects bool
}

type Response struct {
	StatusCode  int
	Headers     map[string]string
	Body        string
	Size        int64
	Duration    time.Duration
	Success     bool
	Error       error
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers:          make(map[string]string),
		timeout:          30 * time.Second,
		retries:          3,
		delay:            1 * time.Second,
		verifySSL:        true,
		followRedirects:  true,
	}
}

func (hc *HTTPClient) SetBaseURL(url string) {
	hc.baseURL = url
}

func (hc *HTTPClient) SetHeaders(headers map[string]string) {
	hc.headers = headers
}

func (hc *HTTPClient) AddHeader(key, value string) {
	if hc.headers == nil {
		hc.headers = make(map[string]string)
	}
	hc.headers[key] = value
}

func (hc *HTTPClient) SetTimeout(timeout time.Duration) {
	hc.timeout = timeout
	hc.client.Timeout = timeout
}

func (hc *HTTPClient) SetRetries(retries int) {
	hc.retries = retries
}

func (hc *HTTPClient) SetDelay(delay time.Duration) {
	hc.delay = delay
}

func (hc *HTTPClient) SetVerifySSL(verify bool) {
	hc.verifySSL = verify
}

func (hc *HTTPClient) SetProxyURL(url string) {
	hc.proxyURL = url
}

func (hc *HTTPClient) SetFollowRedirects(follow bool) {
	hc.followRedirects = follow
}

func (hc *HTTPClient) Request(options RequestOptions) (*Response, error) {
	// Build full URL
	url := options.URL
	if hc.baseURL != "" && !strings.HasPrefix(url, "http") {
		url = hc.baseURL + "/" + strings.TrimPrefix(url, "/")
	}
	
	// Create request
	req, err := hc.createRequest(options.Method, url, options.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	hc.setRequestHeaders(req, options.Headers)
	
	// Execute request with retries
	return hc.executeWithRetries(req, options)
}

func (hc *HTTPClient) Get(url string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  GET,
		URL:     url,
		Headers: headers,
	})
}

func (hc *HTTPClient) Post(url string, body string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  POST,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

func (hc *HTTPClient) Put(url string, body string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  PUT,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

func (hc *HTTPClient) Patch(url string, body string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  PATCH,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

func (hc *HTTPClient) Delete(url string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  DELETE,
		URL:     url,
		Headers: headers,
	})
}

func (hc *HTTPClient) Head(url string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  HEAD,
		URL:     url,
		Headers: headers,
	})
}

func (hc *HTTPClient) Options(url string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  OPTIONS,
		URL:     url,
		Headers: headers,
	})
}

func (hc *HTTPClient) Trace(url string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  TRACE,
		URL:     url,
		Headers: headers,
	})
}

func (hc *HTTPClient) Connect(url string, headers map[string]string) (*Response, error) {
	return hc.Request(RequestOptions{
		Method:  CONNECT,
		URL:     url,
		Headers: headers,
	})
}

func (hc *HTTPClient) createRequest(method HTTPMethod, url string, body string) (*http.Request, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	
	req, err := http.NewRequest(string(method), url, bodyReader)
	if err != nil {
		return nil, err
	}
	
	return req, nil
}

func (hc *HTTPClient) setRequestHeaders(req *http.Request, optionsHeaders map[string]string) {
	// Set default headers
	for key, value := range hc.headers {
		req.Header.Set(key, value)
	}
	
	// Set request-specific headers
	for key, value := range optionsHeaders {
		req.Header.Set(key, value)
	}
}

func (hc *HTTPClient) executeWithRetries(req *http.Request, options RequestOptions) (*Response, error) {
	var lastErr error
	retries := options.Retries
	if retries <= 0 {
		retries = hc.retries
	}
	
	for attempt := 0; attempt <= retries; attempt++ {
		start := time.Now()
		
		resp, err := hc.client.Do(req)
		duration := time.Since(start)
		
		if err != nil {
			lastErr = err
			if attempt < retries {
				time.Sleep(hc.delay)
				continue
			}
			return &Response{
				Duration: duration,
				Success:  false,
				Error:    err,
			}, err
		}
		
		// Read response body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			lastErr = err
			if attempt < retries {
				time.Sleep(hc.delay)
				continue
			}
			return &Response{
				StatusCode: resp.StatusCode,
				Duration:   duration,
				Success:   false,
				Error:     err,
			}, err
		}
		
		// Convert headers
		headers := make(map[string]string)
		for key, values := range resp.Header {
			headers[key] = strings.Join(values, ", ")
		}
		
		success := resp.StatusCode >= 200 && resp.StatusCode < 300
		
		return &Response{
			StatusCode: resp.StatusCode,
			Headers:    headers,
			Body:       string(body),
			Size:       int64(len(body)),
			Duration:   duration,
			Success:    success,
		}, nil
	}
	
	return nil, lastErr
}

func (hc *HTTPClient) BatchRequest(requests []RequestOptions) ([]*Response, error) {
	var responses []*Response
	
	for _, req := range requests {
		resp, err := hc.Request(req)
		if err != nil {
			responses = append(responses, &Response{
				Success: false,
				Error:   err,
			})
		} else {
			responses = append(responses, resp)
		}
	}
	
	return responses, nil
}

func (hc *HTTPClient) GetSupportedMethods() []HTTPMethod {
	return []HTTPMethod{
		GET, POST, PUT, PATCH, DELETE,
		HEAD, OPTIONS, TRACE, CONNECT,
	}
}

func (hc *HTTPClient) IsMethodSupported(method string) bool {
	supported := hc.GetSupportedMethods()
	for _, m := range supported {
		if strings.ToUpper(method) == string(m) {
			return true
		}
	}
	return false
}

func (hc *HTTPClient) GetMethodDescription(method HTTPMethod) string {
	descriptions := map[HTTPMethod]string{
		GET:     "Retrieve data from the server",
		POST:    "Submit data to the server",
		PUT:     "Update or create a resource",
		PATCH:   "Partially update a resource",
		DELETE:  "Remove a resource",
		HEAD:    "Get response headers only",
		OPTIONS: "Get allowed methods for a resource",
		TRACE:   "Echo the request for debugging",
		CONNECT: "Establish a tunnel to the server",
	}
	
	if desc, exists := descriptions[method]; exists {
		return desc
	}
	return "Unknown method"
}

func (hc *HTTPClient) ValidateMethod(method string) error {
	if !hc.IsMethodSupported(method) {
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}
	return nil
}

func (hc *HTTPClient) GetMethodRequirements(method HTTPMethod) []string {
	requirements := map[HTTPMethod][]string{
		GET:     {"URL"},
		POST:    {"URL", "Body (optional)"},
		PUT:    {"URL", "Body (optional)"},
		PATCH:  {"URL", "Body (optional)"},
		DELETE: {"URL"},
		HEAD:   {"URL"},
		OPTIONS: {"URL"},
		TRACE:  {"URL"},
		CONNECT: {"URL"},
	}
	
	if req, exists := requirements[method]; exists {
		return req
	}
	return []string{}
}
