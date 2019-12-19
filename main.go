package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	config "github.com/neolee/mixin-wop/config"
)

const defaultResponse = "I'm ready"

var client *bot.BlazeClient

// Handler is an implementation for interface bot.BlazeListener
// For more details: https://github.com/MixinNetwork/bot-api-go-client/blob/master/blaze.go#L89
type Handler struct{}

// OnMessage is a general method of bot.BlazeListener
func (r Handler) OnMessage(ctx context.Context, msgView bot.MessageView, botID string) error {
	// I handle PLAIN_TEXT message only and make sure respond to current conversation.
	if msgView.Category == bot.MessageCategoryPlainText &&
		msgView.ConversationId == bot.UniqueConversationId(config.GetConfig().ClientID, msgView.UserId) {
		var data []byte
		var err error
		if data, err = base64.StdEncoding.DecodeString(msgView.Data); err != nil {
			log.Panicf("Error: %s\n", err)
			return err
		}
		inst := string(data)
		log.Printf("I got a message from %s, it said: `%s`\n", msgView.UserId, inst)

		if "sync" == inst {
			// Sync? Ack!
			Respond(ctx, msgView, "ack")
		} else {
			Respond(ctx, msgView, defaultResponse)
		}
	}
	return nil
}

// Respond to user.
func Respond(ctx context.Context, msgView bot.MessageView, msg string) {
	if err := client.SendPlainText(ctx, msgView, msg); err != nil {
		log.Panicf("Error: %s\n", err)
	}
}

func main() {
	ctx := context.Background()
	log.Println("start bot")
	handler := Handler{}

	// Create a bot client
	client = bot.NewBlazeClient(config.GetConfig().ClientID, config.GetConfig().SessionID, config.GetConfig().PrivateKey)

	// Start the loop
	for {
		if err := client.Loop(ctx, handler); err != nil {
			log.Printf("Error: %v\n", err)
		}
		log.Println("connection loop end")
		time.Sleep(time.Second)
	}
}
