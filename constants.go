package upstoxapi

import "time"

const (
	base_url1 = "https://api.upstox.com/v2"
	// base_url2 = "https://api.upstox.com/v3"
	base_url3 = "https://hft-api.upstox.com/v3"
	base_url4 = "https://api.upstox.com/v3"

	// sandbox_url               = "https://api-hft.upstox.com/v3"
	request_timeout              = 30 * time.Second
	request_code_url             = "/login/authorization/dialog"
	access_token_endpoint        = "/login/authorization/token"
	user_profile_endpoint        = "/user/profile"
	user_fund_margin_endpoint    = "/user/get-funds-and-margin"
	place_order_endpoint         = "/order/place"
	modify_order_endpoint        = "/order/modify"
	cancel_order_endpoint        = "/order/cancel"
	get_candle_endpoint          = "/historical-candle"
	get_intraday_candle_endpoint = "/historical-candle/intraday"
)
