package stations

import (
	"fmt"
	"os"
	"github.com/gocarina/gocsv"
)

func GetStation(stopID string) (MtaStation, error) {
	station := MtaStation{}
	stations, err := readStations("data/nyc-subway-stations.csv")
	if err != nil {
		return station, err
	}

	// Find station, return error if not found
	for _, s := range stations {
		if s.GTFSStopID == stopID {
			return s, nil
		} 
	}

	return station, fmt.Errorf("Could not find station %s", stopID)
}

func readStations(filepath string) ([]MtaStation, error) {
	stations := []MtaStation{}
	f, err := os.Open(filepath)
	if err != nil {
		return stations, err
	}
	defer f.Close()

	if err := gocsv.UnmarshalFile(f, &stations); err != nil {
    return stations, err
  }

	return stations, nil
}
