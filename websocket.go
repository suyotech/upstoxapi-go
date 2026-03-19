package upstoxapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/suyotech/upstoxapi-go/pb"
	"google.golang.org/protobuf/proto"
)

type SocketClient struct {
	conn           *websocket.Conn
	accessToken    string
	instrumentkeys []string
	mode           string
	callback       func(Tick)

	ready chan struct{}
}

const (
	ModeLTPC   = "ltpc"
	ModeFull   = "full"
	ModeOption = "option_greeks"
	ModeD30    = "full_d30"

	MethodSubscribe   = "sub"
	MethodUnsubscribe = "unsub"
	MethodChangeMode  = "change_mode"
)

func NewSocketClient(accessToken string) *SocketClient {
	return &SocketClient{
		conn:           nil,
		instrumentkeys: make([]string, 0),
		accessToken:    accessToken,
		mode:           ModeFull,
		ready:          make(chan struct{}),
	}
}

func (s *SocketClient) Connect() error {

	header := http.Header{}
	header.Set("Authorization", "Bearer "+s.accessToken)
	header.Set("Accept", "*/*")

	dialer := websocket.DefaultDialer

	url := "wss://api.upstox.com/v3/feed/market-data-feed"

	conn, resp, err := dialer.Dial(url, header)
	if err != nil {

		// handle redirect
		if resp != nil && resp.StatusCode == http.StatusFound {

			redirectURL := resp.Header.Get("Location")
			log.Println("redirecting to:", redirectURL)

			conn, _, err = dialer.Dial(redirectURL, header)
			if err != nil {
				return err
			}

		} else {
			return err
		}
	}

	log.Println("websocket connected")

	s.conn = conn
	close(s.ready)
	go s.readLoop()

	return nil
}

func (s *SocketClient) readLoop() {

	for {
		msgType, data, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("ws read error:", err)
			return
		}

		switch msgType {

		case websocket.TextMessage:
			log.Println("TEXT MESSAGE:", string(data))

		case websocket.BinaryMessage:
			// log.Println("BINARY MESSAGE:", len(data))

			var resp pb.FeedResponse
			err = proto.Unmarshal(data, &resp)
			if err != nil {
				log.Println("protobuf decode error:", err)
				continue
			}

			s.handleFeed(&resp)

		default:
			log.Println("UNKNOWN MESSAGE TYPE:", msgType)
		}
	}
}

type wsRequest struct {
	Guid   string `json:"guid"`
	Method string `json:"method"`
	Data   struct {
		InstrumentKeys []string `json:"instrumentKeys"`
		Mode           string   `json:"mode,omitempty"`
	} `json:"data"`
}

func (s *SocketClient) send(req wsRequest) error {

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// log.Println("SEND:", string(b))

	return s.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (s *SocketClient) Subscribe(keys []string, mode string) error {

	<-s.ready

	req := wsRequest{
		Guid:   uuid.NewString(),
		Method: MethodSubscribe,
	}

	req.Data.InstrumentKeys = keys
	req.Data.Mode = mode

	return s.send(req)
}

func (s *SocketClient) Unsubscribe(keys []string) error {
	<-s.ready

	req := wsRequest{
		Guid:   uuid.NewString(),
		Method: MethodUnsubscribe,
	}

	req.Data.InstrumentKeys = keys

	return s.send(req)
}

func (s *SocketClient) Disconnect() error {
	return s.conn.Close()
}

type Tick struct {
	Instrument string
	LTP        float64
	Volume     int64
	Timestamp  int64
}

func (s *SocketClient) SetOnTick(f func(Tick)) {
	s.callback = f
}

func (s *SocketClient) handleFeed(resp *pb.FeedResponse) {

	// log.Println("=== FEED RESPONSE ===")
	// log.Println("feed type:", resp.Type)
	// log.Println("number of feeds:", len(resp.Feeds))

	// Check if this is MarketInfo
	if resp.Type == pb.Type_market_info { // MarketInfo
		log.Println("Received MarketInfo - connection established")
		return
	}

	for instrument, feed := range resp.Feeds {
		// log.Println("Processing instrument:", instrument)

		// LTPC mode
		if ltpc := feed.GetLtpc(); ltpc != nil {
			// log.Println("Got LTPC feed for:", instrument)

			tick := Tick{
				Instrument: instrument,
				LTP:        ltpc.Ltp,
				Volume:     ltpc.Ltq,
				Timestamp:  ltpc.Ltt,
			}

			if s.callback != nil {
				s.callback(tick)
			}

			continue
		}

		// FULL mode
		if full := feed.GetFullFeed(); full != nil {
			// log.Println("Got FULL feed for:", instrument)

			market := full.GetMarketFF()
			if market == nil {
				log.Println("market is nil for:", instrument)
				continue
			}

			ltpc := market.GetLtpc()
			if ltpc == nil {
				log.Println("ltpc is nil for:", instrument)
				continue
			}

			tick := Tick{
				Instrument: instrument,
				LTP:        ltpc.GetLtp(),
				Volume:     market.GetVtt(),
				Timestamp:  ltpc.GetLtt(),
			}

			if s.callback != nil {
				s.callback(tick)
			}
		}
	}
}
