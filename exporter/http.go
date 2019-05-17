package exporter

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func asyncHTTPGets(targets []string, user string, pass string) ([]*Response, error) {

	// Channels used to enable concurrent requests
	ch := make(chan *Response, len(targets))

	responses := []*Response{}

	for _, url := range targets {

		go func(url string) {
			err := getResponse(url, user, pass, ch)
			if err != nil {
				ch <- &Response{url, nil, []byte{}, err}
			}
		}(url)

	}

	for {
		select {
		case r := <-ch:
			if r.err != nil {
				log.Errorf("Error scraping API, Error: %v", r.err)
				break
			}
			responses = append(responses, r)

			if len(responses) == len(targets) {
				return responses, nil
			}
		}

	}
}

// getResponse collects an individual http.response and returns a *Response
func getResponse(url string, user string, pass string, ch chan<- *Response) error {

	log.Infof("Fetching %s \n", url)

	resp, err := getHTTPResponse(url, user, pass) // do this earlier

	if err != nil {
		return fmt.Errorf("Error converting body to byte array: %v", err)
	}

	// Read the body to a byte array so it can be used elsewhere
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("Error converting body to byte array: %v", err)
	}

	if resp.StatusCode == 404 {
		return fmt.Errorf("Error: Received 404 status from Sonus API, ensure the URL is correct. ")
	}

	ch <- &Response{url, resp, body, err}

	return nil
}

// getHTTPResponse handles the http client creation, token setting and returns the *http.response
func getHTTPResponse(url string, user string, pass string) (*http.Response, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}
