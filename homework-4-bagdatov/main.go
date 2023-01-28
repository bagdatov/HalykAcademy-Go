package main

import (
	"homework-4/airport"
	"log"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)

	airport := airport.NewAirport()

	go StopAirport(airport)

	<-airport.Done
}

func StopAirport(a *airport.AirPort) {
	time.Sleep(20 * time.Second)
	a.Stop()
}
