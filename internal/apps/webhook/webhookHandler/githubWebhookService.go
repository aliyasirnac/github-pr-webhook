package webhookHandler

import (
	"context"
	"go.uber.org/zap"
)

type GithubWebhookRequest struct{}

type GithubWebhookResponse struct {
	Id int64 `json:"id"`
}

type GithubWebhookHandler struct{}

func NewGithubWebhookHandler() *GithubWebhookHandler {
	return &GithubWebhookHandler{}
}

func (w *GithubWebhookHandler) Handle(ctx context.Context, req *GithubWebhookRequest) (*GithubWebhookResponse, error) {
	zap.L().Info("githubWebhookHandler.Handle", zap.Any("req", req))
	return &GithubWebhookResponse{Id: 2}, nil
}
