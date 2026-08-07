package upstoxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
)

type Client struct {
	API_KEY      string
	API_SECRECT  string
	Redirect_URI string
	ACCESS_TOKEN string
	debug        bool
	httpClient   *http.Client
}

type serverResponse[T any] struct {
	Status string     `json:"status"`
	Errors []ApiError `json:"errors,omitempty"`
	Data   T          `json:"data,omitempty"`
}

type ApiError struct {
	ErrorCode    string `json:"error_code"`
	Message      string `json:"message"`
	PropertyPath any    `json:"property_path"`
	InvalidValue any    `json:"invalid_value"`
}

func NewClient(apikey, apisecrect, redirect_uri string) *Client {
	httpclient := &http.Client{
		Timeout: request_timeout,
	}

	return &Client{
		API_KEY:      apikey,
		API_SECRECT:  apisecrect,
		Redirect_URI: redirect_uri,
		debug:        false,
		httpClient:   httpclient,
	}
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// SetProxy configures a static HTTP(S) proxy for all REST requests.
// Pass an empty string to restore the default HTTP transport.
func (c *Client) SetProxy(proxyURL string) error {
	if proxyURL == "" {
		c.httpClient.Transport = nil
		return nil
	}

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}
	if proxy.Host == "" || (proxy.Scheme != "http" && proxy.Scheme != "https") {
		return fmt.Errorf("invalid proxy URL: expected http:// or https:// proxy with host")
	}

	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return fmt.Errorf("unsupported default HTTP transport type %T", http.DefaultTransport)
	}
	proxyTransport := transport.Clone()
	proxyTransport.Proxy = http.ProxyURL(proxy)
	c.httpClient.Transport = proxyTransport
	return nil
}

func (c *Client) doRequest(
	baseurl string,
	method string,
	path string,
	query url.Values,
	body any,
	contentType string,
	result any,
) error {

	// Build URL
	u, err := url.Parse(baseurl + path)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	finalURL := u.String()

	// Prepare Body
	var bodyReader io.Reader
	if body != nil {

		switch contentType {

		case "application/json":
			jsonBytes, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal json body: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonBytes)

			if c.debug {
				fmt.Printf("[DEBUG] Request Body (JSON): %s\n", string(jsonBytes))
			}

		case "application/x-www-form-urlencoded":
			formEncoder := form.NewEncoder()
			values, err := formEncoder.Encode(body)
			if err != nil {
				return fmt.Errorf("failed to encode form body: %w", err)
			}
			bodyReader = strings.NewReader(values.Encode())

			if c.debug {
				fmt.Printf("[DEBUG] Request Body (FORM): %s\n", values.Encode())
			}

		default:
			return fmt.Errorf("unsupported content type: %s", contentType)
		}
	}

	if c.debug {
		fmt.Printf("[DEBUG] Request: %s %s\n", method, finalURL)
	}

	req, err := http.NewRequest(method, finalURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")

	if c.ACCESS_TOKEN != "" {
		req.Header.Set("Authorization", "Bearer "+c.ACCESS_TOKEN)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Request Headers:\n")
		for k, v := range req.Header {
			fmt.Printf(" %s: %v\n", k, v)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response Status: %s\n", resp.Status)
		fmt.Printf("[DEBUG] Response Body: %s\n", string(respBody))
	}

	// HTTP error check
	// If HTTP error, decode API error body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		var errResp serverResponse[any]
		if err := json.Unmarshal(respBody, &errResp); err == nil && len(errResp.Errors) > 0 {
			return fmt.Errorf("api error [%s]: %s",
				errResp.Errors[0].ErrorCode,
				errResp.Errors[0].Message)
		}

		return fmt.Errorf("http error %d: %s",
			resp.StatusCode, string(respBody))
	}

	if result == nil {
		return nil
	}

	// Try to detect wrapped response
	var envelope struct {
		Status string          `json:"status"`
		Data   json.RawMessage `json:"data"`
		Errors json.RawMessage `json:"errors"`
	}

	if err := json.Unmarshal(respBody, &envelope); err == nil && envelope.Data != nil {
		// Wrapped response → decode only data
		if err := json.Unmarshal(envelope.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal wrapped data: %w; raw: %s",
				err, string(envelope.Data))
		}
	} else {
		// Direct response → decode full body
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response JSON: %w; raw: %s",
				err, string(respBody))
		}
	}

	if c.debug {
		fmt.Printf("[DEBUG] Decoded Result: %+v\n\n", result)
	}

	return nil
}

func (c *Client) doJSON(baseurl, method string,
	path string,
	query url.Values,
	body any,
	result any) error {
	return c.doRequest(baseurl, method, path, query, body, "application/json", result)
}

func (c *Client) doForm(baseurl, method string,
	path string,
	query url.Values,
	body any,
	result any) error {
	return c.doRequest(baseurl, method, path, query, body, "application/x-www-form-urlencoded", result)
}
