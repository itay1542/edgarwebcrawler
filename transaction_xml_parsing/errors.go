package transaction_xml_parsing

import "fmt"

type OfficerTypeNotFoundError struct {
	officerText string
}

func (o *OfficerTypeNotFoundError) Error() string {
	return fmt.Sprintf("no keywords found in officer text %s", o.officerText)
}
