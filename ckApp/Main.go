package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// Order is a basic class for an incoming order
type Order struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Temp      string  `json:"temp"`
	ShelfLive int     `json:"shelfLife"`
	DecayRate float32 `json:"decayRate"`
}

// Orders is an array of orders
type Orders struct {
	Orders []Order `json:"orders"`
}

type shelf struct {
	temp  string
	index int
	items []cookedOrder
}

var shelves map[string]*shelf

func (sh *shelf) store(co cookedOrder) bool {
	if sh.index >= len(sh.items) {
		return false
	}
	sh.items[sh.index] = co
	sh.index++
	return true
}

func (sh *shelf) peek() cookedOrder {
	var co cookedOrder
	if sh.index > 0 {
		co = sh.items[0]
	}
	return co
}

func (sh *shelf) get() cookedOrder {
	var co cookedOrder
	if sh.index > 0 {
		co = sh.items[0]
		sh.items = sh.items[1:]
		sh.index--
	}
	return co
}

func shelfWithOldestItem() *shelf {
	var sh *shelf
	for _, shelfValue := range shelves {
		item := shelfValue.peek()
		if (item != cookedOrder{}) && ((sh == nil) || (item.completed.Before(sh.peek().completed))) {
			sh = shelfValue
		}
	}
	return sh
}

func (sh shelf) print() {
	for i := 0; i < sh.index; i++ {
		fmt.Printf("\t %v name:%v\n", i, sh.items[i].Name)
	}
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

func cook(order Order) {
	var co cookedOrder
	co.Order = order
	co.completed = time.Now()
	fmt.Printf("\tBefore shelve %v contains: %v\n", shelves[order.Temp].temp, shelves[order.Temp].index)
	if !shelves[order.Temp].store(co) {
		if !shelves["any"].store(co) {
			fmt.Println("ALL Shelves are full")
		}
	}
	fmt.Printf("\tAfter shelve %v contains: %v\n", shelves[order.Temp].temp, shelves[order.Temp].index)
}

func printShelves() {
	for temp, value := range shelves {
		fmt.Printf("Shelf: %v size: %v\n", temp, value.index)
		value.print()
	}
}

func pickupOldestItem() {

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

func initShelves() {
	var hotShelf shelf
	hotShelf.temp = "hot"
	hotShelf.items = make([]cookedOrder, 10)

	var coldShelf shelf
	coldShelf.temp = "cold"
	coldShelf.items = make([]cookedOrder, 10)

	var frozenShelf shelf
	frozenShelf.temp = "frozen"
	frozenShelf.items = make([]cookedOrder, 10)

	var anyShelf shelf
	anyShelf.temp = "any"
	anyShelf.items = make([]cookedOrder, 15)
	shelves = make(map[string]*shelf)
	shelves["hot"] = &hotShelf
	shelves["cold"] = &coldShelf
	shelves["frozen"] = &frozenShelf
	shelves["any"] = &anyShelf
}

func main() {
	initShelves()
	eventChannel := make(chan event)
	go generateOrder(eventChannel)
	go generatePickup(eventChannel, 100)
	for i := 0; i < 100; i++ {
		e := <-eventChannel
		e.do()
		printShelves()
	}
	fmt.Println("Done")
}
