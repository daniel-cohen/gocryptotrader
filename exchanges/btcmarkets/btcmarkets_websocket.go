package btcmarkets

import (
	"fmt"
	"log"

	"github.com/thrasher-/gocryptotrader/common"
	"github.com/thrasher-/socketio"
)

const (
	BTCMARKETS_SOCKETIO_ADDRESS = "https://socket.btcmarkets.net"
)

var BTCMarketsSocket *socketio.SocketIO

func (b *BTCMarkets) OnConnect(output chan socketio.Message) {
	if b.Verbose {
		log.Printf("%s Connected to Websocket.", b.GetName())
	}

	currencies := []string{}
	for _, x := range b.EnabledPairs {
		currency := common.StringToUpper(x)
		//currency := common.StringToUpper(x[3:] + x[0:3])
		currencies = append(currencies, currency)
	}

	// Note the format is differnt between channels. E.g:
	// Orderbook_XRPAUD
	// TRADE_XRPAUD
	// Ticker-BTCMarkets-XRP-AUD
	endpoints := []string{"Orderbook"}
	//TODO:
	//endpoints := []string{"Orderbook_", "TRADE_"}

	// Register to orderbook Events:
	for _, x := range endpoints {
		for _, y := range currencies {
			channel := fmt.Sprintf(`"%s_%s"`, x, y)
			if b.Verbose {
				log.Printf("%s Websocket subscribing to channel: %s.", b.GetName(), channel)
			}
			output <- socketio.CreateMessageEvent("subscribe", channel, b.OnMessage, BTCMarketsSocket.Version)
		}
	}
}

func (b *BTCMarkets) OnDisconnect(output chan socketio.Message) {
	log.Printf("%s Disconnected from websocket server.. Reconnecting.\n", b.GetName())
	b.WebsocketClient()
}

func (b *BTCMarkets) OnError() {
	log.Printf("%s Error with Websocket connection.. Reconnecting.\n", b.GetName())
	b.WebsocketClient()
}

func (b *BTCMarkets) OnMessage(message []byte, output chan socketio.Message) {
	if b.Verbose {
		log.Printf("%s Websocket message received which isn't handled by default.\n", b.GetName())
		log.Println(string(message))
	}
}

// func (b *BTCMarkets) OnTicker(message []byte, output chan socketio.Message) {
// 	type Response struct {
// 		Ticker WebsocketTicker `json:"ticker"`
// 	}
// 	var resp Response
// 	err := common.JSONDecode(message, &resp)

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

func (b *BTCMarkets) OnOrderbook(message []byte, output chan socketio.Message) {
	orderbook := WebsocketOrderbok{}
	err := common.JSONDecode(message, &orderbook)

	if err != nil {
		log.Println(err)
		return
	}
}

// func (b *BTCMarkets) OnTrade(message []byte, output chan socketio.Message) {
// 	trade := WebsocketTrade{}
// 	err := common.JSONDecode(message, &trade)

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

func (b *BTCMarkets) WebsocketClient() {
	events := make(map[string]func(message []byte, output chan socketio.Message))
	events["OrderBookChange"] = b.OnOrderbook
	//events["ticker"] = b.OnTicker
	//events["trade"] = b.OnTrade

	BTCMarketsSocket = &socketio.SocketIO{
		Version:      1,
		OnConnect:    b.OnConnect,
		OnEvent:      events,
		OnError:      b.OnError,
		OnMessage:    b.OnMessage,
		OnDisconnect: b.OnDisconnect,
	}

	for b.Enabled && b.Websocket {
		err := socketio.ConnectToSocket(BTCMARKETS_SOCKETIO_ADDRESS, BTCMarketsSocket)
		if err != nil {
			log.Printf("%s Unable to connect to Websocket. Err: %s\n", b.GetName(), err)
			continue
		}
		log.Printf("%s Disconnected from Websocket.\n", b.GetName())
	}
}
