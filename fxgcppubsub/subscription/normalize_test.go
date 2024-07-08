package subscription_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeSubscriptionName(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		"projects/test-project/subscriptions/test-subscription",
		subscription.NormalizeSubscriptionName("test-project", "test-subscription"),
	)
}
