# Sviago project.
In Development

Sviago is an aggregator of data for Swedish airports operated by Swedavia.

## data collector
The data collector fetches data and saves/prints it.
To run the data collector.\
`go run cmd/fetch/fetch.go <sub key> <airportIATA> <date> [filename]`

## server
To run the server
- set env variable `SWEDAVIA_SUBSCRIPTION_KEY`\
- run `go run cmd/server/server.go`
- visit `localhost:8080/arrivals/{iata}/{date}`
