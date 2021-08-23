package DAL

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfiguration struct {
	Host, User, Password, DBName string
	Port                         uint
}

type InsideOutDB interface {
	Init() error
	GetOfficerTitle(officerText string) (*Officer, error)
	DoesCompanyExist(symbol string) (bool, error)
	DoesInsiderPositionExist(companyId, insiderId uint) (bool, error)
	AddCompany(company *Company) error
	AddTransactions([]Transaction) error
	AddInsider(insider *Insider) error
	AddInsiderPositions(positions []InsiderPosition) error
}

type PostgresInsideOut struct {
	config DBConfiguration
	db     *gorm.DB
}

func (p *PostgresInsideOut) Init() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jerusalem",
		p.config.Host, p.config.User, p.config.Password, p.config.DBName, p.config.Port)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *PostgresInsideOut) GetOfficerTitle(officerText string) (*Officer, error) {

	p.db.Get
}

func (p *PostgresInsideOut) DoesCompanyExist(symbol string) (bool, error) {
	panic("implement me")
}

func (p *PostgresInsideOut) AddCompany(company *Company) error {
	panic("implement me")
}

func (p *PostgresInsideOut) AddTransactions(transactions []Transaction) error {
	panic("implement me")
}

func (p *PostgresInsideOut) AddInsider(insider *Insider) error {
	panic("implement me")
}

func (p *PostgresInsideOut) DoesInsiderPositionExist(companyId, insiderId uint) (bool, error) {
	panic("implement me")
}

func (p *PostgresInsideOut) AddInsiderPositions(positions []InsiderPosition) error {
	panic("implement me")
}
