package DAL

import "time"

type StockExchange string

const (
	NYSE   StockExchange = "NYSE"
	NASDAQ               = "NASDAQ"
)

type Officer struct {
	ID              uint `gorm:"primaryKey"`
	OfficialTitle   string
	TitleVariations []string
}

type Insider struct {
	ID               uint   `gorm:"primaryKey"`
	CIK              string `gorm:"unique"`
	Name             string
	InsiderPositions []InsiderPosition
}

type Company struct {
	ID            uint   `gorm:"primaryKey"`
	Symbol        string `gorm:"not null;unique"`
	Name          string
	Sector        string
	StockExchange StockExchange `gorm:"not null"`
}

type InsiderPosition struct {
	ID uint `gorm:"primaryKey"`

	//relationships
	InsiderID uint
	Insider   Insider
	CompanyID uint
	Company   Company
	OfficerID uint
	Officer   Officer

	OfficerText       string
	OtherText         string
	IsDirector        bool
	IsTenPercentOwner bool
}

type Transaction struct {
	ID uint `gorm:"primaryKey"`

	Date                 time.Time
	IsAcquired           bool
	NumOfShares          float64
	PricePerShare        float64
	SharesOwnedFollowing float64
	IsDirectOwnership    bool

	InsiderPositionId uint
	InsiderPosition   InsiderPosition
}
