package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type DataFrame struct {
	LastUpdateId uint64        `json:"lastUpdateId"`
	Bids         []Conversions `json:"bids"`
	Asks         []Conversions `json:"asks"`
}

//connect to binance api to get dataframe
func GetDataFrames() (data DataFrame, err error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v3/depth?symbol=BTCUSDT&&limit=20", nil)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return data, err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return data, err
	}
	defer client.CloseIdleConnections()
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return data, err
	}
	data.Asks = data.Asks[:15]
	data.Bids = data.Bids[:15]
	return data, nil
}
