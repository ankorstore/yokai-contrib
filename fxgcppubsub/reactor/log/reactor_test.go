package log_test

import (
	"testing"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/log"
	yokailog "github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogReactor(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := yokailog.NewDefaultLoggerFactory().Create(
		yokailog.WithLevel(zerolog.DebugLevel),
		yokailog.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	react := log.NewLogReactor(logger)

	t.Run("func names", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			[]string{
				"GetTopic",
				"UpdateTopic",
				"ListTopics",
				"ListTopicSubscriptions",
				"DeleteTopic",
				"GetSubscription",
				"UpdateSubscription",
				"ListSubscriptions",
				"DeleteSubscription",
				"DetachSubscription",
				"CreateSchema",
				"GetSchema",
				"ListSchemas",
				"ListSchemaRevisions",
				"CommitSchema",
				"RollbackSchema",
				"DeleteSchemaRevision",
				"DeleteSchema",
				"ValidateSchema",
				"Publish",
				"Acknowledge",
				"ModifyAckDeadline",
				"Pull",
				"Seek",
				"ValidateMessage",
			},
			react.FuncNames(),
		)
	})

	t.Run("react", func(t *testing.T) {
		t.Parallel()

		req := &pubsubpb.AcknowledgeRequest{
			Subscription: "test-subscription",
			AckIds:       []string{"test-id"},
		}

		rHandled, rRet, rErr := react.React(req)

		assert.False(t, rHandled)
		assert.Nil(t, rRet)
		assert.NoError(t, rErr)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"type":    "*pubsubpb.AcknowledgeRequest",
			"data":    "map[ack_ids:[test-id] subscription:test-subscription]",
			"message": "log reactor",
		})
	})
}
