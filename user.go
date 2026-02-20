package upstoxapi

import (
	"fmt"
	"net/http"
	"net/url"
)

type loginRequest struct {
	Code           string `form:"code"`
	Client_ID      string `form:"client_id"`
	Client_Secrect string `form:"client_secret"`
	Redirect_uri   string `form:"redirect_uri"`
	GrantType      string `form:"grant_type"`
}

func (c *Client) GetRedirectURL() string {

	u, _ := url.Parse(bASE_URL + rEQUEST_CODE_ENDPOINT)

	q := u.Query()
	q.Set("client_id", c.API_KEY)
	q.Set("status", "querycode")
	// q.Set("redirect_uri", c.Redirect_URI)
	q.Set("response_type", "code")

	u.RawQuery = fmt.Sprintf("%s&redirect_uri=%s", q.Encode(), c.Redirect_URI)
	// u.RawQuery = q.Encode()

	return u.String()
}

type UserProfileResponse struct {
	Email         string   `json:"email"`
	Exchanges     []string `json:"exchanges"`
	Products      []string `json:"products"`
	Broker        string   `json:"broker"`
	UserID        string   `json:"user_id"`
	UserName      string   `json:"user_name"`
	OrderTypes    []string `json:"order_types"`
	UserType      string   `json:"user_type"`
	Poa           bool     `json:"poa"`
	DDPI          bool     `json:"ddpi,omitempty"`
	IsActive      bool     `json:"is_active"`
	AccessToken   string   `json:"access_token,omitempty"`
	ExtendedToken string   `json:"extended_token,omitempty"`
}

func (c *Client) GenerateSession(requestCode string) error {

	var lr = loginRequest{
		Code:           requestCode,
		Client_ID:      c.API_KEY,
		Client_Secrect: c.API_SECRECT,
		Redirect_uri:   c.Redirect_URI,
		GrantType:      "authorization_code",
	}

	var upfr UserProfileResponse
	err := c.doForm(http.MethodPost, aCCESS_TOKEN_ENDPOINT, lr, &upfr)
	if err != nil {
		return err
	}

	c.ACCESS_TOKEN = upfr.AccessToken

	return nil
}

func (c *Client) UserProfile() (*UserProfileResponse, error) {

	var upfr UserProfileResponse
	err := c.doJSON(http.MethodGet, uUSER_PROFILE_ENDPOINT, nil, &upfr)
	if err != nil {
		return nil, err
	}

	return &upfr, err
}

type FundMargin struct {
	UsedMargin      float64 `json:"used_margin"`
	PayInAmount     float64 `json:"payin_amount"`
	SpanMargin      float64 `json:"span_margin"`
	AdhocMargin     float64 `json:"adhoc_margin"`
	NotionalCash    float64 `json:"notional_cash"`
	AvailableMargin float64 `json:"available_margin"`
	ExposureMargin  float64 `json:"exposure_margin"`
}

func (c *Client) UserFundAndMargin() (map[string]FundMargin, error) {

	var fundMargin map[string]FundMargin
	err := c.doJSON(http.MethodGet, uSer_FUND_MARGIN_ENDPOINT, nil, &fundMargin)
	if err != nil {
		return nil, err
	}

	return fundMargin, err
}

func (c *Client) SetAccessToken(accessToken string) {
	c.ACCESS_TOKEN = accessToken
}
