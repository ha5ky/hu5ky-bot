/**
 * @Author Nil
 * @Description api/collection/create.go
 * @Date 2023/4/2 18:05
 **/

package collection

import (
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/model"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"net/http"
)

type CreateRequest struct {
	Name        string `json:"name" form:"name"`
	FileId      uint   `json:"file_id" form:"file_id"`
	Description string `json:"description" form:"description"`
}

func Create(ctx *gin.Context) {
	var (
		req CreateRequest
		err error
	)
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	collection := &model.Collection{
		Name:        req.Name,
		FileId:      req.FileId,
		Description: req.Description,
	}
	c := model.NewController()
	c.Begin()
	if err = c.CollectionModel(collection).Save(); err != nil {
		logger.Error(err)
		_ = c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	c.Commit()
	api.OK(ctx, nil, 0)
}
