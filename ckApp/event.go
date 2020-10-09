package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type event interface {
	do()
}

type pickupEvent struct {
	courierID int
}

func (pe pickupEvent) do() {
	fmt.Printf("\n***Pickup Event courierId:%v\n", pe.courierID)
	sh := shelfWithOldestItem()
	if sh != nil {
		co := sh.get()
		fmt.Printf("\tPickup: shelf=%v item=%v\n", sh.temp, co.Name)
	}
}

type orderEvent struct {
	Order
}

func (oe orderEvent) do() {
	fmt.Printf("\n***Order Event id:%v name:%v\n", oe.ID, oe.Name)
	cook(oe.Order)
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
