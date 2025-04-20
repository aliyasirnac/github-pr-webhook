package webhookHandler

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/pubsubinterface"
	"go.uber.org/zap"
)

type GithubWebhookRequest struct {
	Event        string `header:"X-GitHub-Event"`
	Signature    string `header:"X-Hub-Signature"`
	Signature256 string `header:"X-Hub-Signature-256"`
}

type GithubWebhookResponse struct {
	Id int64 `json:"id"`
}

type GithubWebhookHandler struct {
	PubSub pubsubinterface.PubSub
}

func NewGithubWebhookHandler(sub pubsubinterface.PubSub) *GithubWebhookHandler {
	return &GithubWebhookHandler{
		PubSub: sub,
	}
}

func (w *GithubWebhookHandler) Handle(ctx context.Context, req *GithubWebhookRequest) (*GithubWebhookResponse, error) {
	zap.L().Info("githubWebhookHandler.Handle", zap.Any("req", req))
	err := w.PubSub.Publish("bot", []byte("deneme"))
	if err != nil {
		zap.L().Error("githubWebhookHandler.Handle", zap.Error(err))
		return nil, err
	}
	return &GithubWebhookResponse{Id: 2}, nil
}
