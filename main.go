package main

import (
	"fmt"
	"os"

	swedavia "github.com/seb-emmot/sviago/swedavia"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <sub key> <airportIATA> <date>")
		return
	}

	sKey := os.Args[1]
	airport := os.Args[2]
	date := os.Args[3]

	client := swedavia.Client{
		URL:             "https://api.swedavia.se",
		SubscriptionKey: sKey,
	}

	arrivalsInfo, err := client.GetArrivals(airport, date)

	if err != nil {
		fmt.Println("Error getting flight info:", err)
		return
	}

	PrintArrivals(*arrivalsInfo)

	departuresInfo, err := client.GetDepartures(airport, date)

	if err != nil {
		fmt.Println("Error getting flight info:", err)
		return
	}

	PrintDepartures(*departuresInfo)
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
