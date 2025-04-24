package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/ulugbek0217/startup-house/filters"
)

func Greeting(ctx context.Context, b *bot.Bot, upd *models.Update) {
	if !filters.IsGroup(upd) {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   "Привет! Добро пожаловать в Startup House! 🚀",
		ReplyParameters: &models.ReplyParameters{
			MessageID: upd.Message.ID,
			ChatID:    upd.Message.Chat.ID,
		},
	})
}

func AnswerAI(ctx context.Context, b *bot.Bot, upd *models.Update) {
	currentTime := time.Now()
	fmt.Printf("[%d-%02d-%02dT%02d:%02d:%-2d] in answer ai\n",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	if !filters.IsGroup(upd) {
		return
	}
	currentTime = time.Now()
	if strings.EqualFold(upd.Message.Text, "/ai") {
		fmt.Printf("[%d-%02d-%02dT%02d:%02d:%-2d] got /ai command\n",
			currentTime.Year(), currentTime.Month(), currentTime.Day(),
			currentTime.Hour(), currentTime.Minute(), currentTime.Second())
		var message string
		if upd.Message.ReplyToMessage != nil {
			message = upd.Message.ReplyToMessage.Text
			response, err := AIResponse(message)
			if err != nil {
				log.Printf("err getting response: %v\n", err)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: upd.Message.Chat.ID,
					Text:   "Извините, произошла внутренняя ошибка.",
					ReplyParameters: &models.ReplyParameters{
						MessageID: upd.Message.ID,
						ChatID:    upd.Message.Chat.ID,
					},
				})
				return
			}
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upd.Message.Chat.ID,
				Text:   response,
				// ParseMode: models.ParseModeMarkdown,
				ReplyParameters: &models.ReplyParameters{
					MessageID: upd.Message.ID,
					ChatID:    upd.Message.Chat.ID,
				},
			})
			if err != nil {
				log.Printf(">>>>err sending response: %v\n", err)
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upd.Message.Chat.ID,
				Text:   "Пожалуйста, ответьте на какой нибудь текст.",
				ReplyParameters: &models.ReplyParameters{
					MessageID: upd.Message.ID,
					ChatID:    upd.Message.Chat.ID,
				},
			})
		}
	}

	if filters.IsAboutAI(upd) {
		interesting, err := AIResponse(fmt.Sprintf("Задавай короткий вопрос, добавь эмодзи в свой вопрос:\n%s",
			upd.Message.Text))
		if err != nil {
			log.Printf("err interesting in ai: %v", err)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			Text:   interesting,
			ReplyParameters: &models.ReplyParameters{
				MessageID: upd.Message.ID,
				ChatID:    upd.Message.Chat.ID,
			},
		})
	} else {
		return
	}
}

func AIResponse(message string) (string, error) {
	var (
		token = os.Getenv("OPEN_AI_TOKEN")
		url   = os.Getenv("URL")
	)

	body := map[string]any{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": message + " \nwrite without formatting characters",
			},
		},
		"max_tokens": 512,
		"stream":     false,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	bodyContent := make(map[string]interface{})
	json.Unmarshal(responseBody, &bodyContent)

	var prettyJson bytes.Buffer
	err = json.Indent(&prettyJson, responseBody, "", "  ")
	if err != nil {
		fmt.Println("Failed to format JSON:", err)
		return "", err
	}

	if choices, ok := bodyContent["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					// fmt.Println("Content:", content)
					return content, nil
				}
			}
		}
	}
	fmt.Println(prettyJson.String())
	return "Извините, произошла внутренняя ошибка.", err
}
