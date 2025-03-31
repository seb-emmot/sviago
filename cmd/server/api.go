package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/seb-emmot/sviago/swedavia"
)

func getHourDist(iata string) (map[int]int, error) {
	// list all files in directory
	files, err := os.ReadDir("../../data/")
	if err != nil {
		return nil, err
	}

	toUse := make([]string, 0)

	for _, file := range files {
		// check of filename contains IATA
		if file.IsDir() {
			continue
		}

		fpath := filepath.Join("../../data/", file.Name())

		if strings.Contains(fpath, iata) {
			toUse = append(toUse, fpath)
		}
	}

	allArrivals := make([]swedavia.ArrivalsInfo, 0)

	for _, fname := range toUse {
		log.Println("Reading file", fname)
		file, err := os.Open(fname)

		if err != nil {
			return nil, err
		}

		arrivals := parseArrivals(file)
		allArrivals = append(allArrivals, *arrivals)
	}

	hourDist := make(map[int]int, 0)

	// intialize all hours to 0 flights.
	for i := 0; i < 24; i++ {
		hourDist[i] = 0
	}

	for _, arrival := range allArrivals {
		hourMap := arrivalHours(arrival)
		for k, v := range hourMap {
			hourDist[k] += v
		}
	}

	log.Println(hourDist)

	return hourDist, nil
}

func parseArrivals(r io.Reader) *swedavia.ArrivalsInfo {
	var arrivals swedavia.ArrivalsInfo
	err := json.NewDecoder(r).Decode(&arrivals)

	if err != nil {
		log.Fatal(err)
	}

	return &arrivals
}

func getArrivalEntries(iataFilter string) ([]ArrivalEntry, error) {
	files, err := os.ReadDir("../../data/")
	if err != nil {
		return nil, err
	}

	// get list of files

	arrivalEntries := make([]ArrivalEntry, 0)

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		// assuming file is called arrivals_IATA_DATE.json
		parts := strings.Split(file.Name(), "_")

		iata := parts[1]

		if iata != iataFilter {
			continue
		}

		d := strings.Trim(parts[2], ".json")

		entry := ArrivalEntry{
			Iata: iata,
			Date: d,
			Link: fmt.Sprintf("/arrivals/%s/%s", iata, d)}

		arrivalEntries = append(arrivalEntries, entry)

		fmt.Printf("found file %s\n", file.Name())
	}

	return arrivalEntries, nil
}

func getArrivalInfo(iata, date string) (*swedavia.ArrivalsInfo, error) {
	fname := fmt.Sprintf("../../data/arrivals_%s_%s.json", iata, date)

	var file *os.File
	file, err := os.Open(fname)

	if err != nil {
		// if data does not exist, create it.
		if os.IsNotExist(err) {
			// read environment variable
			sKey, ok := os.LookupEnv("SWEDAVIA_SUBSCRIPTION_KEY")
			if !ok {
				return nil, fmt.Errorf("Missing Key")
			}

			client := swedavia.Client{
				URL:             "https://api.swedavia.se",
				SubscriptionKey: sKey,
			}

			arrivalsInfo, err := client.GetArrivals(iata, date)
			fmt.Println("polled arrival data")

			if err != nil {
				fmt.Println("Error fetching arrivals:", err)
				return nil, err
			}

			file, err := os.Create(fname)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return nil, err
			}

			defer file.Close()

			encoder := json.NewEncoder(file)
			err = encoder.Encode(arrivalsInfo)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return nil, err
			}
			fmt.Printf("Arrivals info written to %s", fname)
		} else {
			return nil, err
		}
		file, err = os.Open(fname)

		if err != nil {
			return nil, err
		}
	} else {
		defer file.Close()
	}

	arrivals := parseArrivals(file)

	return arrivals, nil
}
