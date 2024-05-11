package swedavia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ArrivalsInfo struct {
	To              ArrivalAirport  `json:"to"`
	NumberOfFlights int             `json:"numberOfFlights"`
	Flights         []ArrivalFlight `json:"flights"`
}

type ArrivalAirport struct {
	ArrivalAirportIata string `json:"arrivalAirportIata"`
	FlightArrivalDate  string `json:"flightArrivalDate"`
}

type ArrivalFlight struct {
	FlightId                string      `json:"flightId"`
	ArrivalTime             ArrivalTime `json:"arrivalTime"`
	DepartureAirportEnglish string      `json:"departureAirportEnglish"`
	AirlineOperator         Airline     `json:"airlineOperator"`
}

type ArrivalTime struct {
	ScheduledUtc string `json:"scheduledUtc"`
	ActualUtc    string `json:"actualUtc"`
}

type Airline struct {
	Iata string `json:"iata"`
	Icao string `json:"icao"`
	Name string `json:"name"`
}

type DeparturesInfo struct {
	From            DepartureAirport  `json:"from"`
	NumberOfFlights int               `json:"numberOfFlights"`
	Flights         []DepartureFlight `json:"flights"`
}

type DepartureAirport struct {
	DepartureAirportIata string `json:"departureAirportIata"`
	FlightDepartureDate  string `json:"flightDepartureDate"`
}

type DepartureFlight struct {
	FlightId              string        `json:"flightId"`
	DepartureTime         DepartureTime `json:"departureTime"`
	ArrivalAirportEnglish string        `json:"arrivalAirportEnglish"`
	AirlineOperator       Airline       `json:"airlineOperator"`
}

type DepartureTime struct {
	ScheduledUtc string `json:"scheduledUtc"`
	ActualUtc    string `json:"actualUtc"`
}

type FlightInfoLoader interface {
	GetArrivals(airport, date string) (*ArrivalsInfo, error)
	GetDepartures(airport, date string) (*DeparturesInfo, error)
}

type Client struct {
	// URL of the API
	URL string
	// Subscription key
	SubscriptionKey string
}

// GetArrivals makes an HTTP GET request to the API and returns the response body
func (c *Client) GetArrivals(airport, date string) (*ArrivalsInfo, error) {

	url := fmt.Sprintf("%s/flightinfo/v2/%s/arrivals/%s", c.URL, airport, date)
	data, err := getInfo(url, c.SubscriptionKey)

	if err != nil {
		fmt.Println("Error calling:", c.URL, err)
	}

	// parse body into json
	var arrivals ArrivalsInfo
	err = json.Unmarshal(data, &arrivals)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return &arrivals, nil
}

func (c *Client) GetDepartures(airport, date string) (*DeparturesInfo, error) {

	url := fmt.Sprintf("%s/flightinfo/v2/%s/departures/%s", c.URL, airport, date)
	data, err := getInfo(url, c.SubscriptionKey)

	if err != nil {
		fmt.Println("Error calling:", c.URL, err)
	}

	// parse body into json
	var departures DeparturesInfo
	err = json.Unmarshal(data, &departures)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return &departures, nil
}

func getInfo(url, subKey string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		url, nil)

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	// Configure headers for Swedavia api.
	req.Header.Set("Ocp-Apim-Subscription-Key", subKey)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/json")

	// Make the HTTP request
	res, err := http.DefaultClient.Do(req)

	fmt.Println("request", req)

	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}

	defer res.Body.Close()

	// Check the status code
	if res.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", res.Status)
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}
