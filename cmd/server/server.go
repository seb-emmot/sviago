package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/seb-emmot/sviago/swedavia"
)

type ArrivalEntries struct {
	Entries []ArrivalEntry
}

type ArrivalEntry struct {
	Iata string
	Date string
	Link string
}

func main() {
	log.Println("Starting server.")
	mux := http.NewServeMux()

	mux.HandleFunc("/", getIndex)
	mux.HandleFunc("/arrivals/{IATA}/", getArrivalsAirport)
	mux.HandleFunc("/arrivals/{IATA}/{Date}", getArrivalsAirportDate)
	mux.HandleFunc("/arrivals/{IATA}/distributions", getArrivalDist)

	http.ListenAndServe(":8080", mux)
}

func getIndex(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("../../static/index.html")

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	data := struct {
		IATA []string
	}{
		IATA: []string{"GOT", "ARN"},
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func getArrivalsAirport(w http.ResponseWriter, r *http.Request) {
	iata := r.PathValue("IATA")

	log.Println(r.URL.Path)

	if iata == "" {
		http.Error(w, "IATA must be provided", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("../../static/arrivalairport.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	arrivalEntries, err := getArrivalEntries(iata)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	err = tmpl.Execute(w, ArrivalEntries{Entries: arrivalEntries})

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}

func getArrivalDist(w http.ResponseWriter, r *http.Request) {
	iata := r.PathValue("IATA")

	if iata == "" {
		http.Error(w, "IATA must be provided", http.StatusBadRequest)
		return
	}

	hourDist, err := getHourDist(iata)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	tmpl, err := template.ParseFiles("../../static/dist.html")
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

func getArrivalsAirportDate(w http.ResponseWriter, r *http.Request) {
	iata := r.PathValue("IATA")
	date := r.PathValue("Date")

	log.Println(r.URL.Path)

	if iata == "" || date == "" {
		http.Error(w, "IATA and Date must be provided", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("../../static/arrivalairportdate.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	arrivals, err := getArrivalInfo(iata, date)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Print(err)
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

func arrivalHours(info swedavia.ArrivalsInfo) map[int]int {
	m := make(map[int]int)
	for _, flight := range info.Flights {
		// get hour of arrival
		t, err := time.Parse("2006-01-02T15:04:05Z", flight.ArrivalTime.ScheduledUtc)

		if err != nil {
			log.Fatal(err)
		}

		hr := t.Hour()
		m[hr]++
	}

	return m
}
