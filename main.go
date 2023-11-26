package main

import (
	"flag"
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
	"github.com/lindsaylandry/go-mta-train-sign/src/trainfeed"
)

func main() {
	stop := flag.String("s", "D30", "stop to parse")
	direction := flag.String("d", "N", "direction of stop")
	key := flag.String("k", "foobar", "access key")
	cont := flag.Bool("c", true, "continue printing arrivals")

	flag.Parse()

	station, err := stations.GetStation(*stop)
	if err != nil {
		panic(err)
	}

	// Get subway feeds from station trains
	feeds := decoder.GetMtaFeeds(station.DaytimeRoutes)

	for {
		arrivals := []trainfeed.Arrival{}
		for _, f := range *feeds {
			t, err := trainfeed.NewTrainFeed(station, *key, *direction, f.URL)
			if err != nil {
				panic(err)
			}

			arr := t.GetArrivals()
			for _, a := range arr {
				arrivals = append(arrivals, a)
			}
		}

		// Print all arrivals
		trainfeed.PrintArrivals(arrivals, station.StopName)

		if !*cont {
			break
		}

		time.Sleep(5 * time.Second)
	}
}
