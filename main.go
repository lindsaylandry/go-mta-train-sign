package main

import (
	"flag"

	"github.com/lindsaylandry/go-mta-train-sign/src/traininfo"
)

func main() {
	stop := flag.String("s", "D30N", "stop to parse")
	key := flag.String("k", "foobar", "access key")
	
	flag.Parse()
		
	t, err := traininfo.NewTrainInfo(*stop, *key)
	if err != nil {
		panic(err)
	}

	//err = t.GetVehicleStopInfo()
	//if err != nil {
	//	panic(err)
	//}

	err = t.GetTripUpdateInfo(*stop)
	if err != nil {
		panic(err)
	}
}
