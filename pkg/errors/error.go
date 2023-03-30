/**
 * @Author Nil
 * @Description pkg/errors/error.go
 * @Date 2023/3/28 21:06
 **/

package errors

type APIError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Details Details `json:"details"`
}

func (a *APIError) Error() string {
	return a.Details.Message
}

type Details struct {
	Type    ErrorType    `json:"type"`
	Reason  StatusReason `json:"reason"`
	Message string       `json:"message"`
}

type StatusReason string
type ErrorType string

const (
	KubernetesError ErrorType = "KubernetesError"
	BotError        ErrorType = "BotError"
	NormalError     ErrorType = "NormalError"
)

var (
	InvalidParams = "invalid params"

	DatabaseQueryError   = "database query error"
	NotFoundError        = "not found error"
	InternalError        = "internal error"
	ResourceOperateError = "resource operate failed"

	SensitiveDataError = "can not modified sensitive data"
)
