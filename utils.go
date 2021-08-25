package edgarwebcrawler

import (
	"regexp"
	"strconv"
)

const (
	TRANSACTION_DATE_LAYOUT = "2006-01-02"
)

func filterSlice(vs []EdgarItem, f func(EdgarItem) bool) []EdgarItem {
	vsf := make([]EdgarItem, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func filterYearDirectories(items *EdgarDirectory, startFromYear uint16) {
	re := regexp.MustCompile(`(\d{4})`)
	items.Item = filterSlice(items.Item, func(e EdgarItem) bool {
		submatchall := re.FindAllString(e.Name, -1)
		if len(submatchall) > 0 {
			year, err := strconv.Atoi(submatchall[0])
			if err != nil {
				return false
			}
			return uint16(year) >= startFromYear
		}
		return false
	})
}
