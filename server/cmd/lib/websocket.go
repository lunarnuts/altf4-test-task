package lib

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/gosuri/uilive"
)

type request struct {
	Method string    `json:"method"`
	Params [1]string `json:"params"`
	ID     int       `json:"id"`
}

type Message struct {
	EventType     string        `json:"e"`
	EventTime     uint64        `json:"E"`
	Symbol        string        `json:"s"`
	FirstUpdateID uint64        `json:"U"`
	LastUpdateId  uint64        `json:"u"`
	Bids          []Conversions `json:"b"`
	Asks          []Conversions `json:"a"`
}

type Conversions struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"qty"`
}

//custom Unmarshaller to properly encode and decode Message from WS
func (c *Conversions) UnmarshalJSON(b []byte) error {
	tmp := [2]json.Number{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	var err error
	c.Price, err = tmp[0].Float64()
	if err != nil {
		return err
	}
	c.Quantity, err = tmp[1].Float64()
	if err != nil {
		return err
	}
	return nil
}

//opens connection to binance api websocket
func BinanceWS() {
	writer := uilive.New()
	writer.Start()
	conn, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws", nil)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	depth := request{"SUBSCRIBE", [1]string{"btcusdt@depth"}, 1}
	conn.WriteJSON(depth)
	db := ConnectToDb()
	db.DownloadSnapshot()
	for {
		var message Message
		readErr := conn.ReadJSON(&message)
		if readErr != nil {
			log.Println(readErr)
			return
		}
		if message.LastUpdateId <= db.LastUpdateID {
			continue
		}
		if message.FirstUpdateID <= db.LastUpdateID+1 &&
			message.LastUpdateId >= db.LastUpdateID+1 {
			db.Update(message)
		}
		fmt.Fprintf(writer, "%v\n", db)
	}
}
