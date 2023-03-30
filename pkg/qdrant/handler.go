/**
 * @Author Nil
 * @Description pkg/qdrant/handler.go
 * @Date 2023/3/30 16:43
 **/

package qdrant

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"time"
)

func InitHandler(ctx *gin.Context, collectionName string) (clientSet *ClientSet, err error) {
	clientSet, err = NewClientSet(ctx)
	if err != nil {
		logger.Fatalf("get qdrant clientset error: %v", err)
	}
	err = recreateCollection(ctx, clientSet, collectionName)
	return
}

func recreateCollection(ctx context.Context, clientSet *ClientSet, collectionName string) (err error) {
	var (
		deleteInfoResp = new(pb.CollectionOperationResponse)
		createInfoResp = new(pb.CollectionOperationResponse)
	)
	var timeout uint64 = 10
	if deleteInfoResp, err = clientSet.CollectionsClient.Delete(ctx,
		&pb.DeleteCollection{
			CollectionName: collectionName,
			Timeout:        &timeout,
		}); err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Infof(
		"delete collection: %s result: %s; timing: %s;",
		collectionName,
		deleteInfoResp.Result,
		deleteInfoResp.Time,
	)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	if createInfoResp, err = clientSet.Create(ctxWithTimeout, &pb.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     1536,
					Distance: pb.Distance_Cosine,
				},
			},
		},
	}); err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Infof(
		"create collection: %s result: %s; timing: %s;",
		collectionName,
		createInfoResp.Result,
		createInfoResp.Time,
	)
	return
}

func ToEmbeddings(ctx context.Context, client *openai.Client, m *DataModel) (prompt, completion string, embedding []float32, err error) {
	var (
		resp openai.EmbeddingResponse
	)
	if resp, err = client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{m.Completion},
		Model: openai.AdaEmbeddingV2,
	}); err != nil {
		logger.Fatalf("create embeddings error: %v", err)
		return
	}

	if len(resp.Data) == 0 {
		err = errors.New("data error")
		return
	}
	prompt = m.Prompt
	completion = m.Completion
	embedding = resp.Data[0].Embedding
	return
}

func ToUpsertPoint(collectionName, prompt, completion string, embedding []float32, count uint64) *pb.UpsertPoints {
	waitUpsert := true
	return &pb.UpsertPoints{
		CollectionName: collectionName,
		Wait:           &waitUpsert,
		Points: []*pb.PointStruct{
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{
						Num: count,
					},
				},
				Vectors: &pb.Vectors{
					VectorsOptions: &pb.Vectors_Vector{
						Vector: &pb.Vector{
							Data: embedding,
						},
					},
				},
				Payload: map[string]*pb.Value{
					"title": {
						Kind: &pb.Value_StringValue{
							StringValue: prompt,
						},
					},
					"text": {
						Kind: &pb.Value_StringValue{
							StringValue: completion,
						},
					},
				},
			},
		},
		Ordering: nil,
	}
}
