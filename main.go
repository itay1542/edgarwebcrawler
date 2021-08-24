package main

import (
	"fmt"
	"github.com/itay1542/edgarwebcrawler/DAL"
	"github.com/itay1542/edgarwebcrawler/edgarwebcrawler"
	"github.com/itay1542/edgarwebcrawler/requests"
	"github.com/itay1542/edgarwebcrawler/transaction_xml_parsing"
	"github.com/itay1542/edgarwebcrawler/transactions"
	"github.com/itay1542/edgarwebcrawler/utils"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func main() {
	var cfg Config
	utils.ReadConfigFile(&cfg, "config.yml")
	fmt.Printf("loaded configuration: %+v", cfg)
	program := initDependencies(&cfg)
	err := program.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func saveFiles(parser edgarwebcrawler.IdxReader, destFile, addPrefix string) error {
	var row *edgarwebcrawler.IdxRow
	//skip the annoying ------ row
	_, err := parser.ReadRow()
	hasSeenFormFour := false
	f, err := os.OpenFile(destFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	for {
		row, err = parser.ReadRow()
		if err != nil {
			log.Fatal(err)
		}
		if hasSeenFormFour && row.FormType != "4" {
			break
		}
		if row.FormType == "4" {
			hasSeenFormFour = true
			if _, err := f.WriteString(fmt.Sprintf("%s\n", addPrefix+row.FileName)); err == nil {
				log.Printf("Successfully appended file name")
			}

		}
		fmt.Printf("%+v \n", row)
	}
	return nil
}

func initDependencies(config *Config) edgarwebcrawler.Orchestrator {
	dbConfig := DAL.DBConfiguration{
		Host:     config.DB.Host,
		User:     config.DB.Username,
		Password: config.DB.Password,
		DBName:   config.DB.Name,
		Port:     config.DB.Port,
	}
	companyGetter := requests.NewAlphaVantageRequester(config.AlphaVantage.Host, config.AlphaVantage.ApiKey)
	dal := &DAL.PostgresInsideOut{Config: dbConfig}
	officerClassifier := &transaction_xml_parsing.KeyTokensOfficerClassifier{}
	commonStockTypeFilter := &transactions.CommonStockTypeTransactionFilter{
		TargetString: "Common Stock",
	}
	stockExchangeTypeFilter := transactions.NewStockExchangeTypeFilter(config.Filter.StockExchanges, companyGetter)
	filters := []transactions.TransactionFilterer{
		stockExchangeTypeFilter, commonStockTypeFilter,
	}
	submissionHandler := edgarwebcrawler.NewSecSubmissionHandler(dal, filters, officerClassifier, companyGetter)
	urlProvider := requests.NewUrlProvider("storage\\form4_submission_uris_2016-.txt")
	edgarRequester := requests.New(2 * time.Second)
	xmlParser := &transaction_xml_parsing.LoadBufferToMemoryXMLExtractor{
		OpeningTag: "<ownershipDocument>",
		ClosingTag: "</ownershipDocument>",
	}
	err := dal.Init()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	orchestrator := edgarwebcrawler.NewOrchestrator(2000, urlProvider, *edgarRequester, xmlParser, submissionHandler)
	return orchestrator
}
