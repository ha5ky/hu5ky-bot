/**
 * @Author Nil
 * @Description model/init.go
 * @Date 2023/3/28 20:21
 **/

package model

import (
	"fmt"
	"github.com/ha5ky/hu5ky-bot/model/base"
)

func Registry() {
	controller := NewController()
	fmt.Println(base.TableRegister)
	if err := controller.CreateTables(base.TableRegister); err != nil {
		panic(err)
	}
}
