package upstoxapi

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const instrument_url = "https://assets.upstox.com/market-quote/instruments/exchange/complete.json.gz"
const filepath = "./instruments.json"

func DownloadInstruments() error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(instrument_url)
	if err != nil {
		return err
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	_, err = io.Copy(out, gzReader)
	if err != nil {
		return nil
	}

	return nil
}

func LoadInstruments() ([]Instrument, error) {

	data, err := os.ReadFile("./instruments.json")
	if err != nil {
		return nil, err
	}

	var instruments []Instrument
	if err := json.Unmarshal(data, &instruments); err != nil {
		return nil, err
	}

	return instruments, nil
}

var (
	SEGMENT_NSE_INDEX = "NSE_INDEX"
	SEGMENT_NSE_EQ    = "NSE_EQ"
	SEGMENT_NSE_FO    = "NSE_FO"
	SEGMENT_BSE_INDEX = "BSE_INDEX"
	SEGMENT_BSE_EQ    = "BSE_EQ"
	SEGMENT_BSE_FO    = "BSE_FO"
	SEGMENT_NCD_FO    = "NCD_FO"
	SEGMENT_NSE_COM   = "NSE_COM"
	SEGMENT_BCD_FO    = "BCD_FO"
	SEGMENT_MCX_FO    = "MCX_FO"

	INSTRUMENT_TYPE_EQ    = "EQ"
	INSTRUMENT_TYPE_FUT   = "FUT"
	INSTRUMENT_TYPE_CE    = "CE"
	INSTRUMENT_TYPE_PE    = "PE"
	INSTRUMENT_TYPE_INDEX = "INDEX"
	INSTRUMENT_TYPE_COM   = "COM"
)

type InstrumentFilter struct {
	Segment        string
	Name           string
	TradingSymbol  string
	Exchange       string
	InstrumentType string
	Expiry         int64
	StrikePrice    float64
}

func FindInstruments(filter InstrumentFilter, instruments []Instrument) ([]Instrument, error) {

	if instruments == nil {
		var err error
		instruments, err = LoadInstruments()
		if err != nil {
			return nil, err
		}
	}

	var filtered []Instrument

	for _, instrument := range instruments {

		//by segment
		if filter.Segment != "" && instrument.Segment != filter.Segment {
			continue
		}
		//by name
		if filter.Name != "" && !strings.HasPrefix(instrument.Name, filter.Name) {
			continue
		}

		//by trading symbol
		if filter.TradingSymbol != "" && !strings.HasPrefix(instrument.TradingSymbol, filter.TradingSymbol) {
			continue
		}

		//by exchange
		if filter.Exchange != "" && instrument.Exchange != filter.Exchange {
			continue
		}

		//by instrument type
		if filter.InstrumentType != "" && instrument.InstrumentType != filter.InstrumentType {
			continue
		}

		//by expiry
		if filter.Expiry != 0 && (instrument.Expiry == nil || *instrument.Expiry != filter.Expiry) {
			continue
		}

		//by strike price
		if filter.StrikePrice != 0 && (instrument.StrikePrice == nil || *instrument.StrikePrice != filter.StrikePrice) {
			continue
		}

		filtered = append(filtered, instrument)
	}
	return filtered, nil
}

type Instrument struct {
	// Common
	Segment        string `json:"segment"`
	Name           string `json:"name"`
	Exchange       string `json:"exchange"`
	InstrumentType string `json:"instrument_type"`
	InstrumentKey  string `json:"instrument_key"`
	ExchangeToken  string `json:"exchange_token"`
	TradingSymbol  string `json:"trading_symbol"`

	// Equity specific
	ISIN          *string  `json:"isin,omitempty"`
	ShortName     *string  `json:"short_name,omitempty"`
	SecurityType  *string  `json:"security_type,omitempty"`
	QtyMultiplier *float64 `json:"qty_multiplier,omitempty"`

	// Derivatives (Futures / Options)
	Weekly           *bool    `json:"weekly,omitempty"`
	Expiry           *int64   `json:"expiry,omitempty"` // epoch millis
	UnderlyingSymbol *string  `json:"underlying_symbol,omitempty"`
	UnderlyingKey    *string  `json:"underlying_key,omitempty"`
	UnderlyingType   *string  `json:"underlying_type,omitempty"`
	StrikePrice      *float64 `json:"strike_price,omitempty"`
	MinimumLot       *int     `json:"minimum_lot,omitempty"`

	// Lot & Tick
	LotSize        *int     `json:"lot_size,omitempty"`
	FreezeQuantity *float64 `json:"freeze_quantity,omitempty"`
	TickSize       *float64 `json:"tick_size,omitempty"`

	// MTF
	MTFEnabled *bool    `json:"mtf_enabled,omitempty"`
	MTFBracket *float64 `json:"mtf_bracket,omitempty"`

	// MIS / Intraday
	IntradayMargin   *float64 `json:"intraday_margin,omitempty"`
	IntradayLeverage *float64 `json:"intraday_leverage,omitempty"`
}

func (i *Instrument) ExpiryToTime(expiry string) (*int64, error) {

	epoch, err := time.Parse("2006-01-02", expiry)
	if err != nil {
		return nil, err
	}

	unixMilli := epoch.UnixMilli()
	return &unixMilli, nil
}

//Equity
// {
//   "segment": "NSE_EQ",
//   "name": "JOCIL LIMITED",
//   "exchange": "NSE",
//   "isin": "INE839G01010",
//   "instrument_type": "EQ",
//   "instrument_key": "NSE_EQ|INE839G01010",
//   "lot_size": 1,
//   "freeze_quantity": 100000.0,
//   "exchange_token": "16927",
//   "tick_size": 5.0,
//   "trading_symbol": "JOCIL",
//   "short_name": "JOCIL",
//   "security_type": "NORMAL"
// }

//Futures
// {
//   "weekly": false,
//   "segment": "NSE_FO",
//   "name": "071NSETEST",
//   "exchange": "NSE",
//   "expiry": 2111423399000,
//   "instrument_type": "FUT",
//   "underlying_symbol": "071NSETEST",
//   "instrument_key": "NSE_FO|36702",
//   "lot_size": 50,
//   "freeze_quantity": 100000.0,
//   "exchange_token": "36702",
//   "minimum_lot": 50,
//   "underlying_key": "NSE_EQ|DUMMYSAN011",
//   "tick_size": 5.0,
//   "underlying_type": "EQUITY",
//   "trading_symbol": "071NSETEST FUT 27 NOV 36",
//   "strike_price": 0.0
// }

//Options
// {
//   "weekly": false,
//   "segment": "NSE_FO",
//   "name": "VODAFONE IDEA LIMITED",
//   "exchange": "NSE",
//   "expiry": 1706207399000,
//   "instrument_type": "CE",
//   "underlying_symbol": "IDEA",
//   "instrument_key": "NSE_FO|36708",
//   "lot_size": 80000,
//   "freeze_quantity": 1600000.0,
//   "exchange_token": "36708",
//   "minimum_lot": 80000,
//   "underlying_key": "NSE_EQ|INE669E01016",
//   "tick_size": 5.0,
//   "underlying_type": "EQUITY",
//   "trading_symbol": "IDEA 22 CE 25 JAN 24",
//   "strike_price": 22.0
// }

//Index
// {
//   "segment": "BSE_INDEX",
//   "name": "AUTO",
//   "exchange": "BSE",
//   "instrument_type": "INDEX",
//   "instrument_key": "BSE_INDEX|AUTO",
//   "exchange_token": "13",
//   "trading_symbol": "AUTO"
// }

//suspended
// {
//   "segment": "NSE_EQ",
//   "name": "JOCIL LIMITED",
//   "exchange": "NSE",
//   "isin": "INE839G01010",
//   "instrument_type": "BE",
//   "instrument_key": "NSE_EQ|INE839G01010",
//   "lot_size": 1,
//   "freeze_quantity": 100000.0,
//   "exchange_token": "16931",
//   "tick_size": 1.0,
//   "trading_symbol": "JOCIL",
//   "qty_multiplier": 1.0
// }

//MTF
// {
//   "segment": "NSE_EQ",
//   "name": "RELIANCE INDUSTRIES LTD",
//   "exchange": "NSE",
//   "isin": "INE002A01018",
//   "instrument_type": "EQ",
//   "instrument_key": "NSE_EQ|INE002A01018",
//   "lot_size": 1,
//   "freeze_quantity": 100000.0,
//   "exchange_token": "2885",
//   "tick_size": 5.0,
//   "trading_symbol": "RELIANCE",
//   "short_name": "Reliance Industries",
//   "mtf_enabled": true,
//   "mtf_bracket": 26.5,
//   "security_type": "NORMAL"
// }

//MIS
// {
//   "segment": "NSE_EQ",
//   "name": "RELIANCE INDUSTRIES LTD",
//   "exchange": "NSE",
//   "isin": "INE002A01018",
//   "instrument_type": "EQ",
//   "instrument_key": "NSE_EQ|INE002A01018",
//   "lot_size": 1,
//   "freeze_quantity": 100000,
//   "exchange_token": "2885",
//   "tick_size": 10,
//   "trading_symbol": "RELIANCE",
//   "short_name": "Reliance Industries",
//   "qty_multiplier": 1,
//   "security_type": "NORMAL",
//   "intraday_margin": 20,
//   "intraday_leverage": 5
// }
