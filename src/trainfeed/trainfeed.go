package trainfeed

import (
	"fmt"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"sort"
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
)

type TrainFeed struct {
	Station   stations.MtaStation
	Key       string
	Direction string

	Feed gtfs.FeedMessage
}

type Arrival struct {
	Train string
	Mins  int64
}

func NewTrainFeed(station stations.MtaStation, accessKey, direction, url string) (*TrainFeed, error) {
	t := TrainFeed{}

	t.Key = accessKey
	t.Direction = direction

	t.Station = station

	feed, err := decoder.Decode(accessKey, url)
	t.Feed = feed

	return &t, err
}

func (t *TrainFeed) GetArrivals() []Arrival {
	stopID := t.Station.GTFSStopID + t.Direction
	now := time.Now()
	arrivals := []Arrival{}
	for _, entity := range t.Feed.Entity {
		trip := entity.GetTripUpdate()
		if trip != nil {
			stopTimes := trip.StopTimeUpdate
			for _, s := range stopTimes {
				if *s.StopId == stopID {
					route := ""
					vehicle := trip.Trip
					if vehicle != nil {
						route = *vehicle.RouteId
					}
					secs := *s.Arrival.Time - now.Unix()
					mins := secs / 60

					a := Arrival{}
					a.Train = route
					a.Mins = mins

					arrivals = append(arrivals, a)
				}
			}
		}
	}

	return arrivals
}

func PrintArrivals(arrivals []Arrival, name string) {
	fmt.Printf("STATION %s\n", name)

	if len(arrivals) == 0 {
		fmt.Println("No trains arriving at this station today")
		return
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].Mins < arrivals[j].Mins })
	for _, a := range arrivals {
		fmt.Printf("%s %d mins\n", a.Train, a.Mins)
	}
	fmt.Println()
}
