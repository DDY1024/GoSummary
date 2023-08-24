package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("/tmp/yao1069174586/bas.xlsx")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rows, _ := f.GetRows("Sheet1")
	fmt.Println(rows)
	// fmt.Println(rows[0])
	// fmt.Println(rows[1])
	// fmt.Println(rows[2])

	// ff := excelize.NewFile()
	// defer ff.Close()
	// // idx, err := ff.NewSheet("Sheet2")
	// // fmt.Println(idx, err)
	// ff.SetCellValue("Sheet1", "A2", "Name")
	// ff.SetCellValue("Sheet1", "B2", "Age")
	// // fmt.Println(ff.SetSheetRow("s1", "A1", &[]interface{}{"Name", "Age"}))
	// // fmt.Println(ff.SetSheetRow("s1", "A2", &[]interface{}{"wxy", 18}))
	// // ff.SetActiveSheet(idx)
	// fmt.Println(ff.SaveAs("./data/xx.xlsx"))
}

// func main() {
// 	f := excelize.NewFile()
// 	// Create a new sheet.
// 	index, err := f.NewSheet("Sheet2")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Set value of a cell.
// 	//  f.SetCellValue("Sheet2", "A1", "Name")
// 	// f.SetCellValue("Sheet2", "B1", "Age")
// 	// Set active sheet of the workbook.
// 	f.SetSheetRow("Sheet2", "A1", &[]string{"Name", "Age"})
// 	f.SetSheetRow("Sheet2", "A2", &[]interface{}{"wxy", 18})
// 	f.SetSheetRow("Sheet2", "A3", &[]interface{}{"wfl", 18})
// 	f.SetActiveSheet(index)
// 	// Save spreadsheet by the given path.
// 	if err := f.SaveAs("test.xlsx"); err != nil {
// 		fmt.Println(err)
// 	}
// 	f.Close()

// 	ff, err := excelize.OpenFile("test.xlsx")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer ff.Close()

// 	fmt.Println(ff.GetRows("Sheet2"))
// }
