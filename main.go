package main

import (
	"flag"
	"time"
	"fmt"

	"github.com/lindsaylandry/go-mta-train-sign/src/traininfo"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
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

	for {
		t, err := traininfo.NewTrainInfo(station, *key, *direction)
		if err != nil {
			panic(err)
		}

		t.PrintArrivals()

		if !*cont {
			break
		}
		
		fmt.Println()
		time.Sleep(5 * time.Second)
	}
}
