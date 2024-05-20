package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"

	"github.com/seb-emmot/sviago/swedavia"
)

func main() {
	// setup server serving http files located at /static/

	mux := http.NewServeMux()
	mux.HandleFunc("/arrivals/{IATA}/{Date}", getArrivals)

	http.ListenAndServe(":8080", mux)
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

	var arrivals swedavia.ArrivalsInfo
	err = json.NewDecoder(file).Decode(&arrivals)

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

type Flight struct {
	From    string
	To      string
	Airline string
	SAT     string
	AAT     string
}
