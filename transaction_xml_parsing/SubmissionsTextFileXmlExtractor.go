package transaction_xml_parsing

import (
	"encoding/xml"
	"errors"
	"strings"
)

type XMLExtractor interface {
	ExtractXML(rawSubmissionText []byte) (*RawOwnershipDocument, error)
}

type LoadBufferToMemoryXMLExtractor struct {
	OpeningTag, ClosingTag string
}

func (extractor *LoadBufferToMemoryXMLExtractor) ExtractXML(rawSubmissionText []byte) (*RawOwnershipDocument, error) {
	submission := string(rawSubmissionText)
	xmlOpeningTagIndex, xmlClosingTagIndex, err := extractor.getOpeningAndClosingIndices(submission)
	if err != nil {
		return nil, err
	}
	xmlDocumentText := submission[xmlOpeningTagIndex : xmlClosingTagIndex+len(extractor.ClosingTag)]
	var parsedXML RawOwnershipDocument
	err = xml.Unmarshal([]byte(xmlDocumentText), &parsedXML)
	if err != nil {
		return nil, err
	}
	return &parsedXML, nil
}

func (extractor *LoadBufferToMemoryXMLExtractor) getOpeningAndClosingIndices(text string) (opening, closing int, err error) {
	opening = strings.Index(text, extractor.OpeningTag)
	if opening == -1 {
		return -1, -1, errors.New("Submission does not contain tag " + extractor.OpeningTag)
	}
	closing = strings.Index(text, extractor.ClosingTag)
	if closing == -1 {
		return -1, -1, errors.New("Submission does not contain tag " + extractor.ClosingTag)
	}
	return
}
