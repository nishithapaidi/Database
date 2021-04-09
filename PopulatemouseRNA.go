package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tealeg/xlsx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func populateRnaDb() error {
	// CONNECT TO DATABASE WITH SPECIAL VALUES EXAMPLE
	content, err := ioutil.ReadFile(DATABASE_CONFIG)
	if err != nil {
		fmt.Println(err)
	}
	dsn := string(content)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	excelFileName := MOUSE_XLSX
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(xlFile.Sheets)-1; i++ {
		sheet := xlFile.Sheets[i]

		result := []map[string]interface{}{}

		db.Raw("SELECT sheet.name FROM gene.sheet sheet WHERE sheet.name = '" + sheet.Name + "';").Scan(&result)

		if len(result) >= 1 {
			continue
		}

		fmt.Println("Populating: ", sheet.Name)

		sheetUuid := generateUUID()
		dbSheet := map[string]interface{}{}
		dbSheet["id"] = sheetUuid
		dbSheet["name"] = sheet.Name
		db.Table("sheet").Create(&dbSheet)

		// name row
		classNameRow, err := sheet.Row(1)
		if err != nil {
			fmt.Println(err)
		}

		var classMap map[string]string
		classMap = make(map[string]string)
		for n := 1; n < sheet.MaxCol; n++ {
			className := classNameRow.GetCell(n)
			name := className.Value
			if classMap[name] == "" {
				classUuid := generateUUID()
				classMap[name] = classUuid
				classType := map[string]interface{}{}
				classType["id"] = classUuid
				classType["name"] = name
				db.Table("class_type").Create(&classType)
			}
		}

		rowName, err := sheet.Row(0)
		if err != nil {
			fmt.Println(err)
		}
		rowNameMap := make(map[int]string)
		for n := 1; n < sheet.MaxCol; n++ {
			trialName := rowName.GetCell(n).Value

			sampleUuid := generateUUID()
			sample := map[string]interface{}{}
			sample["id"] = sampleUuid
			sample["sheet_id"] = sheetUuid
			sample["class_id"] = classMap[classNameRow.GetCell(n).Value]
			sample["name"] = trialName
			sample["index"] = n
			db.Table("sample").Create(&sample)

			rowNameMap[n] = sampleUuid

		}

		rowExpressionMap := make(map[string]string)
		for j := 2; j < sheet.MaxRow; j++ {
			rowIndex := j
			row, err := sheet.Row(j)
			if err != nil {
				fmt.Println(err)
			}

			expressionName := row.GetCell(0).String()
			if rowExpressionMap[expressionName] == "" {
				// ADD TO DATABASE EXAMPLE
				expressionUuid := generateUUID()
				expression := map[string]interface{}{}
				expression["id"] = expressionUuid
				expression["sheet_id"] = dbSheet["id"]
				expression["name"] = expressionName
				db.Table("expression").Create(&expression)

				rowExpressionMap[expressionName] = expressionUuid
			}

			rowVisitor(db, row, dbSheet, rowNameMap, rowExpressionMap[expressionName], rowIndex)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}

	sqlDB.Close()

	return nil
}

func rowVisitor(db *gorm.DB, r *xlsx.Row, dbSheet map[string]interface{}, rowNameMap map[int]string, expressionUuid string, rowIndex int) error {
	// Operate initial value
	trials := []map[string]interface{}{}
	for i := 1; i < r.Sheet.MaxCol; i++ {
		trialValue := r.GetCell(i).Value

		trialUuid := generateUUID()
		trial := map[string]interface{}{}
		trial["id"] = trialUuid
		trial["expression_id"] = expressionUuid
		trial["sample_id"] = rowNameMap[i]
		trial["sheet_id"] = dbSheet["id"]
		trial["row_index"] = rowIndex
		if trialValue == "" {
			trial["value"] = nil
		} else {
			trial["value"] = trialValue
		}
		trials = append(trials, trial)

	}
	db.Table("trial").Create(&trials)

	return nil
}
