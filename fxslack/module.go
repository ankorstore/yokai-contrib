package fxslack

import (
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "slack"

// FxSlack is the [Fx] slack module.
//
// [Fx]: https://github.com/uber-go/fx
var FxSlackModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewSlackClient,
	),
)

// FxSlackClientParam allows injection of the required dependencies in [NewSlackClient].
type FxSlackClientParam struct {
	fx.In
	LifeCycle        fx.Lifecycle
	HttpRoundTripper http.RoundTripper
	Config           *config.Config
}

// NewSlackClient returns a [slack.Client].
func NewSlackClient(p FxSlackClientParam) *slack.Client {
	if p.Config.IsTestEnv() {
		return createTestClient(p)
	} else {
		return createClient(p)
	}
}

func createClient(p FxSlackClientParam) *slack.Client {
	httpClient := &http.Client{
		Transport: p.HttpRoundTripper,
	}

	client := slack.New(p.Config.GetString("modules.slack.token"), slack.OptionHTTPClient(httpClient))

	return client
}

func createTestClient(p FxSlackClientParam) *slack.Client {
	server := slacktest.NewTestServer()

	go server.Start()

	httpClient := &http.Client{
		Transport: p.HttpRoundTripper,
	}

	client := slack.New(p.Config.GetString("modules.slack.token"), slack.OptionHTTPClient(httpClient), slack.OptionAPIURL(server.GetAPIURL()))

	return client
}
