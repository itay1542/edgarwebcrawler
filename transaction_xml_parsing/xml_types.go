package transaction_xml_parsing

import "encoding/xml"

type RawOwnershipDocument struct {
	XMLName            xml.Name            `xml:"ownershipDocument"`
	Issuer             Issuer              `xml:"issuer"`
	ReportingOwner     ReportingOwner      `xml:"reportingOwner"`
	NonDerivativeTable *NonDerivativeTable `xml:"nonDerivativeTable,omitempty"`
}

type Issuer struct {
	XMLName             xml.Name `xml:"issuer"`
	IssuerCIK           string   `xml:"issuerCik"`
	IssuerName          string   `xml:"issuerName"`
	IssuerTradingSymbol string   `xml:"issuerTradingSymbol"`
}

type ReportingOwner struct {
	XMLName                    xml.Name                   `xml:"reportingOwner"`
	ReportingOwnerId           ReportingOwnerId           `xml:"reportingOwnerId"`
	ReportingOwnerRelationship ReportingOwnerRelationship `xml:"reportingOwnerRelationship"`
}

type ReportingOwnerId struct {
	XMLName            xml.Name `xml:"reportingOwnerId"`
	ReportingOwnerCIK  string   `xml:"rptOwnerCik"`
	ReportingOwnerName string   `xml:"rptOwnerName"`
}

type ReportingOwnerRelationship struct {
	XMLName           xml.Name `xml:"reportingOwnerRelationship"`
	IsDirector        byte     `xml:"isDirector"`
	IsOfficer         byte     `xml:"isOfficer"`
	IsTenPercentOwner byte     `xml:"isTenPercentOwner"`
	IsOther           byte     `xml:"isOther"`
	OfficerTitle      *string  `xml:"officerTitle,omitempty"`
	OtherText         *string  `xml:"otherText,omitempty"`
}

type NonDerivativeTable struct {
	XMLName      xml.Name                    `xml:"nonDerivativeTable"`
	Transactions *[]NonDerivativeTransaction `xml:"nonDerivativeTransaction,omitempty"`
}

type NonDerivativeTransaction struct {
	XMLName                xml.Name               `xml:"nonDerivativeTransaction"`
	SecurityTitle          HasStringValue         `xml:"securityTitle"`
	TransactionDate        HasStringValue         `xml:"transactionDate"`
	TransactionAmounts     TransactionAmounts     `xml:"transactionAmounts"`
	PostTransactionAmounts PostTransactionAmounts `xml:"postTransactionAmounts"`
	OwnerShipNature        OwnerShipNature        `xml:"ownerShipNature"`
}

type HasStringValue struct {
	Value string `xml:"value"`
}

type HasFloatValue struct {
	Value float64 `xml:"value"`
}

type TransactionAmounts struct {
	XMLName                         xml.Name       `xml:"transactionAmounts"`
	TransactionShares               HasFloatValue  `xml:"transactionShares"`
	TransactionPricePerShare        HasFloatValue  `xml:"transactionPricePerShare"`
	TransactionAcquiredDisposedCode HasStringValue `xml:"transactionAcquiredDisposedCode"`
}

type PostTransactionAmounts struct {
	XMLName                         xml.Name      `xml:"postTransactionAmounts"`
	SharesOwnedFollowingTransaction HasFloatValue `xml:"sharesOwnedFollowingTransaction"`
}

type OwnerShipNature struct {
	DirectOrIndirectOwnership HasStringValue `xml:"directOrIndirectOwnership"`
}
