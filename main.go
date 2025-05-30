package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/ulugbek0217/startup-house/filters"
	"github.com/ulugbek0217/startup-house/handlers"
	"github.com/ulugbek0217/startup-house/misc"
)

func main() {
	err := misc.LoadEnv("config/.env")
	if err != nil {
		log.Fatalf("err loading env: %v\n", err)
	}

	var botToken string = os.Getenv("BOT_TOKEN")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		// bot.WithMessageTextHandler("привет", bot.MatchTypeContains, handlers.Greeting),
		bot.WithMessageTextHandler("ии", bot.MatchTypeContains, handlers.AnswerAI),
		bot.WithDefaultHandler(handlers.AnswerAI),
		bot.WithSkipGetMe(),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandlerMatchFunc(filters.IsGreeting, handlers.Greeting)

	fmt.Println("Listening...")

	b.Start(ctx)
}
