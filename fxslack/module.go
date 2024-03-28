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
		NewSlackTestServer,
		NewSlackClient,
	),
)

// FxSlackTestServerParam allows injection of the required dependencies in [NewSlackTestServer].
type FxSlackTestServerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
}

// NewSlackTestServer returns a [slacktest.Server].
func NewSlackTestServer(p FxSlackTestServerParam) *slacktest.Server {
	if p.Config.IsTestEnv() {
		return slacktest.NewTestServer()
	}

	return nil
}

// FxSlackClientParam allows injection of the required dependencies in [NewSlackClient].
type FxSlackClientParam struct {
	fx.In
	LifeCycle        fx.Lifecycle
	HttpRoundTripper http.RoundTripper
	Config           *config.Config
	TestServer       *slacktest.Server
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

	client := slack.New(
		p.Config.GetString("modules.slack.auth_token"),
		slack.OptionHTTPClient(httpClient),
		slack.OptionAppLevelToken(p.Config.GetString("modules.slack.app_level_token")),
		slack.OptionDebug(p.Config.AppDebug()),
	)

	return client
}

func createTestClient(p FxSlackClientParam) *slack.Client {
	server := p.TestServer

	httpClient := &http.Client{
		Transport: p.HttpRoundTripper,
	}

	client := slack.New(
		p.Config.GetString("modules.slack.auth_token"),
		slack.OptionHTTPClient(httpClient),
		slack.OptionAppLevelToken(p.Config.GetString("modules.slack.app_level_token")),
		slack.OptionAPIURL(server.GetAPIURL()),
		slack.OptionDebug(p.Config.AppDebug()),
	)

	return client
}
