package edgarwebcrawler

import (
	"github.com/mmcdole/gofeed"
)

type SampleDiscarder interface {
	CheckSample(entry *gofeed.Item) bool
}

type InMemorySampleDiscarderById struct {
	lastSampleId string
}

func NewInMemorySampleDiscarderById() *InMemorySampleDiscarderById{
	return &InMemorySampleDiscarderById{
		lastSampleId: "",
	}
}

func (i *InMemorySampleDiscarderById) CheckSample(entry *gofeed.Item) bool {
	entryId := entry.GUID
	if i.lastSampleId == entryId{
		return false
	}
	if !i.validateEntryIsFormFour(entry){
		return false
	}
	i.lastSampleId = entryId
	return true
}

func (i *InMemorySampleDiscarderById) validateEntryIsFormFour(entry *gofeed.Item) bool {
	return entry.Categories[0] == "4"
}


