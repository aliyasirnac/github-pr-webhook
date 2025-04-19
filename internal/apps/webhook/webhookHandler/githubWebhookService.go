package webhookHandler

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

type GithubWebhookRequest struct{}

type GithubWebhookResponse struct {
	Id int64 `json:"id"`
}

type GithubWebhookHandler struct {
	RetryClient *retryablehttp.Client
}

func NewGithubWebhookHandler(retryClient *retryablehttp.Client) *GithubWebhookHandler {
	return &GithubWebhookHandler{
		RetryClient: retryClient,
	}
}

func (w *GithubWebhookHandler) Handle(ctx context.Context, req *GithubWebhookRequest) (*GithubWebhookResponse, error) {
	zap.L().Info("githubWebhookHandler.Handle", zap.Any("req", req))
	return &GithubWebhookResponse{Id: 2}, nil
}
