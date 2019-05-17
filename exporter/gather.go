package exporter

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
)

// gatherData - Collects the data from the API and stores into struct
func (e *Exporter) gatherData() ([]*TGcollection, error) {

	data := []*TGcollection{}

	responses, err := asyncHTTPGets(e.TargetURLs, e.APIUser, e.APIPass)

	if err != nil {
		return data, err
	}

	for _, response := range responses {

		d := new(TGcollection)
		xml.Unmarshal(response.body, &d)
		data = append(data, d)

		log.Infof("API data fetched: %s", response.url)
	}

	//return data, rates, err
	return data, nil

}
