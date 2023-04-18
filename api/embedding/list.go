/**
 * @Author Nil
 * @Description api/embedding/list.go
 * @Date 2023/3/30 17:37
 **/

package embedding

import (
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/model"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/ha5ky/hu5ky-bot/pkg/qdrant"
	pb "github.com/qdrant/go-client/qdrant"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {
	collectionIdStr := ctx.Query("collection_id")
	var (
		searchResp   = new(pb.ScrollResponse)
		collectionId int

		collection model.Collection

		err error
	)
	c := model.NewController()
	if collectionId, err = strconv.Atoi(collectionIdStr); err != nil {
		logger.Fatalf("to embeddings error: %v", err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	uintCollectionId := uint(collectionId)
	if collection, err = c.CollectionModel(&model.Collection{}).Get(&model.CollectionQuery{ID: &uintCollectionId}); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	collectionName := collection.Name
	qdrantClientSet, err := qdrant.NewClientSet(ctx)
	if err != nil {
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	//withPayload:=true
	if searchResp, err = qdrantClientSet.PointsClient.Scroll(ctx, &pb.ScrollPoints{
		CollectionName: collectionName,
		WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{
			Enable: true,
		}},
		WithVectors: &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{
			Enable: false,
		}},
	}); err != nil {
		logger.Errorf("search error: %v", err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	api.OK(ctx, searchResp, len(searchResp.Result))
}
