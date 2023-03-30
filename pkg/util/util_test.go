/**
 * @Author Nil
 * @Description pkg/util/util_test.go
 * @Date 2023/3/28 14:33
 **/

package util

import (
	"fmt"
	"testing"
)

func TestMergeAny(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {
		s := []any{"qqq", "www"}
		fmt.Println(MergeAny(s...))
	})

}
