package main

import "fmt"

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

func printShelves() {
	for temp, value := range shelves {
		fmt.Printf("Shelf: %v size: %v\n", temp, value.index)
		value.print()
	}
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
