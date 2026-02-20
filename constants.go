package upstoxapi

import "time"

const (
	bASE_URL                  = "https://api.upstox.com/v2"
	rEQUEST_TIMEOUT           = 30 * time.Second
	rEQUEST_CODE_ENDPOINT     = "/login/authorization/dialog"
	aCCESS_TOKEN_ENDPOINT     = "/login/authorization/token"
	uUSER_PROFILE_ENDPOINT    = "/user/profile"
	uSer_FUND_MARGIN_ENDPOINT = "/user/get-funds-and-margin"
)
