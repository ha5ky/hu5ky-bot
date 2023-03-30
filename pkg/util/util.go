/**
 * @Author Nil
 * @Description pkg/util/util.go
 * @Date 2023/3/28 14:29
 **/

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func MergeAny(args ...any) (ret string) {
	for _, arg := range args {
		ret += fmt.Sprintf(" %+v ", arg)
	}
	return
}

func DumpPretty(input interface{}) {
	bs, _ := json.Marshal(input)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("%v\n", out.String())
}
