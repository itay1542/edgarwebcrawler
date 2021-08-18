package requests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Requester interface {
	Start(urls <-chan string, responseChannel chan<- []byte) error
	Stop() error
}

type EdgarRequester struct {
	ticker      *time.Ticker
	done        chan bool
	errorC      chan error
	requestRate time.Duration
}

func New(requestRate time.Duration) *EdgarRequester {
	return &EdgarRequester{
		requestRate: requestRate,
		done:        make(chan bool),
		errorC:      make(chan error),
		ticker:      nil,
	}
}

func (requester *EdgarRequester) Start(urls <-chan string, responseChannel chan<- []byte) {
	requester.ticker = time.NewTicker(requester.requestRate)
	go func() {
		for {
			select {
			case <-requester.done:
				return
			case <-requester.ticker.C:
				log.Println("Tick complete, getting response body")
				select {
				case url := <-urls:
					go requester.getResponseBody(url, responseChannel)
				case <-requester.done:
					return
				}
			}
		}
	}()
}

func (requester *EdgarRequester) Stop() error {
	requester.done <- true
	return nil
}

func (requester *EdgarRequester) getResponseBody(url string, respChannel chan<- []byte) {
	fmt.Printf("Requesting %s \n", url)
	request, err := http.NewRequest("GET", url, nil)

	request.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("got error when requesting resource %s \n", err)
		requester.errorC <- err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		requester.errorC <- err
	}
	respChannel <- data
}
