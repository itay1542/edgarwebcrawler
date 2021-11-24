package edgarwebcrawler

import (
	edgarwebcrawler "github.com/itay1542/edgarwebcrawler/submission_url/rss"
	"sync"
	"testing"
	"time"
)

func TestSecSubmissionHandler_ExtractTransactions(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	discarder := edgarwebcrawler.NewInMemorySampleDiscarderById()
	edgarwebcrawler.NewUrlFromRssProvider("https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=4&company=&dateb=&owner=include&start=-1&count=100&output=rss",
		time.Second * 2, discarder).Start(make(chan string))
	wg.Wait()
}


