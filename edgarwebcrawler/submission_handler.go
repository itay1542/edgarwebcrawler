package edgarwebcrawler

import (
	"fmt"
	"github.com/itay1542/edgarwebcrawler/DAL"
	"github.com/itay1542/edgarwebcrawler/requests"
	"github.com/itay1542/edgarwebcrawler/transaction_xml_parsing"
	"github.com/itay1542/edgarwebcrawler/transactions"
	"log"
	"sort"
	"time"
)

type SubmissionHandler interface {
	HandleSubmission(submission *transaction_xml_parsing.RawOwnershipDocument) error
}

//SecSubmissionHandler implements SubmissionHandler
type SecSubmissionHandler struct {
	companyGetter     requests.CompanyGetter
	officerClassifier transaction_xml_parsing.OfficerClassifier
	dal               DAL.InsideOutDB
	filters           []transactions.TransactionFilterer
}

func NewSecSubmissionHandler(dal DAL.InsideOutDB, filters []transactions.TransactionFilterer,
	officerClassifier transaction_xml_parsing.OfficerClassifier,
	companyGetter requests.CompanyGetter) *SecSubmissionHandler {
	return &SecSubmissionHandler{
		companyGetter:     companyGetter,
		officerClassifier: officerClassifier,
		dal:               dal,
		filters:           filters,
	}
}

func (s *SecSubmissionHandler) HandleSubmission(submission *transaction_xml_parsing.RawOwnershipDocument) error {
	log.Printf("Received submission from %+v", submission.Issuer)
	if !s.shouldProcess(submission) {
		log.Printf("Dropping submission due to filters")
		return nil
	}

	companyId, err := s.saveCompany(submission.Issuer.IssuerTradingSymbol)
	if err != nil {
		return err
	}
	insiderId, err := s.saveInsider(&submission.ReportingOwner.ReportingOwnerId)
	if err != nil {
		return err
	}
	insiderPositionId, err := s.saveInsiderPositions(insiderId, companyId, &submission.ReportingOwner.ReportingOwnerRelationship)
	if err != nil {
		return err
	}
	transactionModels := s.extractTransactions(*submission.NonDerivativeTable.Transactions, insiderPositionId)
	if err != nil {
		return err
	}
	transactionModels, err = s.dal.AddTransactions(transactionModels)
	if err != nil {
		return err
	}
	return nil
}

func (s SecSubmissionHandler) extractTransactions(transactions []transaction_xml_parsing.NonDerivativeTransaction,
	insiderPositionId uint) []DAL.Transaction {
	transactionModels := make([]DAL.Transaction, 0)
	for _, val := range transactions {
		parsedTransactionDate, err := time.Parse(TRANSACTION_DATE_LAYOUT, val.TransactionDate.Value)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		newModel := DAL.Transaction{
			Date:                 parsedTransactionDate,
			IsAcquired:           val.TransactionAmounts.TransactionAcquiredDisposedCode.Value == "A",
			NumOfShares:          val.TransactionAmounts.TransactionShares.Value,
			PricePerShare:        val.TransactionAmounts.TransactionPricePerShare.Value,
			SharesOwnedFollowing: val.PostTransactionAmounts.SharesOwnedFollowingTransaction.Value,
			IsDirectOwnership:    val.OwnerShipNature.DirectOrIndirectOwnership.Value == "D",
			InsiderPositionId:    insiderPositionId,
		}
		transactionModels = append(transactionModels, newModel)
	}
	return transactionModels
}

func (s *SecSubmissionHandler) saveCompany(symbol string) (uint, error) {
	companyId, err := s.dal.DoesCompanyExist(symbol)
	if err != nil {
		return 0, err
	}
	if companyId != 0 {
		log.Printf("Found company with symbol %s in the db, skipping save stage", symbol)
		return companyId, nil
	}

	companyDetails, err := s.companyGetter.GetCompanyDetails(symbol)
	if err != nil {
		return 0, err
	}
	log.Printf("Received company details for %s", companyDetails.Name)
	company := &DAL.Company{
		Symbol:        symbol,
		Name:          companyDetails.Name,
		Sector:        companyDetails.Sector,
		StockExchange: DAL.StockExchange(companyDetails.Exchange),
	}
	company, err = s.dal.AddCompany(company)
	if err != nil {
		return 0, err
	}
	log.Printf("Successfuly inserted company into the database, id: %d", company.ID)
	return company.ID, nil
}

func (s *SecSubmissionHandler) shouldProcess(submission *transaction_xml_parsing.RawOwnershipDocument) bool {
	sort.Slice(s.filters, func(i, j int) bool {
		return s.filters[i].Priority() > s.filters[j].Priority()
	})
	/*filterResults := make(chan bool)
	wg := sync.WaitGroup{}*/
	for _, filter := range s.filters {
		//wg.Add(1)
		filter := filter
		//func() {
		shouldKeep, _ := filter.ShouldKeep(submission)
		if !shouldKeep {
			return false
		}
		//filterResults <- shouldKeep
		//wg.Done()
		//}()
	}
	/*wg.Wait()
	close(filterResults)
	for shouldKeep := range filterResults {
		if !shouldKeep {
			return false
		}
	}*/
	return true
}

func (s *SecSubmissionHandler) saveInsider(insider *transaction_xml_parsing.ReportingOwnerId) (uint, error) {
	insiderId, err := s.dal.DoesInsiderExist(insider.ReportingOwnerCIK)
	if err != nil {
		return 0, err
	}
	if insiderId == 0 {
		var insiderModel *DAL.Insider = &DAL.Insider{
			CIK:  insider.ReportingOwnerCIK,
			Name: insider.ReportingOwnerName,
		}
		insiderModel, err = s.dal.AddInsider(insiderModel)
		if err != nil {
			return 0, nil
		}
		insiderId = insiderModel.ID
	}
	return insiderId, nil
}

func (s *SecSubmissionHandler) saveInsiderPositions(
	insiderId, companyId uint, position *transaction_xml_parsing.ReportingOwnerRelationship) (uint, error) {
	insiderPositionId, err := s.dal.DoesInsiderPositionExist(companyId, insiderId)
	if err != nil {
		return 0, err
	}
	if insiderPositionId == 0 {
		positionModel := DAL.InsiderPosition{
			InsiderID:         insiderId,
			CompanyID:         companyId,
			IsDirector:        position.IsDirector > 0,
			IsTenPercentOwner: position.IsTenPercentOwner > 0,
		}
		officers := make([]*DAL.Officer, 0)
		if position.IsOfficer > 0 {
			positionModel.OfficerText = *position.OfficerTitle
			officerTitleIds := s.getOfficerTitleIds(position)
			if len(officerTitleIds) > 0 {
				for _, officerTitleId := range officerTitleIds {
					officers = append(officers, &DAL.Officer{ID: officerTitleId})
				}
			}
		}
		if position.IsOther > 0 {
			positionModel.OtherText = *position.OtherText
		}
		insiderPosition, err := s.dal.AddInsiderPosition(&positionModel, officers)
		if err != nil {
			return 0, err
		}
		insiderPositionId = insiderPosition.ID
	}

	return insiderPositionId, nil
}

func (s SecSubmissionHandler) getOfficerTitleIds(position *transaction_xml_parsing.ReportingOwnerRelationship) []uint {
	var officerTitleIds []uint = make([]uint, 0)
	var officerTypes []transaction_xml_parsing.OfficerType
	officerTypes, err := s.officerClassifier.GetOfficerType(*position.OfficerTitle)
	if err != nil {
		log.Print(err)
	}
	for _, officerType := range officerTypes {
		officerTypeId, err := s.dal.GetOfficerTitleId(officerType)
		if err != nil {
			log.Print(err)
			continue
		}
		officerTitleIds = append(officerTitleIds, officerTypeId)
	}
	return officerTitleIds
}
