package exporter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Response struct is used to store http.Response and associated data
type Response struct {
	url      string
	response *http.Response
	body     *[]byte
}

// doHTTPRequest makes an individual HTTP request and returns a *Response
func doHTTPRequest(client *http.Client, url string, user string, pass string) (*Response, error) {
	log.Infof("Fetching %s \n", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Read the body to a byte array so it can be used elsewhere
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Error: Received 404 status from Sonus API, ensure the URL is correct. ")
	}

	return &Response{url, resp, &body}, nil
}
