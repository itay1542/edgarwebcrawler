package edgarwebcrawler

import (
	"github.com/itay1542/edgarwebcrawler/requests"
	"github.com/itay1542/edgarwebcrawler/transaction_xml_parsing"
	"log"
)

type Orchestrator interface {
	Run() error
}

type ReadyFileOrchestrator struct {
	urls           chan string
	submissions    chan []byte
	running        bool
	urlProvider    requests.SubmissionUrlProvider
	edgarRequester requests.EdgarRequester
	xmlExtractor   transaction_xml_parsing.XMLExtractor
}

func NewOrchestrator(urlBufferSize uint,
	urlProvider requests.SubmissionUrlProvider,
	edgarRequester requests.EdgarRequester,
	xmlExtractor transaction_xml_parsing.XMLExtractor) *ReadyFileOrchestrator {
	return &ReadyFileOrchestrator{
		urls:           make(chan string, urlBufferSize),
		submissions:    make(chan []byte),
		running:        false,
		urlProvider:    urlProvider,
		edgarRequester: edgarRequester,
		xmlExtractor:   xmlExtractor,
	}
}

func (o *ReadyFileOrchestrator) Run() error {
	err := o.urlProvider.Start(o.urls)
	if err != nil {
		return err
	}
	o.edgarRequester.Start(o.urls, o.submissions)
	for {
		select {
		case submission := <-o.submissions:
			parsed, err := o.xmlExtractor.ExtractXML(submission)
			if err != nil {
				log.Print(err)
			}
			log.Println(parsed.ReportingOwner.ReportingOwnerRelationship)
		}
	}
}
