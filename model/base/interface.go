/**
 * @Author Nil
 * @Description model/interface.go
 * @Date 2023/3/28 17:27
 **/

package base

type AutoRegister interface {
	Registry()
	TableName() string
}
