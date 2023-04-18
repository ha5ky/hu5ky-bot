/**
 * @Author Nil
 * @Description api/completion/get.go
 * @Date 2023/3/30 10:06
 **/

package completion

import (
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/model"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/ha5ky/hu5ky-bot/pkg/qdrant"
	"github.com/ha5ky/hu5ky-bot/pkg/util"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"strconv"
)

func Get(ctx *gin.Context) {
	search := ctx.Query("search")
	collectionIdStr := ctx.Query("collection_id")
	var (
		c = model.NewController()

		collectionId int

		collection model.Collection
		answer     *queryResp

		err error
	)
	c.Begin()
	if collectionId, err = strconv.Atoi(collectionIdStr); err != nil {
		logger.Fatalf("to embeddings error: %v", err)
		c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	uintCollectionId := uint(collectionId)
	if collection, err = c.CollectionModel(&model.Collection{}).Get(&model.CollectionQuery{ID: &uintCollectionId}); err != nil {
		logger.Error(err)
		c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	if answer, err = query(ctx, search, collection.Name); err != nil {
		logger.Error(err)
		c.Rollback()
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	c.Commit()
	api.OK(ctx, answer, 1)
}

func prompt(question string, answers []*Answer) (ret []openai.ChatCompletionMessage) {
	system := "你是hu5ky智能助手"
	q := "使用以下段落来回答问题，如果段落内容不相关就返回未查到相关信息：\""
	q += question + "\"\n"
	for index, answer := range answers {
		q += strconv.Itoa(index+1) + ". " + answer.Title + ": " + answer.Text + "\n"
	}
	ret = []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: system,
		},
		{
			Role:    "user",
			Content: q,
		},
	}
	return
}

type Answer struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type queryResp struct {
	Answer string   `json:"answer"`
	Tags   []string `json:"tags"`
}

func query(ctx *gin.Context, text, collectionName string) (ret *queryResp, err error) {
	var (
		embedding      []float32
		searchResp     = new(pb.SearchResponse)
		completionResp openai.ChatCompletionResponse
		answers        = make([]*Answer, 0)
		tags           = make([]string, 0)
	)
	ret = new(queryResp)
	qdrantClientSet, err := qdrant.NewClientSet(ctx)
	if err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	defer qdrantClientSet.ConnClose()
	m := &qdrant.DataModel{
		Prompt:     "",
		Completion: text,
	}
	openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)
	if _, _, embedding, err = qdrant.ToEmbeddings(ctx, openaiClient, m); err != nil {
		logger.Fatalf("to embeddings error: %v", err)
		return
	}

	hnswEf := uint64(128)
	exact := false

	if searchResp, err = qdrantClientSet.PointsClient.Search(ctx, &pb.SearchPoints{
		CollectionName: collectionName,
		Vector:         embedding,
		Limit:          20,
		WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{
			Enable: true,
		}},
		Params: &pb.SearchParams{
			HnswEf: &hnswEf,
			Exact:  &exact,
		},
	}); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}

	for _, item := range searchResp.Result {
		answers = append(answers, &Answer{
			Title: item.Payload["title"].String(),
			Text:  item.Payload["text"].String(),
		})
	}
	promptMessage := prompt(text, answers)
	util.DumpPretty(promptMessage)

	if completionResp, err = openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    promptMessage,
		Temperature: 0.7,
	}); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	answer := ""
	if len(completionResp.Choices) > 0 {
		answer = completionResp.Choices[0].Message.Content
	}
	ret = &queryResp{
		Answer: answer,
		Tags:   tags,
	}
	return
}
