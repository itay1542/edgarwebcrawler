package DAL

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/itay1542/edgarwebcrawler/transaction_xml_parsing"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strings"
)

type DBConfiguration struct {
	Host, User, Password, DBName string
	Port                         uint
}

type InsideOutDB interface {
	Init() error
	GetOfficerTitleId(officerType transaction_xml_parsing.OfficerType) (uint, error)
	DoesCompanyExist(symbol string) (uint, error)
	DoesInsiderPositionExist(companyId, insiderId uint) (uint, error)
	DoesInsiderExist(cik string) (uint, error)
	AddCompany(company *Company) (*Company, error)
	AddTransactions([]Transaction) ([]Transaction, error)
	AddInsider(insider *Insider) (*Insider, error)
	AddInsiderPosition(position *InsiderPosition, officerPositions []*Officer) (*InsiderPosition, error)
}

type PostgresInsideOut struct {
	Config DBConfiguration
	db     *gorm.DB
}

func (p *PostgresInsideOut) Init() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jerusalem",
		p.Config.Host, p.Config.User, p.Config.Password, p.Config.DBName, p.Config.Port)
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}))
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Officer{}, &Company{}, &Insider{}, &InsiderPosition{}, &Transaction{})
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *PostgresInsideOut) GetOfficerTitleId(officerType transaction_xml_parsing.OfficerType) (uint, error) {
	var officer Officer
	if err := p.db.Where("official_title = ?", strings.ToLower(string(officerType))).
		First(&officer).Error; err != nil {
		return 0, err
	}
	return officer.ID, nil
}

func (p *PostgresInsideOut) DoesCompanyExist(symbol string) (uint, error) {
	var company Company
	err := p.db.Where("symbol = ?", symbol).First(&company).Error
	if err == nil {
		return company.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else {
		return 0, err
	}
}

func (p *PostgresInsideOut) AddCompany(company *Company) (*Company, error) {
	var tempCompany Company = *company
	err := p.db.Create(&tempCompany).Error
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Printf("Added company with name %s, id %d", tempCompany.Name, tempCompany.ID)
	return &tempCompany, nil

}

func (p *PostgresInsideOut) AddTransactions(transactions []Transaction) ([]Transaction, error) {
	var tempTransactions []Transaction = make([]Transaction, len(transactions))
	copy(tempTransactions, transactions)
	err := p.db.Create(tempTransactions).Error
	if err != nil {
		return nil, err
	}
	log.Printf("Added %d transactions, first transaction with id %d", len(tempTransactions), tempTransactions[0].ID)
	return tempTransactions, nil

}

func (p *PostgresInsideOut) AddInsider(insider *Insider) (*Insider, error) {
	var temp Insider = *insider
	err := p.db.Create(&temp).Error
	if err != nil {
		return nil, err
	}
	log.Printf("Added insider %s, id: %d", temp.Name, temp.ID)
	return &temp, nil
}

func (p *PostgresInsideOut) DoesInsiderPositionExist(companyId, insiderId uint) (uint, error) {
	var insiderPosition InsiderPosition
	err := p.db.Where("company_id = ? AND insider_id = ?", companyId, insiderId).First(&insiderPosition).Error
	if err == nil {
		return insiderPosition.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else {
		return 0, err
	}
}

func (p *PostgresInsideOut) AddInsiderPosition(position *InsiderPosition, officerPositions []*Officer) (*InsiderPosition, error) {
	var temp InsiderPosition = *position
	err := p.db.Create(&temp).Error
	if err != nil {
		return nil, err
	}
	if len(officerPositions) > 0 {
		err = p.db.Model(&temp).Association("Officers").Append(officerPositions)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Added Insider Position %+v", temp)
	return &temp, nil
}

func (p *PostgresInsideOut) DoesInsiderExist(cik string) (uint, error) {
	var insider Insider
	err := p.db.Where("cik = ?", cik).First(&insider).Error
	if err == nil {
		return insider.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	} else {
		return 0, err
	}
}
