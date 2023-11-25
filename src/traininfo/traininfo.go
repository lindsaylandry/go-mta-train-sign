package traininfo

import (
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"fmt"

	"github.com/lindsaylandry/go-mta-train-sign/src/decoder"
)

type TrainInfo struct {
	Train string
	Key string

	Feed gtfs.FeedMessage
}

func NewTrainInfo(train, accessKey string) (*TrainInfo, error) {
	t := TrainInfo{}

	t.Key = accessKey
	t.Train = train

	feed, err := decoder.Decode(accessKey)
	t.Feed = feed

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
