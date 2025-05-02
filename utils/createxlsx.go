package utils

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter interface {
	SetHeaders(headers []string)
	SetRows(rows [][]interface{})
	Generate() (*bytes.Buffer, error)
}

type TaomExcelExporter struct {
	file      *excelize.File
	sheetName string
}

func NewTaomExcelExporter() *TaomExcelExporter {
	f := excelize.NewFile()
	sheet := "Sheet1"
	index, err := f.NewSheet(sheet)
	if err != nil {
		fmt.Println(err.Error())
	}
	f.SetActiveSheet(index)
	return &TaomExcelExporter{
		file:      f,
		sheetName: sheet,
	}
}
func (t *TaomExcelExporter) SetHeaders(headers []string) {
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		t.file.SetCellValue(t.sheetName, cell, header)
	}
}

func (t *TaomExcelExporter) SetRows(rows [][]interface{}) {
	for rowIdx, row := range rows {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			t.file.SetCellValue(t.sheetName, cell, value)
		}
	}
}

func (t *TaomExcelExporter) Generate() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := t.file.Write(buf); err != nil {
		return nil, err
	}
	return buf, nil
}
