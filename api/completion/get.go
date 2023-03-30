/**
 * @Author Nil
 * @Description api/completion/get.go
 * @Date 2023/3/30 10:06
 **/

package completion

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/ha5ky/hu5ky-bot/pkg/qdrant"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"strconv"
)

func Get(ctx *gin.Context) {
	search := ctx.Query("search")
	answer, err := query(ctx, search)
	if err != nil {
		logger.Fatalf("to embeddings error: %v", err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}

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

func query(ctx *gin.Context, text string) (ret *queryResp, err error) {
	var (
		collectionName = config.SysCache.DB.QdRant.CollectionName
		embedding      []float32
		searchResp     = new(pb.SearchResponse)
		completionResp openai.ChatCompletionResponse
		answers        = make([]*Answer, 0)
		tags           = make([]string, 0)
	)
	ret = new(queryResp)
	qdrantClientSet, err := qdrant.NewClientSet(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer qdrantClientSet.ConnClose()
	//m := &qdrant.DataModel{
	//	Prompt:     "",
	//	Completion: text,
	//}
	openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)
	//if _, _, embedding, err = qdrant.ToEmbeddings(ctx, openaiClient, m); err != nil {
	//	logger.Fatalf("to embeddings error: %v", err)
	//	return
	//}
	var (
		resp openai.EmbeddingResponse
	)
	if resp, err = openaiClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	}); err != nil {
		logger.Fatalf("create embeddings error: %v", err)
		return
	}

	if len(resp.Data) == 0 {
		err = errors.New("data error")
		return
	}
	embedding = resp.Data[0].Embedding

	hnswEf := uint64(128)
	exact := false

	if searchResp, err = qdrantClientSet.PointsClient.Search(ctx, &pb.SearchPoints{
		CollectionName: collectionName,
		Vector:         embedding,
		Limit:          3,
		WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{
			Enable: true,
		}},
		Params: &pb.SearchParams{
			HnswEf: &hnswEf,
			Exact:  &exact,
		},
	}); err != nil {
		logger.Errorf("to embeddings error: %v", err)
		return
	}

	for _, item := range searchResp.Result {
		answers = append(answers, &Answer{
			Title: item.Payload["title"].String(),
			Text:  item.Payload["text"].String(),
		})
	}
	promptMessage := prompt(text, answers)

	if completionResp, err = openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    promptMessage,
		Temperature: 0.7,
	}); err != nil {
		logger.Errorf("create completion error: %v", err)
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
