/**
 * @Author Nil
 * @Description api/handler.go
 * @Date 2023/3/28 21:04
 **/

package api

import (
	"github.com/gin-gonic/gin"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"net/http"
	"reflect"
)

const RespMsg = "message"

const (
	KubernetesErrCode = -1
	BotErrCode        = -2

	Success StatusType = "success"
	Failure StatusType = "failure"
)

type Resp map[string]interface{}
type StatusType string

type Response struct {
	Code   int                `json:"code"`
	Status StatusType         `json:"status"`
	Error  boterrors.APIError `json:"error"`
	Data   interface{}        `json:"data"`
}

type Data struct {
	TotalSize int         `json:"total_size"`
	Items     interface{} `json:"items"`
}

// AssertErrorType is used to get actually type fo error, and get the status code and messages.
func AssertErrorType(err interface{}) (code uint, reason, message string) {
	if v, ok := err.(*boterrors.APIError); ok {
		return uint(v.Code), string(v.Details.Reason), v.Details.Message
	}
	return 0, "", err.(error).Error()
}

// ErrorRender when statusCode is -1, it means e for k8s error,
// ErrorRender will extract error code automatically .
func ErrorRender(ctx *gin.Context, statusCode int, e error, errMsg string) {
	if code, reason, msg := AssertErrorType(e); statusCode == BotErrCode {
		ctx.JSON(http.StatusOK, Response{
			Code:   statusCode,
			Status: Failure,
			Error: boterrors.APIError{
				Code:    statusCode,
				Message: errMsg,
				Details: boterrors.Details{
					Type:    boterrors.BotError,
					Reason:  boterrors.StatusReason(reason),
					Message: msg,
				},
			},
		})
	} else {
		ctx.JSON(http.StatusOK, Response{
			Code:   int(code),
			Status: Failure,
			Error: boterrors.APIError{
				Code:    statusCode,
				Message: errMsg,
				Details: boterrors.Details{
					Type:    boterrors.NormalError,
					Reason:  boterrors.StatusReason(reason),
					Message: msg,
				},
			},
		})
	}
}

func Error(ctx *gin.Context, e boterrors.APIError) {
	ctx.JSON(http.StatusOK, Response{
		Code:   e.Code,
		Status: Failure,
		Error:  e,
	})
}

func OK(ctx *gin.Context, data any, total int) {
	var kind reflect.Kind
	ttype := reflect.TypeOf(data)
	if ttype == nil {
		kind = reflect.Invalid
	} else {
		kind = ttype.Kind()
	}
	switch kind {
	case reflect.Struct, reflect.Invalid:
		ctx.JSON(http.StatusOK, Response{
			Code:   http.StatusOK,
			Status: Success,
			Error:  boterrors.APIError{},
			Data:   data,
		})
		return
	}
	ctx.JSON(http.StatusOK, Response{
		Code:   http.StatusOK,
		Status: Success,
		Error:  boterrors.APIError{},
		Data: Data{
			TotalSize: total,
			Items:     data,
		},
	})
}
