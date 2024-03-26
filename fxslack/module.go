package fxslack

import (
	"context"
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
		NewSlackTestServer,
	),
)

type FxSlackClientParam struct {
	fx.In
	LifeCycle        fx.Lifecycle
	HttpRoundTripper http.RoundTripper
	Config           *config.Config
	TestServer       *slacktest.Server
}

func NewSlackClient(p FxSlackClientParam) *slack.Client {
	if p.Config.IsTestEnv() {
		return createTestClient(p)
	} else {
		return createClient(p)
	}
}

type FxSlackTestServerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
}

func NewSlackTestServer(p FxSlackTestServerParam) *slacktest.Server {
	if p.Config.IsTestEnv() {
		server := slacktest.NewTestServer()

		p.LifeCycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				go server.Start()

				return nil
			},
			OnStop: func(context.Context) error {
				server.Stop()

				return nil
			},
		})

		return server
	}

	return nil
}

func createClient(p FxSlackClientParam) *slack.Client {
	httpClient := &http.Client{
		Transport: p.HttpRoundTripper,
	}

	client := slack.New(p.Config.GetString("modules.slack.auth_token"), slack.OptionHTTPClient(httpClient), slack.OptionAppLevelToken(p.Config.GetString("modules.slack.app_level_token")))

	return client
}

func createTestClient(p FxSlackClientParam) *slack.Client {
	server := p.TestServer

	httpClient := &http.Client{
		Transport: p.HttpRoundTripper,
	}

	client := slack.New(p.Config.GetString("modules.slack.auth_token"), slack.OptionHTTPClient(httpClient), slack.OptionAppLevelToken(p.Config.GetString("modules.slack.app_level_token")), slack.OptionAPIURL(server.GetAPIURL()))

	return client
}
