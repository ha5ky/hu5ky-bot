/**
 * @Author Nil
 * @Description api/embedding/get.go
 * @Date 2023/3/28 20:32
 **/

package embedding

import (
	"bufio"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/model"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/ha5ky/hu5ky-bot/pkg/qdrant"
	"github.com/ha5ky/hu5ky-bot/pkg/util"
	pb "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"os"
	"time"
)

type OutPutModel struct {
	qdrant.DataModel `json:",inline"`
	Result           *pb.UpdateResult `json:"result"`
	Timing           float64          `json:"timing"`
}

type CreateRequest struct {
	FileId       uint `json:"file_id" form:"file_id"`
	CollectionId uint `json:"collection_id" form:"collection_id"`
}

func Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	var (
		collection model.Collection
		err        error
	)
	c := model.NewController()
	if collection, err = c.CollectionModel(&model.Collection{}).Get(&model.CollectionQuery{ID: &req.CollectionId}); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	var (
		collectionName  = collection.Name
		qdrantClientSet *qdrant.ClientSet
		f               model.File
	)
	qdrantClientSet, err = qdrant.InitHandler(ctx, collectionName)
	if err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	defer qdrantClientSet.ConnClose()
	if f, err = c.FileModel(&model.File{}).Get(&model.FileQuery{ID: &req.FileId}); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	//openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)

	fileHandler, err := os.OpenFile(f.Path, os.O_RDONLY, 0666)
	if err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}

	defer fileHandler.Close()

	reader := bufio.NewReader(fileHandler)
	count := 0
	// 按行处理txt
	for {
		var (
			prompt, completion string
			embedding          []float32
			line               []byte
			pointsResp         = new(pb.PointsOperationResponse)
		)
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}
		count++
		m := new(qdrant.DataModel)
		if err = json.Unmarshal(line, m); err != nil {
			logger.Error(err)
			api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
			return
		}
		if m.Completion == "" {
			continue
		}
		openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)
		if prompt, completion, embedding, err = qdrant.ToEmbeddings(ctx, openaiClient, m); err != nil {
			logger.Error(err)
			api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
			return
		}

		// Upsert points
		if pointsResp, err = qdrantClientSet.PointsClient.Upsert(
			ctx,
			qdrant.ToUpsertPoint(collectionName, prompt, completion, embedding, uint64(count))); err != nil {
			logger.Error(err)
			api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
			return
		}
		util.DumpPretty(OutPutModel{
			DataModel: *m,
			Result:    pointsResp.Result,
			Timing:    pointsResp.Time,
		})
		time.Sleep(time.Second)
	}
	api.OK(ctx, nil, 0)

}
