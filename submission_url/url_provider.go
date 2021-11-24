package edgarwebcrawler

import (
	"bufio"
	"github.com/itay1542/edgarwebcrawler"
	"log"
	"os"
	"sync"
)

type SubmissionUrlProvider interface {
	Start(urlC chan<- string) error
	Stop()
}

type TextFileSubmissionsUrlProvider struct {
	filePath string
	quit     chan bool
}

func NewUrlProvider(filePath string) *TextFileSubmissionsUrlProvider {
	return &TextFileSubmissionsUrlProvider{
		filePath: filePath,
		quit:     make(chan bool),
	}
}

func (t *TextFileSubmissionsUrlProvider) Start(urlC chan<- string) error {
	file, err := os.Open(t.filePath)
	if err != nil {
		return err
	}

	waitgroup := sync.WaitGroup{}
	waitgroup.Add(1)
	go func() {
		fileReader := bufio.NewReader(file)
		for {
			submission, err := edgarwebcrawler.ReadLine(fileReader)
			if err != nil {
				log.Println(err)
				break
			}
			urlC <- submission
		}
		waitgroup.Done()
	}()

	return nil
}

func (t *TextFileSubmissionsUrlProvider) Stop() {
	t.quit <- true
}
