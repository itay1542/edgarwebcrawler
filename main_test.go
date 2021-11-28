package edgarwebcrawler

import (
	"fmt"
	edgarwebcrawler "github.com/itay1542/edgarwebcrawler/submission_url/rss"
	"log"
	"sync"
	"testing"
	"time"
)

func TestSecSubmissionHandler_ExtractTransactions(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var entryCacheSize uint = 40
	discarder := edgarwebcrawler.NewInMemorySampleDiscarderById(entryCacheSize)
	channel := make(chan string)
	edgarwebcrawler.NewUrlFromRssProvider(fmt.Sprintf("https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=4&company=&dateb=&owner=include&start=-1&count=%d&output=rss", entryCacheSize),
		time.Second*2, discarder).Start(channel)
	for {
		select {
		case c := <-channel:
			log.Println(c)
		}
	}
	wg.Wait()
}
