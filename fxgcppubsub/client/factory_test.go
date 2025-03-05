package client_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub/pstest"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/client"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDefaultClientFactory(t *testing.T) {
	t.Parallel()

	createFactory := func(tb testing.TB) (*client.DefaultClientFactory, logtest.TestLogBuffer) {
		tb.Helper()

		cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("../testdata/config"))
		assert.NoError(tb, err)

		logBuffer := logtest.NewDefaultTestLogBuffer()
		logger, err := log.NewDefaultLoggerFactory().Create(
			log.WithLevel(zerolog.DebugLevel),
			log.WithOutputWriter(logBuffer),
		)
		assert.NoError(tb, err)

		return client.NewDefaultClientFactory(cfg, logger), logBuffer
	}

	t.Run("creation success", func(t *testing.T) {
		t.Parallel()

		factory, logBuffer := createFactory(t)

		server := pstest.NewServer()
		defer server.Close()

		psClient, psErr := factory.Create(
			context.Background(),
			"test-project",
			option.WithEndpoint(server.Addr),
			option.WithoutAuthentication(),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)

		assert.NoError(t, psErr)
		assert.NotNil(t, psClient)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"attempt": 1,
			"message": "pubsub client creation success",
		})
	})

	t.Run("creation error with retries", func(t *testing.T) {
		t.Parallel()

		factory, logBuffer := createFactory(t)

		server := pstest.NewServer()
		defer server.Close()

		psClient, psErr := factory.Create(
			context.Background(),
			"test-project",
			option.WithEndpoint(server.Addr),
			// option.WithoutAuthentication(), <- removed to make the client fail
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)

		assert.Error(t, psErr)
		assert.Equal(t, "pubsub client creation error after 3 attempts", psErr.Error())
		assert.Nil(t, psClient)

		logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "warn",
			"attempt": 1,
			"error":   "grpc: the credentials require transport level security",
			"message": "pubsub client creation error, attempting again in 1 seconds",
		})

		logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "warn",
			"attempt": 2,
			"error":   "grpc: the credentials require transport level security",
			"message": "pubsub client creation error, attempting again in 1 seconds",
		})

		logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "warn",
			"attempt": 3,
			"error":   "grpc: the credentials require transport level security",
			"message": "pubsub client creation error, attempting again in 1 seconds",
		})
	})
}
