package log

import (
	"github.com/ankorstore/yokai/log"
)

// LogReactor is a pub/sub test server reactor for logging server events.
type LogReactor struct {
	logger *log.Logger
}

// NewLogReactor returns a new LogReactor instance.
func NewLogReactor(logger *log.Logger) *LogReactor {
	return &LogReactor{
		logger: logger,
	}
}

// FuncNames returns the list of function names this reactor will react to.
func (r *LogReactor) FuncNames() []string {
	return []string{
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
	}
}

// React is the reactor logic.
func (r *LogReactor) React(req interface{}) (bool, any, error) {
	r.logger.Debug().Interface("req", req).Msg("log reactor")

	return false, nil, nil
}
