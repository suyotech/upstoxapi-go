package upstoxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
)

type Client struct {
	API_KEY      string
	API_SECRECT  string
	Redirect_URI string
	ACCESS_TOKEN string
	baseUrl      string
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

		Timeout: rEQUEST_TIMEOUT,
	}

	return &Client{
		API_KEY:      apikey,
		API_SECRECT:  apisecrect,
		Redirect_URI: redirect_uri,
		baseUrl:      bASE_URL,
		debug:        false,
		httpClient:   httpclient,
	}
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) doRequest(
	method string,
	path string,
	body any,
	contentType string,
	result any,
) error {

	// Build URL
	url := c.baseUrl + path

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
		fmt.Printf("[DEBUG] Request: %s %s\n", method, url)
	}

	req, err := http.NewRequest(method, url, bodyReader)
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

	// Decode response JSON
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %w; raw: %s",
			err, string(respBody))
	}

	if c.debug {
		fmt.Printf("[DEBUG] Decoded Result: %+v\n\n\n", result)
	}

	return nil
}

func (c *Client) doJSON(method string,
	path string,
	body any,
	result any) error {
	return c.doRequest(method, path, body, "application/json", result)
}

func (c *Client) doForm(method string,
	path string,
	body any,
	result any) error {
	return c.doRequest(method, path, body, "application/x-www-form-urlencoded", result)
}

