package edgarwebcrawler

type StockExchange string

const (
	NYSE   StockExchange = "NYSE"
	NASDAQ               = "NASDAQ"
)

type OfficerType struct {
	Id int
	OfficialTitle string
	TitleVariations []string
}

type CompanyType struct {
	
}

type StockExchangeType struct {
	Id int
	Symbol
}