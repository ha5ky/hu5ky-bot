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
	"io/fs"
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

	fileSystem := os.DirFS(*dirPath)

	if err = fs.WalkDir(fileSystem, ".", func(fileName string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ef, err := excelize.OpenFile(path.Join(*dirPath, fileName))
		if err != nil {
			panic(err)
		}

		for _, sheet := range ef.GetSheetList() {
			rows, _ := ef.GetRows(sheet)
			rowCount := len(rows)
			fmt.Println(rowCount)
			for i := *startRow; i <= rowCount; i++ {
				prompt, _ := ef.GetCellValue(sheet, getCellKey(*promptCol, i))
				completion, _ := ef.GetCellValue(sheet, getCellKey(*completionCol, i))
				models = append(models, model{
					Prompt:     prompt,
					Completion: completion,
				})
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
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

func getKey(index int) string {
	colCode := ""
	key := 'A'
	loop := index / 26
	if loop > 0 {
		colCode += getKey(loop - 1)
	}
	return colCode + string(key+int32(index)%26)
}

func getCellKey(x string, y int) string {
	return strings.ToUpper(x) + strconv.Itoa(y)
}

func DumpPretty(input interface{}) {
	bs, _ := json.Marshal(input)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("%v\n", out.String())
}
