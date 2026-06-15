package upstoxapi

import (
	"net/http"
	"net/url"
)

// PlaceOrderReq represents the payload used to submit an order to the broker.
// All fields map directly to broker order parameters.
type PlaceOrderReq struct {

	// Quantity is the number of units/lots to trade.
	// Must be a valid lot size for the instrument.
	// Example: 1, 25, 50
	Quantity int64 `json:"quantity"`

	// Product defines the product/margin type used for the order.
	// Common values:
	//   I   - Intraday
	//   D   - Delivery
	//   MTF - Margin Trading Facility
	Product string `json:"product"`

	// Validity defines how long the order remains active.
	// Supported values:
	//   DAY - Valid for the trading day
	//   IOC - Immediate Or Cancel
	Validity string `json:"validity"`

	// Price is the order price.
	// Required for LIMIT orders.
	// Should be 0 for MARKET orders.
	Price float64 `json:"price"`

	// Tag is an optional custom identifier used by client systems
	// to track or group orders (strategy ID, user ID, etc).
	Tag string `json:"tag"`

	// InstrumentToken uniquely identifies the tradable instrument
	// in the broker system.
	// Example: "NSE_FO|12345"
	InstrumentToken string `json:"instrument_token"`

	// OrderType defines how the order is executed.
	// Supported values:
	//   MARKET - Executes at best available price
	//   LIMIT  - Executes at specified price
	//   SL     - Stop Loss Limit
	//   SL-M   - Stop Loss Market
	OrderType string `json:"order_type"`

	// TransactionType specifies order direction.
	// Supported values:
	//   BUY
	//   SELL
	TransactionType string `json:"transaction_type"`

	// DisclosedQuantity is the portion of the total quantity
	// visible in the market order book.
	// Default: 0 (entire quantity disclosed).
	DisclosedQuantiry int64 `json:"disclosed_quantity"`

	// TriggerPrice is used for stop-loss orders (SL / SL-M).
	// Ignored for MARKET and LIMIT orders.
	// Default: 0
	TriggerPrice float64 `json:"trigger_price"`

	// IsAMO indicates whether the order is an After Market Order.
	// true  -> AMO order
	// false -> Regular market order
	IsAMO bool `json:"is_amo"`

	// Slice enables order slicing when quantity exceeds exchange limits.
	// true  -> order will be automatically split
	// false -> send as a single order
	Slice bool `json:"slice"`
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
	//ex. 1, 2 etc
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

type OrderBookItem struct {
	Exchange          string  `json:"exchange"`
	Product           string  `json:"product"`
	Price             float64 `json:"price"`
	Quantity          int64   `json:"quantity"`
	Status            string  `json:"status"`
	Guid              *string `json:"guid"`
	Tag               *string `json:"tag"`
	InstrumentToken   string  `json:"instrument_token"`
	PlacedBy          string  `json:"placed_by"`
	TradingSymbol     string  `json:"trading_symbol"`
	TradingSymbolAlt  string  `json:"tradingsymbol"`
	OrderType         string  `json:"order_type"`
	Validity          string  `json:"validity"`
	TriggerPrice      float64 `json:"trigger_price"`
	DisclosedQuantity int64   `json:"disclosed_quantity"`
	TransactionType   string  `json:"transaction_type"`
	AveragePrice      float64 `json:"average_price"`
	FilledQuantity    int64   `json:"filled_quantity"`
	PendingQuantity   int64   `json:"pending_quantity"`
	StatusMessage     *string `json:"status_message"`
	StatusMessageRaw  *string `json:"status_message_raw"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	ParentOrderID     *string `json:"parent_order_id"`
	OrderID           string  `json:"order_id"`
	Variety           string  `json:"variety"`
	OrderTimestamp    string  `json:"order_timestamp"`
	ExchangeTimestamp *string `json:"exchange_timestamp"`
	IsAMO             bool    `json:"is_amo"`
	OrderRequestID    string  `json:"order_request_id"`
	OrderRefID        string  `json:"order_ref_id"`
}

func (c *Client) GetOrderBook() ([]OrderBookItem, error) {
	var orders []OrderBookItem

	err := c.doJSON(base_url5, http.MethodGet, get_order_book_endpoint, nil, nil, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
