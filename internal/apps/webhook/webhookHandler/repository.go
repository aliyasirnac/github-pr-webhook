package webhookHandler

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookDomain"
)

type GithubWebhookRepository interface {
	SaveGithubWebhook(ctx context.Context, req *webhookDomain.Github) error
}
