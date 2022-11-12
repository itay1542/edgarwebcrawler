package edgarwebcrawler

import (
	"github.com/mmcdole/gofeed"
	"github.com/thoas/go-funk"
	"log"
	"sync"
)

type SampleDiscarder interface {
	CheckSample(entry *gofeed.Item) bool
}

type InMemorySampleDiscarderById struct {
	seenSamples []string
	mutex       sync.Mutex
}

func NewInMemorySampleDiscarderById(sampleCacheSize uint) *InMemorySampleDiscarderById {
	return &InMemorySampleDiscarderById{
		seenSamples: make([]string, sampleCacheSize),
		mutex:       sync.Mutex{},
	}
}

func (i *InMemorySampleDiscarderById) CheckSample(entry *gofeed.Item) bool {
	entryId := entry.GUID
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if funk.ContainsString(i.seenSamples, entryId) {
		return false
	}
	if !i.validateEntryIsFormFour(entry) {
		log.Println("entry found to not be form 4")
		return false
	}
	i.seenSamples = i.seenSamples[1:]
	i.seenSamples = append(i.seenSamples, entryId)
	return true
}

func (i *InMemorySampleDiscarderById) validateEntryIsFormFour(entry *gofeed.Item) bool {
	return entry.Categories[0] == "4"
}
