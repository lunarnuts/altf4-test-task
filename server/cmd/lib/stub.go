package lib

import (
	"fmt"
	"log"
	"strings"
)

var db_stub DbStub

//for orderbook with capacity of 15, no need to establish DB, stub is sufficient
type DbStub struct {
	LastUpdateID uint64
	Bids         map[float64]float64
	Asks         map[float64]float64
	Count        int
	BVol         float64
	AVol         float64
}

//stub for database connection
func ConnectToDb() *DbStub {
	db_stub = DbStub{
		LastUpdateID: 0,
		Bids:         make(map[float64]float64, 15),
		Asks:         make(map[float64]float64, 15),
		Count:        0,
		BVol:         0,
		AVol:         0,
	}
	return &db_stub
}

//implements Stringer interface for stub
func (db DbStub) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "LastUpdateID: %v\n", db.LastUpdateID)
	fmt.Fprint(&b, "Bids: [")
	for k, v := range db.Bids {
		fmt.Fprintf(&b, " {%.5f %.5f}", k, v)
	}
	fmt.Fprint(&b, " ]\nAsks: [")
	for k, v := range db.Asks {
		fmt.Fprintf(&b, " {%.5f %.5f}", k, v)
	}
	b.WriteString(" ]\n")
	fmt.Fprintf(&b, "Bids Volume: %.5f\nAsks Volume: %.5f", db.BVol, db.AVol)
	return b.String()
}

//calculates volume
func CalculateVolume(c map[float64]float64) float64 {
	sum := float64(0)
	for k, v := range c {
		sum += (k * v)
	}
	return sum
}

//initial download of snapshot from binance rest api
func (db *DbStub) DownloadSnapshot() {
	data, err := GetDataFrames()
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	db.LastUpdateID = data.LastUpdateId
	for _, v := range data.Asks {
		db.Asks[v.Price] = v.Quantity
	}
	for _, v := range data.Bids {
		db.Bids[v.Price] = v.Quantity
	}
	db.Count = len(data.Asks)
	db.BVol = CalculateVolume(db.Bids)
	db.AVol = CalculateVolume(db.Asks)
}

//stub imitatation of update to db
func (db *DbStub) Update(m Message) {
	for _, v := range m.Asks {
		_, ok := db.Asks[v.Price]
		if ok {
			if v.Quantity == 0 {
				delete(db.Asks, v.Price)
			} else {
				db.Asks[v.Price] = v.Quantity
			}
		}
		if len(db.Asks) < 15 && v.Quantity > 0 {
			db.Asks[v.Price] = v.Quantity
		}
	}
	for _, v := range m.Bids {
		_, ok := db.Bids[v.Price]
		if ok {
			if v.Quantity == 0 {
				delete(db.Bids, v.Price)
			} else {
				db.Bids[v.Price] = v.Quantity
			}
		}
		if len(db.Bids) < 15 && v.Quantity > 0 {
			db.Bids[v.Price] = v.Quantity
		}
	}
	db.BVol = CalculateVolume(db.Bids)
	db.AVol = CalculateVolume(db.Asks)
	db.LastUpdateID = m.LastUpdateId
}
