package transaction_xml_parsing

import "encoding/xml"

type RawOwnershipDocument struct {
	XMLName xml.Name `xml:"ownershipDocument"`
	Issuer issuer `xml:"issuer"`
	ReportingOwner reportingOwner `xml:"reportingOwner"`
	NonDerivativeTable nonDerivativeTable `xml:"nonDerivativeTable"`
}

type issuer struct {
	XMLName xml.Name `xml:"issuer"`
	IssuerCIK string `xml:"issuerCik"`
	IssuerName string `xml:"issuerName"`
	IssuerTradingSymbol string `xml:"issuerTradingSymbol"`
}

type reportingOwner struct {
	XMLName xml.Name `xml:"reportingOwner"`
	ReportingOwnerId reportingOwnerId `xml:"reportingOwnerId"`
	ReportingOwnerRelationship reportingOwnerRelationship `xml:"reportingOwnerRelationship"`
}

type reportingOwnerId struct {
	XMLName xml.Name `xml:"reportingOwnerId"`
	ReportingOwnerCIK string `xml:"rptOwnerCik"`
	ReportingOwnerName string `xml:"rptOwnerName"`
}

type reportingOwnerRelationship struct {
	XMLName xml.Name `xml:"reportingOwnerRelationship"`
	IsDirector byte `xml:"isDirector"`
	IsOfficer byte `xml:"isOfficer"`
	IsTenPercentOwner byte `xml:"isTenPercentOwner"`
	IsOther byte `xml:"isOther"`
	OfficerTitle *string `xml:"officerTitle,omitempty"`
	OtherText *string `xml:"otherText,omitempty"`
}

type nonDerivativeTable struct {
	XMLName xml.Name `xml:"nonDerivativeTable"`
	Transactions *[]nonDerivativeTransaction `xml:"nonDerivativeTransaction,omitempty"`
}

type nonDerivativeTransaction struct {
	XMLName xml.Name `xml:"nonDerivativeTransaction"`
	SecurityTitle hasStringValue `xml:"securityTitle"`
	TransactionDate hasStringValue `xml:"transactionDate"`
	TransactionAmounts transactionAmounts `xml:"transactionAmounts"`
	PostTransactionAmounts postTransactionAmounts `xml:"postTransactionAmounts"`
	OwnerShipNature ownerShipNature `xml:"ownerShipNature"`
}

type hasStringValue struct {
	Value string `xml:"value"`
}

type hasFloatValue struct {
	Value float64 `xml:"value"`
}

type hasByteValue struct {
	Value byte `xml:"value"`
}

type transactionAmounts struct {
	XMLName xml.Name `xml:"transactionAmounts"`
	TransactionShares hasFloatValue `xml:"transactionShares"`
	TransactionPricePerShare hasFloatValue `xml:"transactionPricePerShare"`
	TransactionAcquiredDisposedCode hasByteValue `xml:"transactionAcquiredDisposedCode"`
}

type postTransactionAmounts struct {
	XMLName xml.Name `xml:"postTransactionAmounts"`
	SharesOwnedFollowingTransaction hasFloatValue `xml:"sharesOwnedFollowingTransaction"`
}

type ownerShipNature struct {
	XMLName xml.Name `xml:"ownershipNature"`
	DirectOrIndirectOwnership hasByteValue `xml:"directOrIndirectOwnership"`
}
