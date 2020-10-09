package main

import (
	"fmt"
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

type cookedOrder struct {
	Order
	completed time.Time
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
