package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/traininfo"
)

func main() {
	stop := flag.String("s", "D30N", "stop to parse")
	key := flag.String("k", "foobar", "access key")
	cont := flag.Bool("c", true, "continue printing arrivals")

	flag.Parse()

	for {
		t, err := traininfo.NewTrainInfo(*stop, *key)
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
