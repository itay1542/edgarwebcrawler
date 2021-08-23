package transactions

import (
	edgarwebcrawler "github.com/itay1542/edgarwebcrawler/DAL"
	"github.com/itay1542/edgarwebcrawler/requests"
	"github.com/itay1542/edgarwebcrawler/transaction_xml_parsing"
	"strings"
)

type Priority uint8

const (
	FIRST  Priority = 1
	SECOND Priority = 2
)

type TransactionFilterer interface {
	ShouldKeep(transaction *transaction_xml_parsing.RawOwnershipDocument) (bool, int)
	Priority() Priority
}

//CommonStockTypeTransactionFilter implements TransactionFilterer
type CommonStockTypeTransactionFilter struct {
	lookForString string
}

func (s *CommonStockTypeTransactionFilter) Priority() Priority {
	return SECOND
}

func (s *CommonStockTypeTransactionFilter) ShouldKeep(transaction *transaction_xml_parsing.RawOwnershipDocument) (bool, int) {
	var newTransactionsArray []transaction_xml_parsing.NonDerivativeTransaction
	newTransactionsArray = *transaction.NonDerivativeTable.Transactions
	shouldKeep := false
	removedCount := 0
	for i, tran := range newTransactionsArray {
		if strings.Contains(strings.ToLower(tran.SecurityTitle.Value), strings.ToLower(s.lookForString)) {
			shouldKeep = true
			continue
		}
		removeElement(&newTransactionsArray, i-removedCount)
		removedCount++
		transaction.NonDerivativeTable.Transactions = &newTransactionsArray
	}
	return shouldKeep, removedCount
}

//IndirectStockFilter implements TransactionFilterer
type IndirectStockFilter struct {
}

func (s *IndirectStockFilter) Priority() Priority {
	return SECOND
}

func (i *IndirectStockFilter) ShouldKeep(transaction *transaction_xml_parsing.RawOwnershipDocument) (bool, int) {
	var newTransactionsArray []transaction_xml_parsing.NonDerivativeTransaction
	newTransactionsArray = *transaction.NonDerivativeTable.Transactions
	shouldKeep := false
	removedCount := 0
	for i, tran := range newTransactionsArray {
		if tran.OwnerShipNature.DirectOrIndirectOwnership.Value == "D" {
			shouldKeep = true
			continue
		}
		removeElement(&newTransactionsArray, i-removedCount)
		removedCount++
		transaction.NonDerivativeTable.Transactions = &newTransactionsArray
	}
	return shouldKeep, removedCount
}

func removeElement(slice *[]transaction_xml_parsing.NonDerivativeTransaction, index int) {
	(*slice)[index] = (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
}

type StockExchangeTypeFilter struct {
	keepExchanges []edgarwebcrawler.StockExchange
	companyGetter requests.CompanyGetter
}

func (s StockExchangeTypeFilter) ShouldKeep(transaction *transaction_xml_parsing.RawOwnershipDocument) (bool, int) {
	companyDetails, err := s.companyGetter.GetCompanyDetails(transaction.Issuer.IssuerTradingSymbol)
	if err != nil {
		return false, 0
	}
	for _, val := range s.keepExchanges {
		if string(val) == strings.ToUpper(string(companyDetails.Exchange)) {
			return true, 0
		}
	}
	return false, 0
}

func (s StockExchangeTypeFilter) Priority() Priority {
	return FIRST
}
