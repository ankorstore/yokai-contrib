package log_test

import (
	"testing"

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

		rHandled, rRet, rErr := react.React("test")

		assert.False(t, rHandled)
		assert.Nil(t, rRet)
		assert.NoError(t, rErr)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"req":     "test",
			"message": "log reactor",
		})
	})
}
