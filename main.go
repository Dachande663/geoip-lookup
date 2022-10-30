package main

import (
	"github.com/joho/godotenv"
	"github.com/oschwald/maxminddb-golang"
	"log"
	"net/http"
	"os"
	"time"
)

var startTime time.Time
var config map[string]string
var ipDatabase *maxminddb.Reader
var totalRequests uint64
var totalLookups uint64
var totalLookupsFound uint64

func main() {
	startTime = time.Now()
	var err error
	config, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	for _, requiredConfigKey := range []string{
		"HTTP_ADDR",
		"MAXMIND_CITY_FILE",
	} {
		if val, ok := config[requiredConfigKey]; ok {
			if val == "" {
				log.Fatal("Required config key is empty in env: " + requiredConfigKey)
			}
		} else {
			log.Fatal("Required config key not found in env: " + requiredConfigKey)
		}
	}

	if _, err := os.Stat(config["MAXMIND_CITY_FILE"]); err != nil {
		log.Fatal("Unable to find maxmind city db file: " + config["MAXMIND_CITY_FILE"])
	}
	ipDatabase, err = maxminddb.Open(config["MAXMIND_CITY_FILE"])
	if err != nil {
		log.Fatal(err)
	}
	defer ipDatabase.Close()

	http.Handle("/", appMiddleware(http.HandlerFunc(indexRoute)))
	http.Handle("/status", appMiddleware(http.HandlerFunc(statusRoute)))
	http.Handle("/lookup/", appMiddleware(http.HandlerFunc(lookupRoute)))
	log.Println("Listening on " + config["HTTP_ADDR"])
	http.ListenAndServe(config["HTTP_ADDR"], nil)
}
