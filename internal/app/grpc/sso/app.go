package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/Ktulhu777/protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
	cc  *grpc.ClientConn // Храним соединение для закрытия
}

func NewClient(
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const fn = "grpc.NewClient"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(interceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		log.Error("failed to create gRPC connection", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
		log: log,
		cc:  cc,
	}, nil
}

// Закрытие соединения
func (c *Client) Close() error {
	return c.cc.Close()
}

func interceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, toSlogLevel(level), msg, fields...)
	})
}

func toSlogLevel(level grpclog.Level) slog.Level {
	switch level {
	case grpclog.LevelDebug:
		return slog.LevelDebug
	case grpclog.LevelInfo:
		return slog.LevelInfo
	case grpclog.LevelWarn:
		return slog.LevelWarn
	case grpclog.LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const fn = "grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{UserId: userID})
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return resp.IsAdmin, nil
}
