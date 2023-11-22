package traininfo

import (
	"google.golang.org/protobuf/proto"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"net/http"
	"io/ioutil"
	"fmt"
	"errors"
)

type TrainInfo struct {
	Train string
	Key string
}

func NewTrainInfo(train, accessKey string) (*TrainInfo, error) {
	t := TrainInfo{}

	t.Key = accessKey
	t.Train = train

	return &t, nil
}

func (t *TrainInfo) Decode() (error) {
	client := &http.Client{}
  req, err := http.NewRequest("GET", "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw", nil)
  req.Header.Add("x-api-key", t.Key)
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    return err
  }

	// read response code
	// TODO: make more robust
	if resp.StatusCode == 401 {
		return errors.New("Unauthorized")
	}

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  feed := gtfs.FeedMessage{}
  err = proto.Unmarshal(body, &feed)
  if err != nil {
    return err
  }

	for _, entity := range feed.Entity {
		vehicle := entity.GetVehicle()
		if vehicle != nil {
			stopId := vehicle.StopId
			fmt.Printf("Stop ID: %s\n", *stopId)
		}
  }

	return nil
}
