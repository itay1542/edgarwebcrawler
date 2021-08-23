package transaction_xml_parsing

import (
	"strings"
)

type OfficerType string

const (
	CEO OfficerType = "CEO"
	CFO OfficerType = "CFO"
	COO OfficerType = "COO"
	CMO OfficerType = "CMO"
	CTO OfficerType = "CTO"
	CAO OfficerType = "CAO"
	COB OfficerType = "COB"
)

var (
	OFFICER_KEY_TOKENS = map[OfficerType][]string{
		CFO: []string{"financial",
			"cfo",
		},
		CEO: []string{
			"executive",
			"ceo",
		},
		COO: []string{
			"operating",
			"coo",
		},
		CMO: []string{
			"marketing",
			"cmo",
		},
		CTO: []string{
			"technology",
			"cto",
			"technical",
		},
		CAO: []string{
			"accounting",
			"cao",
		},
		COB: []string{
			"cob",
			"chairman of the board",
		},
	}
)

type OfficerClassifier interface {
	GetOfficerType(officerText string) OfficerType
}

//KeyTokensOfficerClassifier implements OfficerClassifier
type KeyTokensOfficerClassifier struct {
}

func (k *KeyTokensOfficerClassifier) GetOfficerType(officerText string) ([]OfficerType, error) {
	officerTypes := make([]OfficerType, 0, 1)
	for key, val := range OFFICER_KEY_TOKENS {
		for _, token := range val {
			if strings.Contains(strings.ToLower(officerText), token) {
				officerTypes = append(officerTypes, key)
			}
			break
		}
	}
	if len(officerTypes) == 0 {
		return officerTypes, &OfficerTypeNotFoundError{officerText: officerText}
	}
	return officerTypes, nil
}
