/**
 * @Author Nil
 * @Description api/image/generate.go
 * @Date 2023/4/10 20:43
 **/

package image

import (
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
)

func Generate(ctx *gin.Context) {
	prompt := ctx.Query("prompt")
	openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)
	req := openai.ImageRequest{
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}
	var (
		imageResp     openai.ImageResponse
		imageHttpResp *http.Response
		imageBytes    []byte
		err           error
	)
	if imageResp, err = openaiClient.CreateImage(ctx, req); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InternalError)
		return
	}
	if imageHttpResp, err = http.Get(imageResp.Data[0].URL); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InternalError)
		return
	}
	if imageBytes, err = io.ReadAll(imageHttpResp.Body); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InternalError)
		return
	}
	//imageResp.Data[0].B64JSON = base64.StdEncoding.EncodeToString(imageBytes)
	if err = os.WriteFile(path.Join(config.SysCache.ServerConfig.Storage, "./storage", prompt+".png"), imageBytes, fs.ModePerm); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InternalError)
		return
	}
	api.OK(ctx, imageResp, 1)
}
