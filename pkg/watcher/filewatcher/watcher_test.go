/**
 * @Author Nil
 * @Description pkg/watcher/filewatcher/watcher_test.go
 * @Date 2022/10/5 17:20
 **/

package filewatcher

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestAddFileWatcher(t *testing.T) {

	t.Run("test fileWatcher", func(t *testing.T) {
		watcher := NewWatcher()
		fmt.Println(os.Getwd())
		AddFileWatcher(watcher, "../../../test/fileWatcher", func() {
			contentBytes, err := os.ReadFile("../../../test/fileWatcher")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(contentBytes))
		})
		time.Sleep(time.Hour)
	})

}
