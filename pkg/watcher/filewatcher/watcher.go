/**
 * @Author Nil
 * @Description pkg/watcher/filewatcher/watcher.go
 * @Date 2022/10/5 17:20
 **/

package filewatcher

import (
	"fmt"
	"time"
)

func NewWatcher() FileWatcher {
	return NewFileWatcher()
}

func AddFileWatcher(fileWatcher FileWatcher, file string, callback func()) {
	err := fileWatcher.Add(file)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		var timerC <-chan time.Time
		for {
			select {
			case <-timerC:
				timerC = nil
				callback()
			case <-fileWatcher.Events(file):
				// Use a timer to debounce configuration updates
				if timerC == nil {
					timerC = time.After(100 * time.Millisecond)
				}
			}
		}
	}()
}
