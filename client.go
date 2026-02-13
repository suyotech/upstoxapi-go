package upstoxapi

import (
	"net/http"
	"net/url"
)

type Client struct {
	API_KEY      string
	API_SECRECT  string
	Redirect_URI string
	ACCESS_TOKEN string
	httpClient   *http.Client
}

func NewClient(apikey, apisecrect, redirect_uri string) *Client {

	httpclient := &http.Client{

		Timeout: REQUEST_TIMEOUT,
	}

	return &Client{
		API_KEY:      apikey,
		API_SECRECT:  apisecrect,
		Redirect_URI: redirect_uri,
		httpClient:   httpclient,
	}
}

func (c *Client) GetRedirectURL() string {

	u := url.URL{
		Scheme: "https",
		Host:   BASE_URL,
		Path:   REQUEST_CODE_ENDPOINT,
	}

	q := u.Query()
	q.Set("client_id", c.API_KEY)
	q.Set("redirect_uri", c.Redirect_URI)
	q.Set("status", "querycode")
	q.Set("response_type", "code")

	u.RawQuery = q.Encode()

	return u.String()
}

func (c *Client) GenerateSession(requestCode string) error {

	return nil
}
