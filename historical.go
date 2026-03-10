package upstoxapi

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
	OI     int64
}

type CandleDataRequest struct {
	InsrumentKey string
	Unit         string
	Timeframe    int64
	FromDate     string
	ToDate       string
}

func (c *Client) GetHistCandles(req *CandleDataRequest) (any, error) {

	instrument := url.PathEscape(req.InsrumentKey)
	urlpath := fmt.Sprintf("%s/%s/%s/%d/%s/%s", get_candle_endpoint, instrument, req.Unit, req.Timeframe, req.ToDate, req.FromDate)
	log.Println("encoded path : ", urlpath)

	var resp struct {
		Candles [][]any `json:"candles"`
	}

	err := c.doJSON(base_url4, "GET", urlpath, nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return convertCandles(resp.Candles)
}

func convertCandles(raw [][]any) ([]Candle, error) {
	candles := make([]Candle, 0, len(raw))

	for _, c := range raw {

		t, err := time.Parse(time.RFC3339, c[0].(string))
		if err != nil {
			return nil, err
		}

		candle := Candle{
			Time:   t,
			Open:   c[1].(float64),
			High:   c[2].(float64),
			Low:    c[3].(float64),
			Close:  c[4].(float64),
			Volume: int64(c[5].(float64)), // JSON numbers decode as float64
		}

		candles = append(candles, candle)
	}

	return candles, nil
}

func (c *Client) GetIntradayCandles(req *CandleDataRequest) (any, error) {

	instrument := url.PathEscape(req.InsrumentKey)
	urlpath := fmt.Sprintf("%s/%s/%s/%d", get_intraday_candle_endpoint, instrument, req.Unit, req.Timeframe)
	log.Println("encoded path : ", urlpath)

	var resp struct {
		Candles [][]any `json:"candles"`
	}

	err := c.doJSON(base_url4, "GET", urlpath, nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return convertCandles(resp.Candles)
}
