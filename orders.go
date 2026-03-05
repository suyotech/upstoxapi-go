package upstoxapi

import (
	"net/http"
	"net/url"
)

type PlaceOrderReq struct {
	Quantity          int64   `json:"quantity"`
	Product           string  `json:"product"`
	Validity          string  `json:"validity"`
	Price             float64 `json:"price"`
	Tag               string  `json:"tag"`
	InstrumentToken   string  `json:"intrument_token"`
	OrderType         string  `json:"order_type"`
	TransactionType   string  `json:"transaction_type"`
	DisclosedQuantiry int64   `json:"disclosed_quantity"`
	TriggerPrice      float64 `json:"trigger_price"`
	IsAMO             bool    `json:"is_amo"`
	Slice             bool    `json:"slice"`
}

type OrderResp struct {
	OrderIds []string
	MetaData map[string]string
}

func (c *Client) PlaceOrder(params PlaceOrderReq) (*OrderResp, error) {

	var orderResp OrderResp
	err := c.doJSON(base_url3, http.MethodPost, place_order_endpoint, nil, params, &orderResp)
	if err != nil {
		return nil, err
	}
	return &orderResp, nil
}

type ModifyOrderReq struct {
	Quantity          int64   `json:"quantity"`
	Validity          string  `json:"validity"`
	Price             float64 `json:"price"`
	OrderID           string  `json:"order_id"`
	OrderType         string  `json:"order_type"`
	DisclosedQuantiry int64   `json:"disclosed_quantity"`
	TriggerPrice      float64 `json:"trigger_price"`
}

func (c *Client) ModifyOrder(params ModifyOrderReq) (*OrderResp, error) {
	var orderResp OrderResp
	err := c.doJSON(base_url3, http.MethodPost, modify_order_endpoint, nil, params, &orderResp)
	if err != nil {
		return nil, err
	}
	return &orderResp, nil
}

func (c *Client) CancelOrder(orderdID string) (*OrderResp, error) {
	var orderResp OrderResp

	q := url.Values{}
	q.Set("order_id", orderdID)

	err := c.doJSON(base_url3, http.MethodDelete, cancel_order_endpoint, q, nil, &orderResp)
	if err != nil {
		return nil, err
	}
	return &orderResp, nil
}
