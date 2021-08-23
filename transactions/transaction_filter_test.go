package transactions

import (
	"encoding/xml"
	edgarwebcrawler "github.com/itay1542/edgarwebcrawler/DAL"
	"github.com/itay1542/edgarwebcrawler/requests"
	"github.com/itay1542/edgarwebcrawler/transaction_xml_parsing"
	"testing"
)

func TestShouldKeepBothStocks(t *testing.T) {
	instance := &CommonStockTypeTransactionFilter{lookForString: "common stock"}
	transactions := make([]transaction_xml_parsing.NonDerivativeTransaction, 2)
	transactions[1] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Common Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature:        transaction_xml_parsing.OwnerShipNature{},
	}
	transactions[0] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Class A Common Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature:        transaction_xml_parsing.OwnerShipNature{},
	}
	rawTransaction := &transaction_xml_parsing.RawOwnershipDocument{
		XMLName:        xml.Name{},
		Issuer:         transaction_xml_parsing.Issuer{},
		ReportingOwner: transaction_xml_parsing.ReportingOwner{},
		NonDerivativeTable: &transaction_xml_parsing.NonDerivativeTable{
			XMLName:      xml.Name{},
			Transactions: &transactions,
		},
	}
	expectArrayLength := len(transactions)
	got, _ := instance.ShouldKeep(rawTransaction)
	if got != true {
		t.Fatalf("Expected to get true return type but got %t", got)
	}
	if len(*rawTransaction.NonDerivativeTable.Transactions) != expectArrayLength {
		t.Fatalf("Expected transaction array length %d, received %d", expectArrayLength,
			len(*rawTransaction.NonDerivativeTable.Transactions))
	}
}

func TestShouldKeepOnlyOneStockAndReturnTrue(t *testing.T) {
	instance := &CommonStockTypeTransactionFilter{lookForString: "common stock"}
	transactions := make([]transaction_xml_parsing.NonDerivativeTransaction, 2)
	transactions[1] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Common Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature:        transaction_xml_parsing.OwnerShipNature{},
	}
	transactions[0] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Preferred Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature:        transaction_xml_parsing.OwnerShipNature{},
	}
	rawTransaction := &transaction_xml_parsing.RawOwnershipDocument{
		XMLName:        xml.Name{},
		Issuer:         transaction_xml_parsing.Issuer{},
		ReportingOwner: transaction_xml_parsing.ReportingOwner{},
		NonDerivativeTable: &transaction_xml_parsing.NonDerivativeTable{
			XMLName:      xml.Name{},
			Transactions: &transactions,
		},
	}
	expectArrayLength := len(transactions) - 1
	got, _ := instance.ShouldKeep(rawTransaction)
	if got != true {
		t.Fatalf("Expected to get true return type but got %t", got)
	}
	if len(*rawTransaction.NonDerivativeTable.Transactions) != expectArrayLength {
		t.Fatalf("Expected transaction array length %d, received %d", expectArrayLength,
			len(*rawTransaction.NonDerivativeTable.Transactions))
	}
}

func TestShouldEmptyTransactionArrayAndReturnFalse(t *testing.T) {
	instance := &CommonStockTypeTransactionFilter{lookForString: "common stock"}
	transactions := make([]transaction_xml_parsing.NonDerivativeTransaction, 2)
	transactions[1] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature:        transaction_xml_parsing.OwnerShipNature{},
	}
	transactions[0] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Preferred Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature:        transaction_xml_parsing.OwnerShipNature{},
	}
	rawTransaction := &transaction_xml_parsing.RawOwnershipDocument{
		XMLName:        xml.Name{},
		Issuer:         transaction_xml_parsing.Issuer{},
		ReportingOwner: transaction_xml_parsing.ReportingOwner{},
		NonDerivativeTable: &transaction_xml_parsing.NonDerivativeTable{
			XMLName:      xml.Name{},
			Transactions: &transactions,
		},
	}
	expectArrayLength := 0
	got, _ := instance.ShouldKeep(rawTransaction)
	if got != false {
		t.Fatalf("Expected to get false return type but got %t", got)
	}
	if len(*rawTransaction.NonDerivativeTable.Transactions) != expectArrayLength {
		t.Fatalf("Expected transaction array length %d, received %d", expectArrayLength,
			len(*rawTransaction.NonDerivativeTable.Transactions))
	}
}

//IndirectFilter
func TestIndirectStockFilter_ShouldKeep(t *testing.T) {
	instance := &IndirectStockFilter{}
	transactions := make([]transaction_xml_parsing.NonDerivativeTransaction, 2)
	transactions[1] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature: transaction_xml_parsing.OwnerShipNature{
			DirectOrIndirectOwnership: transaction_xml_parsing.HasStringValue{
				Value: "I",
			},
		},
	}
	transactions[0] = transaction_xml_parsing.NonDerivativeTransaction{
		XMLName: xml.Name{},
		SecurityTitle: transaction_xml_parsing.HasStringValue{
			Value: "Preferred Stock",
		},
		TransactionDate:        transaction_xml_parsing.HasStringValue{},
		TransactionAmounts:     transaction_xml_parsing.TransactionAmounts{},
		PostTransactionAmounts: transaction_xml_parsing.PostTransactionAmounts{},
		OwnerShipNature: transaction_xml_parsing.OwnerShipNature{
			DirectOrIndirectOwnership: transaction_xml_parsing.HasStringValue{
				Value: "D",
			},
		},
	}
	rawTransaction := &transaction_xml_parsing.RawOwnershipDocument{
		XMLName:        xml.Name{},
		Issuer:         transaction_xml_parsing.Issuer{},
		ReportingOwner: transaction_xml_parsing.ReportingOwner{},
		NonDerivativeTable: &transaction_xml_parsing.NonDerivativeTable{
			XMLName:      xml.Name{},
			Transactions: &transactions,
		},
	}
	expectArrayLength := 1
	got, _ := instance.ShouldKeep(rawTransaction)
	if got != true {
		t.Fatalf("Expected to get true return type but got %t", got)
	}
	if len(*rawTransaction.NonDerivativeTable.Transactions) != expectArrayLength {
		t.Fatalf("Expected transaction array length %d, received %d", expectArrayLength,
			len(*rawTransaction.NonDerivativeTable.Transactions))
	}
}

type companyGetterMock struct {
	Return *requests.CompanyDetails
}

func (c companyGetterMock) GetCompanyDetails(symbol string) (*requests.CompanyDetails, error) {
	return c.Return, nil
}

func TestStockExchangeTypeFilter_ShouldKeep(t *testing.T) {
	var companyGetter *companyGetterMock = &companyGetterMock{
		Return: nil,
	}
	exchanges := []edgarwebcrawler.StockExchange{
		edgarwebcrawler.NASDAQ,
		edgarwebcrawler.NYSE,
	}
	filter := &StockExchangeTypeFilter{
		keepExchanges: exchanges,
		companyGetter: companyGetter,
	}
	rawTransaction := &transaction_xml_parsing.RawOwnershipDocument{
		XMLName: xml.Name{},
		Issuer: transaction_xml_parsing.Issuer{
			XMLName:             xml.Name{},
			IssuerCIK:           "",
			IssuerName:          "",
			IssuerTradingSymbol: "KAKI",
		},
		ReportingOwner: transaction_xml_parsing.ReportingOwner{},
		NonDerivativeTable: &transaction_xml_parsing.NonDerivativeTable{
			XMLName:      xml.Name{},
			Transactions: nil,
		},
	}
	t.Run("should return false for non NYSE or NASDAQ company", func(test *testing.T) {
		companyGetter.Return = &requests.CompanyDetails{
			Sector:    "",
			Exchange:  "TLV",
			Name:      "",
			MarketCap: "",
		}
		got, _ := filter.ShouldKeep(rawTransaction)
		if got != false {
			t.Fatalf("Received true value for exchange: TLV")
		}
	})
	t.Run("should return true for NYSE exchange", func(test *testing.T) {
		companyGetter.Return = &requests.CompanyDetails{
			Sector:    "",
			Exchange:  "NYSE",
			Name:      "",
			MarketCap: "",
		}
		got, _ := filter.ShouldKeep(rawTransaction)
		if got != true {
			t.Fatalf("Received false value for exchange: NYSE")
		}
	})
}
