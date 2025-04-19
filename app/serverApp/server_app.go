package serverapp

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookConfig"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookHandler"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

import (
	_ "github.com/aliyasirnac/github-pr-webhook-bot/internal/infra/postgres"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"log"
)

var httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration_seconds",
	Help:    "Duration of HTTP requests in seconds",
	Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
}, []string{"route", "method", "status"})

func init() {
	prometheus.MustRegister(httpRequestDuration)
}

type ServerApp struct {
	Config *webhookConfig.ServerConfig
	app    *fiber.App
}

func New(config *webhookConfig.ServerConfig) *ServerApp {
	return &ServerApp{
		Config: config,
	}
}

func (s *ServerApp) Start(ctx context.Context) error {
	zap.L().Info("Starting database")
	_ = initTracer(s.Config)
	//_, err := postgres.New(s.Config.Postgres.Dsn)
	//if err != nil {
	//	zap.L().Fatal("Failed to connect to database", zap.Error(err))
	//	return err
	//}

	client := httpc()
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = client.Transport
	retryClient.RetryMax = 0
	retryClient.RetryWaitMin = 100 * time.Millisecond
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	githubWebhookHandler := webhookHandler.NewGithubWebhookHandler(retryClient)
	zap.L().Info("Server start")

	errCh := make(chan error)
	app := fiber.New(fiber.Config{})
	s.app = app
	app.Use(otelfiber.Middleware())
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Post("/webhook", handle(githubWebhookHandler))

	go func() {
		err := app.Listen(":8080")
		if err != nil {
			zap.L().Error("Error starting server", zap.Error(err))
			errCh <- err
		}
	}()
	err := <-errCh
	if err != nil {
		zap.L().Error("Server start failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *ServerApp) Stop(ctx context.Context) error {
	//close database
	if err := s.app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Server shutdown failed", zap.Error(err))
		return err
	}

	zap.L().Info("Server shutdown")
	return nil
}

func httpc() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(transport),
	}

	return httpClient
}

func initTracer(appConfig *webhookConfig.ServerConfig) *sdktrace.TracerProvider {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(appConfig.OpenTelemetry.OtelTraceEndpoint),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("webhook-go"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
