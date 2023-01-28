package main

import (
	"fmt"

	"homework4.2/airport"
)

func main() {

	a := airport.NewAirport()

	planes := a.Start()
	a.Close(15)

	for _, plane := range planes {
		fmt.Printf("%#v\n", plane)
	}
}
