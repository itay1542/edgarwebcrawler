package edgarwebcrawler

type TransactionParsed struct {

}

type TransactionFetcher interface {
	FetchTransaction(textFileUrl string) TransactionParsed
}
