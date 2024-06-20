package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/seb-emmot/sviago/swedavia"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/arrivals/{IATA}/", getArrivalDist)
	mux.HandleFunc("/arrivals/{IATA}/{Date}", getArrivals)

	http.ListenAndServe(":8080", mux)
}

func getArrivalDist(w http.ResponseWriter, r *http.Request) {
	// list all files in directory
	files, err := os.ReadDir("data/")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	toUse := make([]string, 0)

	for _, file := range files {
		// check of filename contains IATA
		if file.IsDir() {
			continue
		}

		fpath := filepath.Join("data/", file.Name())
		iata := r.PathValue("IATA")

		if strings.Contains(fpath, iata) {
			toUse = append(toUse, fpath)
		}
	}

	allArrivals := make([]swedavia.ArrivalsInfo, 0)

	for _, fname := range toUse {
		log.Println("Reading file", fname)
		file, err := os.Open(fname)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
			return
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

	tmpl, err := template.ParseFiles("static/dist.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	log.Println(hourDist)
	err = tmpl.Execute(w, hourDist)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}

func getArrivals(w http.ResponseWriter, r *http.Request) {
	iata := r.PathValue("IATA")
	date := r.PathValue("Date")

	log.Println(r.URL.Path)

	if iata == "" || date == "" {
		http.Error(w, "IATA and Date must be provided", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	fname := fmt.Sprintf("data/arrivals_%s_%s.json", iata, date)

	var file *os.File
	file, err = os.Open(fname)

	if err != nil {
		// if data does not exist, create it.
		if os.IsNotExist(err) {
			// read environment variable
			sKey, ok := os.LookupEnv("SWEDAVIA_SUBSCRIPTION_KEY")
			if !ok {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Fatal("SWEDAVIA_SUBSCRIPTION_KEY not set")
				return
			}

			client := swedavia.Client{
				URL:             "https://api.swedavia.se",
				SubscriptionKey: sKey,
			}

			arrivalsInfo, err := client.GetArrivals(iata, date)
			fmt.Println("polled arrival data")

			if err != nil {
				http.Error(w, "Error fetching data", http.StatusInternalServerError)
				fmt.Println("Error fetching arrivals:", err)
				return
			}

			file, err := os.Create(fname)
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
			fmt.Printf("Arrivals info written to %s", fname)
		} else {
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			log.Fatal(err)
		}
		file, err = os.Open(fname)

		if err != nil {
			log.Fatal(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	} else {
		defer file.Close()
	}

	arrivals := parseArrivals(file)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	sort.Slice(arrivals.Flights, func(i, j int) bool {
		return arrivals.Flights[i].AirlineOperator.Name < arrivals.Flights[j].AirlineOperator.Name
	})

	err = tmpl.Execute(w, arrivals)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func parseArrivals(r io.Reader) *swedavia.ArrivalsInfo {
	var arrivals swedavia.ArrivalsInfo
	err := json.NewDecoder(r).Decode(&arrivals)

	if err != nil {
		log.Fatal(err)
	}

	return &arrivals
}

type Flight struct {
	From    string
	To      string
	Airline string
	SAT     string
	AAT     string
}

func arrivalHours(info swedavia.ArrivalsInfo) map[int]int {
	m := make(map[int]int)
	for _, flight := range info.Flights {
		// get hour of arrival
		hr := parseArrivalHour(flight)
		m[hr]++
	}

	return m
}

func parseArrivalHour(fl swedavia.ArrivalFlight) int {
	t, err := time.Parse("2006-01-02T15:04:05Z", fl.ArrivalTime.ScheduledUtc)
	if err != nil {
		log.Fatal(err)
	}
	return t.Hour()
}
