package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	swedavia "github.com/seb-emmot/sviago/swedavia"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run fetch.go <airportIATA> <date> [outputdir]")
		return
	}

	airport := os.Args[1]
	date := os.Args[2]

	// read environment variable
	sKey, ok := os.LookupEnv("SWEDAVIA_SUBSCRIPTION_KEY")
	if !ok {
		log.Fatal("SWEDAVIA_SUBSCRIPTION_KEY not set")
		return
	}

	client := swedavia.Client{
		URL:             "https://api.swedavia.se",
		SubscriptionKey: sKey,
	}

	arrivalsInfo, err := client.GetArrivals(airport, date)

	if err != nil {
		fmt.Println("Error getting flight info:", err)
		return
	}

	departuresInfo, err := client.GetDepartures(airport, date)

	if err != nil {
		fmt.Println("Error getting flight info:", err)
		return
	}

	arrivalFname := fmt.Sprintf("data/arrivals_%s_%s.json", airport, date)
	departFname := fmt.Sprintf("data/departures_%s_%s.json", airport, date)

	if len(os.Args) == 4 {
		dir := os.Args[3]
		arrivalFname = fmt.Sprintf("%s/%s", dir, arrivalFname)
		departFname = fmt.Sprintf("%s/%s", dir, departFname)
	}

	adir := filepath.Dir(arrivalFname)
	os.MkdirAll(adir, os.ModePerm)

	// save arrivals to file
	file, err := os.Create(arrivalFname)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(arrivalsInfo)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Printf("Arrivals written to %s", arrivalFname)

	// save departures to file
	file, err = os.Create(departFname)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder = json.NewEncoder(file)
	err = encoder.Encode(departuresInfo)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Printf("Departures written to %s", departFname)
}

func PrintArrivals(a swedavia.ArrivalsInfo) {
	fmt.Println("Arrival Airport:", a.To.ArrivalAirportIata)
	fmt.Println("Flight Arrival Date:", a.To.FlightArrivalDate)
	fmt.Println("Number of Flights:", a.NumberOfFlights)

	for _, flight := range a.Flights {
		fmt.Println("Flight ID:", flight.FlightId)
		fmt.Println("Scheduled Arrival Time:", flight.ArrivalTime.ScheduledUtc)
		fmt.Println("Actual Arrival Time:", flight.ArrivalTime.ActualUtc)
		fmt.Println("Departure Airport:", flight.DepartureAirportEnglish)
		fmt.Println("Airline Operator:", flight.AirlineOperator.Name)
	}
}

func PrintDepartures(d swedavia.DeparturesInfo) {
	fmt.Println("Departure Airport:", d.From.DepartureAirportIata)
	fmt.Println("Flight Departure Date:", d.From.FlightDepartureDate)
	fmt.Println("Number of Flights:", d.NumberOfFlights)

	for _, flight := range d.Flights {
		fmt.Println("Flight ID:", flight.FlightId)
		fmt.Println("Scheduled Departure Time:", flight.DepartureTime.ScheduledUtc)
		fmt.Println("Arrival Airport:", flight.ArrivalAirportEnglish)
		fmt.Println("Airline Operator:", flight.AirlineOperator.Name)
	}
}
