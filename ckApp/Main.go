package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"container/list"
)

// Order is a basic class for an incoming order
type Order struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ShelfLive int     `json:"shelfLife"`
	DecayRate float32 `json:"decayRate"`
}

// Orders is an array of orders
type Orders struct {
	Orders []Order `json:"orders"`
}

type shelf struct {
	temp         string
	capacity     int
	cookedOrders list
}

func (sh shelf) isFull() bool {
	return
}

func (sh shelf) store(co cookedOrder) {

}

type cookedOrder struct {
	Order
	completed time.Time
}

type event interface {
	do()
}

type pickupEvent struct {
	courierID int
}

func (pe pickupEvent) do() {
	fmt.Printf("Pickup Event courierId:%v\n", pe.courierID)
}

type orderEvent struct {
	Order
}

func (oe orderEvent) do() {
	fmt.Printf("Order Event id:%v name:%v\n", oe.ID, oe.Name)
}

func generatePickup(eventChannel chan event, count int) {
	for i := 0; i < count; i++ {
		msTime := rand.Intn(4000) + 2000
		time.Sleep(time.Duration(msTime) * time.Millisecond)
		var puEvent pickupEvent
		puEvent.courierID = i
		eventChannel <- puEvent
	}
}

func generateOrder(eventChannel chan event) {
	jsonFile, err := os.Open("/Users/sutanto/src/go/src/github.com/herrys/ckApp/orders.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened orders.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var orders Orders

	json.Unmarshal(byteValue, &orders)

	for i := 0; i < len(orders.Orders); i++ {
		time.Sleep(500 * time.Millisecond)
		// fmt.Println(i, " Id: ", orders.Orders[i].ID, " Name: ", orders.Orders[i].Name)
		var oEvent orderEvent
		oEvent.Order = orders.Orders[i]
		eventChannel <- oEvent
	}
	fmt.Println("Done")
}

func main() {
	eventChannel := make(chan event)
	go generateOrder(eventChannel)
	go generatePickup(eventChannel, 100)
	for i := 0; i < 100; i++ {
		e := <-eventChannel
		e.do()
	}
	fmt.Println("Done")
}
