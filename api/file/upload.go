/**
 * @Author Nil
 * @Description api/file/upload.go
 * @Date 2023/4/2 00:41
 **/

package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/model"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"io"
	"mime/multipart"
	"net/http"
	"path"
)

func Upload(ctx *gin.Context) {
	var (
		fileHeader = new(multipart.FileHeader)
		f          multipart.File

		fileModel model.File

		//content []byte

		err error
	)
	md5Hash := md5.New()
	fileHeader, err = ctx.FormFile("file")
	dst := path.Join(config.SysCache.ServerConfig.Storage, "upload", fileHeader.Filename)
	if err = ctx.SaveUploadedFile(fileHeader, dst); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	if f, err = fileHeader.Open(); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	if f, err = fileHeader.Open(); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	if _, err = io.Copy(md5Hash, f); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	hash := hex.EncodeToString(md5Hash.Sum(nil))
	c := model.NewController()
	c.Begin()
	if fileModel, err = c.FileModel(&model.File{}).Get(&model.FileQuery{Hash: &hash}); err != nil {
		logger.Error(err)
		c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	if fileModel.ID != 0 {
		err = errors.New("file has already exists")
		logger.Error(err)
		c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.ResourceOperateError)
		return
	}
	if err = c.FileModel(&model.File{
		Name: fileHeader.Filename,
		Path: dst,
		Hash: hash,
		Type: model.PromptCompletion,
	}).Save(); err != nil {
		logger.Error(err)
		c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	c.Commit()
	api.OK(ctx, nil, 0)
}
