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
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/ha5ky/hu5ky-bot/pkg/qdrant"
	"github.com/ha5ky/hu5ky-bot/pkg/util"
	pb "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"time"
)

type OutPutModel struct {
	qdrant.DataModel `json:",inline"`
	Result           *pb.UpdateResult `json:"result"`
	Timing           float64          `json:"timing"`
}

func Create(ctx *gin.Context) {
	var (
		collectionName = config.SysCache.DB.QdRant.CollectionName
	)
	qdrantClientSet, err := qdrant.InitHandler(ctx, collectionName)
	if err != nil {
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	defer qdrantClientSet.ConnClose()
	trainingPath := config.SysCache.ServerConfig.Storage + "/training"
	fileSystem := os.DirFS(trainingPath)
	//openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)

	if err = fs.WalkDir(fileSystem, ".", func(fileName string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		fileHandler, err := os.OpenFile(path.Join(trainingPath, fileName), os.O_RDONLY, 0666)
		if err != nil {
			return err
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
				return err
			}
			if m.Completion == "" {
				continue
			}
			openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)
			if prompt, completion, embedding, err = qdrant.ToEmbeddings(ctx, openaiClient, m); err != nil {
				logger.Fatalf("to embeddings error: %v", err)
				return err
			}

			// Upsert points
			if pointsResp, err = qdrantClientSet.PointsClient.Upsert(
				ctx,
				qdrant.ToUpsertPoint(collectionName, prompt, completion, embedding, uint64(count))); err != nil {
				logger.Fatalf("upsert error: %v", err)
				return err
			}
			util.DumpPretty(OutPutModel{
				DataModel: *m,
				Result:    pointsResp.Result,
				Timing:    pointsResp.Time,
			})
			time.Sleep(time.Millisecond * 500)
		}
		return nil
	}); err != nil {
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
}
