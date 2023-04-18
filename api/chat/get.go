/**
 * @Author Nil
 * @Description api/chat/get.go
 * @Date 2023/4/10 19:24
 **/

package chat

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/api"
	"github.com/ha5ky/hu5ky-bot/model"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	boterrors "github.com/ha5ky/hu5ky-bot/pkg/errors"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"time"
)

func Get(ctx *gin.Context) {
	prompt := ctx.Query("prompt")
	var (
		err                  error
		completionStreamResp *openai.ChatCompletionStream
		//completionResp openai.ChatCompletionResponse
		messages  = make([]*model.Message, 0)
		pageSize  = 10
		pageIndex = 0
	)
	system := "你是hu5ky智能助手"
	c := model.NewController()
	if messages, _, err = c.MessageModel(&model.Message{}).List(&model.MessageQuery{
		PageSize:  &pageSize,
		PageIndex: &pageIndex,
		Desc:      true,
	}); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	promptMsg := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: system,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
	openaiClient := openai.NewClient(config.SysCache.GPT.OpenaiAPIKey)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    convertToChatCompletion(promptMsg, messages...),
		MaxTokens:   2000,
		Temperature: 0.7,
		TopP:        1,
		Stream:      true,
	}

	if completionStreamResp, err = openaiClient.CreateChatCompletionStream(ctx, req); err != nil {
		logger.Error(err)
		api.ErrorRender(ctx, http.StatusBadRequest, err, boterrors.InvalidParams)
		return
	}
	defer completionStreamResp.Close()
	var (
		recvChan = make(chan string, 1)
		errChan  = make(chan error, 1)

		completion string
	)

	go func() {
		for {
			time.Sleep(time.Microsecond * 100)
			recv, err := completionStreamResp.Recv()
			if errors.Is(err, io.EOF) {
				errChan <- err
				logger.Info("\nStream finished")
				return
			}
			recvChan <- recv.Choices[0].Delta.Content

			if err != nil {
				errChan <- err
				logger.Infof("\nStream error: %v\n", err)
				return
			}

			completion += recv.Choices[0].Delta.Content
		}
	}()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case msg := <-recvChan:
			ctx.SSEvent("completion", msg)
		case <-errChan:
			if err = c.MessageModel(&model.Message{
				Prompt:     prompt,
				Completion: completion,
			}).Save(); err != nil {
				logger.Error(err)
				ctx.SSEvent("error: ", err.Error())
			}
		}
		return true
	})
}

func convertToChatCompletion(promptMsg []openai.ChatCompletionMessage, messages ...*model.Message) []openai.ChatCompletionMessage {
	if len(messages) >= 3 {
		messages = messages[:3]
	}

	for i := len(messages) - 1; i >= 0; i-- {
		promptMsg = append(
			promptMsg,
			openai.ChatCompletionMessage{
				Role:    "user",
				Content: messages[i].Prompt,
			},
			openai.ChatCompletionMessage{
				Role:    "assistant",
				Content: messages[i].Completion,
			},
		)
	}
	return promptMsg
}
