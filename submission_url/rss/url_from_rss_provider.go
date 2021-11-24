package edgarwebcrawler

import (
	"github.com/mmcdole/gofeed"
	"log"
	"strings"
	"time"
)

type UrlFromRssProvider struct {
	rssUrl          string
	quit            chan bool
	parser          *gofeed.Parser
	sampleFreq      time.Duration
	sampleDiscarder SampleDiscarder
}

func NewUrlFromRssProvider(rssUrl string, sampleFreq time.Duration, sampleDiscarder SampleDiscarder) *UrlFromRssProvider {
	return &UrlFromRssProvider{
		rssUrl:     rssUrl,
		quit:       make(chan bool),
		parser:     gofeed.NewParser(),
		sampleFreq: sampleFreq,
		sampleDiscarder: sampleDiscarder,
	}
}

func (t *UrlFromRssProvider) Start(urlC chan<- string) error {
	ticker := time.NewTicker(t.sampleFreq)
	go func() {
		for {
			select {
			case <-t.quit:
				ticker.Stop()
				break
			case <-ticker.C:
				log.Println("starting sampling rss")
				feed, err := t.parser.ParseURL(t.rssUrl)
				if err != nil {
					log.Println(err)
				}
				t.processItems(feed.Items, urlC)
			}
		}
	}()

	return nil
}

func (t *UrlFromRssProvider) Stop() {
	log.Println("stopping rss sampling")
	t.quit <- true
}

func (t *UrlFromRssProvider) processItems(items []*gofeed.Item, urlC chan<- string) {
	for _, item := range items{
		if !t.sampleDiscarder.CheckSample(item){
			log.Println("sample already seen")
			continue
		}
		log.Println("new item detected")
		entryLink := t.getTextEntry(item.Link)
		log.Println("sending link")
		urlC <- entryLink
	}
}

func (t *UrlFromRssProvider) getTextEntry(htmlLink string) string {
	return strings.Replace(htmlLink, "-index.htm", ".txt", 1)
}