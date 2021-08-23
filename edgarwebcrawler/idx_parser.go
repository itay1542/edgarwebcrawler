package edgarwebcrawler

import (
	"bufio"
	"errors"
	"github.com/itay1542/edgarwebcrawler/utils"
	"io"
	"strings"
)

type IdxRow struct {
	FileName, FormType, CompanyName, CIK, DateFiled string
}

type IdxReader interface {
	ReadRow() (*IdxRow, error)
}

//SimpleIdxReader implements IdxReader
type SimpleIdxReader struct {
	formTypeIndex, companyNameIndex, cikIndex, dateIndex, fileNameIndex int
	reader                                                              bufio.Reader
}

func NewIdxReader(reader io.Reader) (*SimpleIdxReader, error) {
	bufReader := bufio.NewReader(reader)
	for {
		line, err := utils.ReadLine(bufReader)
		if err != nil {
			break
		}
		formTypeIndex := strings.Index(line, "Form Type")
		companyNameIndex := strings.Index(line, "Company Name")
		cikIndex := strings.Index(line, "CIK")
		dateIndex := strings.Index(line, "Date Filed")
		fileNameIndex := strings.Index(line, "File Name")
		if formTypeIndex == -1 ||
			companyNameIndex == -1 ||
			cikIndex == -1 ||
			dateIndex == -1 ||
			fileNameIndex == -1 {
			continue
		}
		return &SimpleIdxReader{
			formTypeIndex:    formTypeIndex,
			companyNameIndex: companyNameIndex,
			cikIndex:         cikIndex,
			dateIndex:        dateIndex,
			fileNameIndex:    fileNameIndex,
			reader:           *bufReader,
		}, nil
	}

	return nil, errors.New("could not find IDX Header")
}

func (p *SimpleIdxReader) ReadRow() (*IdxRow, error) {
	idxLine, err := utils.ReadLine(&p.reader)
	if err != nil {
		return nil, err
	}
	formType := strings.Trim(idxLine[p.formTypeIndex:p.companyNameIndex], " ")
	companyName := strings.Trim(idxLine[p.companyNameIndex:p.cikIndex], " ")
	cik := strings.Trim(idxLine[p.cikIndex:p.dateIndex], " ")
	date := strings.Trim(idxLine[p.dateIndex:p.fileNameIndex], " ")
	fileName := strings.Trim(idxLine[p.fileNameIndex:], " ")

	return &IdxRow{
		FileName:    fileName,
		CompanyName: companyName,
		CIK:         cik,
		DateFiled:   date,
		FormType:    formType,
	}, nil
}
