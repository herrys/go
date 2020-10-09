package main

import "fmt"

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
