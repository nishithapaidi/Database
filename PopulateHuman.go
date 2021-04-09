package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tealeg/xlsx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func populateHuman() error {
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

	excelFileName := HUMAN_XLSX
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(xlFile.Sheets); i++ {
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
		db.Table("human_sheet").Create(&dbSheet)

		for i := 1; i < sheet.MaxRow; i++ {
			row, err := sheet.Row(i)
			if err != nil {
				return err
			}

			value := row.GetCell(2).Value
			value_log2foldchange := row.GetCell(3).Value
			value_padj := row.GetCell(4).Value
			if value_padj == "NA" {
				value_padj = ""
			}

			if (value != "") && (value_log2foldchange != "") {
				expressionUuid := generateUUID()
				expression := map[string]interface{}{}
				expression["id"] = expressionUuid
				expression["sheet_id"] = dbSheet["id"]
				expression["gene"] = row.GetCell(0).Value
				expression["name"] = row.GetCell(1).Value
				expression["value"] = value
				expression["value_log2foldchange"] = value_log2foldchange
				if value_padj != "NA" {
					expression["value_padj"] = value_padj
				}
				db.Table("human_gene_data").Create(&expression)
			}
		}

		if err != nil {
			fmt.Println(err)
		}

	}

	db.Table("")

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}

	sqlDB.Close()

	return nil
}
