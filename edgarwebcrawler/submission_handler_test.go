package edgarwebcrawler

import (
	"testing"
	"time"
)

func TestSecSubmissionHandler_ExtractTransactions(t *testing.T) {
	layout := TRANSACTION_DATE_LAYOUT
	date := "2020-08-18"
	parsed, err := time.Parse(layout, date)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Month() != time.August || parsed.Day() != 18 || parsed.Year() != 2020 {
		t.Fatalf("time parsing failed, time received: %s", parsed)
	}
}
