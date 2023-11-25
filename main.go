package main

import (
	"flag"

	"github.com/lindsaylandry/go-mta-train-sign/src/traininfo"
)

func main() {
	train := flag.String("t", "Q", "train to parse")
	key := flag.String("k", "foobar", "access key")
	flag.Parse()
		
	t, err := traininfo.NewTrainInfo(*train, *key)
	if err != nil {
		panic(err)
	}

	err = t.GetVehicleStopInfo()
	if err != nil {
		panic(err)
	}
}
