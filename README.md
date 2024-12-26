# Sviago project.
In Development

Sviago is an aggregator of data for Swedish airports operated by Swedavia.

For running this you can acquire a subscription key from Swedavia dev portal.
https://apideveloper.swedavia.se/


## data collector
The data collector fetches data and saves/prints it.
To run the data collector.\
`go run cmd/fetch/fetch.go <airportIATA> <date> [outputdir]`

## server
To run the server
- set env variable `SWEDAVIA_SUBSCRIPTION_KEY` this should be your FlightInfo key
- run `go run cmd/server/server.go`
- visit `localhost:8080/arrivals/{iata}/{date}`
