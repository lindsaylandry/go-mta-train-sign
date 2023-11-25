package traininfo

import (
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"fmt"
	"sort"
	"time"

	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
	"github.com/lindsaylandry/go-mta-train-sign/src/stations"
)

type TrainInfo struct {
	Station stations.MtaStation
	Key string
	Direction string

	Feed gtfs.FeedMessage

	Arrivals []Arrival
}

type Arrival struct {
  Train string
  Mins int64
}

func NewTrainInfo(station stations.MtaStation, accessKey, direction, url string) (*TrainInfo, error) {
	t := TrainInfo{}

	t.Key = accessKey
	t.Direction = direction

	t.Station = station

	feed, err := decoder.Decode(accessKey, url)
	t.Feed = feed

	t.GetArrivals()

	return &t, err
}

func (t *TrainInfo) GetVehicleStopInfo() (error) {
	for _, entity := range t.Feed.Entity {
		vehicle := entity.GetVehicle()
		if vehicle != nil {
			stopId := vehicle.StopId
			msg := fmt.Sprintf("Stop ID: %s", *stopId)
			status := vehicle.CurrentStatus
			if status != nil {
				msg = fmt.Sprintf("%s %s", msg, *status)
			}	else {
				currentStatus := vehicle.CurrentStatus
				if currentStatus == nil {
					msg = fmt.Sprintf("%s IN_TRANSIT_TO", msg)
				} else {
					switch *currentStatus {
					case gtfs.VehiclePosition_INCOMING_AT:
						msg = fmt.Sprintf("%s INCOMING_AT", msg)
					case gtfs.VehiclePosition_STOPPED_AT:
						msg = fmt.Sprintf("%s STOPPED_AT", msg)
					case gtfs.VehiclePosition_IN_TRANSIT_TO:
						msg = fmt.Sprintf("%s IN_TRANSIT_TO", msg)
					// If current_status is missing IN_TRANSIT_TO is assumed.
					default:
						msg = fmt.Sprintf("%s IN_TRANSIT_TO", msg)
					}
				}
			}
			desc := vehicle.Vehicle
			if desc != nil {
				msg = fmt.Sprintf("Label: %s %s", *desc.Label, msg)
			}
			fmt.Printf("%s\n", msg)
		}
	}
	return nil
}

func (t *TrainInfo) GetArrivals() () {
	stopID := t.Station.GTFSStopID + t.Direction
	now := time.Now()
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
					mins := secs/60

					a := Arrival{}
					a.Train = route
					a.Mins = mins

					t.Arrivals = append(t.Arrivals, a)
				}
			}
		}
  }

	// TODO: sort arrivals by mins
	sort.Slice(t.Arrivals, func(i, j int) bool { return t.Arrivals[i].Mins < t.Arrivals[j].Mins })
}

func (t *TrainInfo) PrintArrivals() {
	fmt.Printf("STATION %s\n", t.Station.StopName)
	for _, a := range t.Arrivals {
		fmt.Printf("%s %d mins\n", a.Train, a.Mins)
	}
}
