/**
 * @Author Nil
 * @Description model/base/var.go
 * @Date 2023/3/28 17:29
 **/

package base

var (
	TableRegister []AutoRegister
)

func init() {
	TableRegister = make([]AutoRegister, 0)
}
