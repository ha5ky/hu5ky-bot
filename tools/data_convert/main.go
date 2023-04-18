/**
 * @Author Nil
 * @Description tools/data_convert/main.go
 * @Date 2023/3/27 19:38
 **/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type model struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

var (
	dirPath       = flag.String("dirPath", "./data", "")
	promptCol     = flag.String("promptCol", "D", "")
	completionCol = flag.String("completionCol", "E", "")
	startRow      = flag.Int("startRow", 1, "")
)

func main() {
	flag.Parse()

	var (
		err error
	)
	var (
		models = make([]model, 0)
	)

	ef, err := excelize.OpenFile(path.Join(*dirPath, "辽宁22年.xlsx"))
	if err != nil {
		panic(err)
	}

	for _, sheet := range ef.GetSheetList() {
		rows, _ := ef.GetRows(sheet)
		rowCount := len(rows)
		rows = nil
		fmt.Println(rowCount)
		for i := *startRow; i <= rowCount; i++ {
			prompt := ""
			completion := sprintf(
				getCellValue(ef, sheet, "B", i),
				getCellValue(ef, sheet, "A", i),
				"2022",
				getCellValue(ef, sheet, "F", i),
				getCellValue(ef, sheet, "D", i),
				getCellValue(ef, sheet, "C", i),
				getCellValue(ef, sheet, "E", i),
			)
			m := model{
				Prompt:     prompt,
				Completion: completion,
			}

			models = append(models, m)
		}
	}
	file, _ := os.OpenFile("./json.txt", os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, model := range models {
		jsonByte, _ := json.Marshal(model)
		write.WriteString(string(jsonByte) + "\n")
	}
	write.Flush()
}

func getCellValue(ef *excelize.File, sheet, key string, row int) string {
	ret, _ := ef.GetCellValue(sheet, getCellKey(key, row))
	return ret
}

func sprintf(collage, no, year, kind, major, majorNum, score string) string {
	return fmt.Sprintf(
		"招生院校%s的院校编号是%s，在%s年，%s科类，招生专业是%s（专业编号%s）的最低投档分是%s",
		collage, no, year, kind, major, majorNum, score,
	)
}

func getKey(index int) string {
	colCode := ""
	key := 'A'
	loop := index / 26
	if loop > 0 {
		colCode += getKey(loop - 1)
	}
	return colCode + string(key+int32(index)%26)
}

func getCellKey(x any, y int) (ret string) {

	switch x.(type) {
	case string:
		ret = strings.ToUpper(x.(string)) + strconv.Itoa(y)
	case int:
		ret = getKey(x.(int))
	}
	return
}

func DumpPretty(input interface{}) {
	bs, _ := json.Marshal(input)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("%v\n", out.String())
}
