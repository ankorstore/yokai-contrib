package subscription

import "fmt"

// NormalizeSubscriptionName normalizes a subscription name for given projectID and subscriptionID
// to be compatible with the "projects/{projectID}/subscriptions/{subscriptionID}" format.
func NormalizeSubscriptionName(projectID string, subscriptionID string) string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subscriptionID)
}
