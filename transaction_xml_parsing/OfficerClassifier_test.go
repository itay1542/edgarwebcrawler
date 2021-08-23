package transaction_xml_parsing

import (
	"testing"
)

func TestKeyTokensOfficerClassifier_GetOfficerType(t *testing.T) {
	t.Run("should return CEO for chief executive officer text", func(test *testing.T) {
		officerText := "chief executive officer"
		expected := CEO
		got, _ := (&KeyTokensOfficerClassifier{}).GetOfficerType(officerText)
		if got[0] != expected {
			test.Fatalf("received %s instead of %s", got[0], expected)
		}
	})

	t.Run("should ignore capital letters and return cfo as result", func(test *testing.T) {
		officerText := "Chief Financial Officer and some other thing"
		expected := CFO
		got, _ := (&KeyTokensOfficerClassifier{}).GetOfficerType(officerText)
		if got[0] != expected {
			test.Fatalf("received %s instead of %s", got[0], expected)
		}
	})

	t.Run("Should Return a slice containing both COB and CEO when given two roles", func(test *testing.T) {
		officerText := "COB and CEO"
		got, _ := (&KeyTokensOfficerClassifier{}).GetOfficerType(officerText)
		for _, val := range got {
			if val != COB && val != CEO {
				test.Fatalf("CEO or COB were not returned")
			}
		}
	})
}
