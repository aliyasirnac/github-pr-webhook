package bot

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/bot/botconfig"
	bot2 "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type Bot struct {
	Config botconfig.Telegram
	bot    *bot2.Bot
}

func NewBot(config botconfig.Telegram) *Bot {
	return &Bot{Config: config}
}

func (b *Bot) Start(ctx context.Context) error {
	opts := []bot2.Option{
		bot2.WithDefaultHandler(handler),
	}

	bot, err := bot2.New(b.Config.ApiKey, opts...)
	if err != nil {
		zap.L().Error("Bot could not started", zap.Error(err))
		return err
	}

	b.bot = bot
	b.bot.Start(ctx)

	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	return nil
}

func (b *Bot) SendMessage(ctx context.Context, message string) {
	if b.bot == nil {
		zap.L().Error("Bot is not initialized")
		return
	}

	_, err := b.bot.SendMessage(ctx, &bot2.SendMessageParams{
		ChatID: b.Config.ChannelId,
		Text:   message,
	})
	if err != nil {
		zap.L().Error("Failed to send message", zap.Error(err))
	}
}

func handler(ctx context.Context, b *bot2.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	switch update.Message.Text {
	case "/ping":
		_, err := b.SendMessage(ctx, &bot2.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "pong 🏓",
		})
		if err != nil {
			zap.L().Error("Failed to handle /ping command", zap.Error(err))
		}
	default:
		_, err := b.SendMessage(ctx, &bot2.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Komutu anlayamadım. 🤖",
		})
		if err != nil {
			zap.L().Error("Failed to handle unknown command", zap.Error(err))
		}
	}
}
