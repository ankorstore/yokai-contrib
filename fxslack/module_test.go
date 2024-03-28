package fxslack_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxslack"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxSlackModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxslack.FxSlackModule,
		fxconfig.FxConfigModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxSlackClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvDev)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var client *slack.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		provideTestRoundTripper(),
		fx.Populate(&conf, &client),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create slack.Client")
	assert.NotNil(t, client)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close slack.Client")
}

func TestFxSlackTestClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var client *slack.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		provideTestRoundTripper(),
		fx.Populate(&conf, &client),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test slack.Client")
	assert.NotNil(t, client)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test slack.Client")
}

func TestFxSlackTestServer(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var server *slacktest.Server

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		provideTestRoundTripper(),
		fx.Populate(&conf, &server),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test slacktest.Server")
	assert.NotNil(t, server)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test slacktest.Server")
}

func TestSlackClientWithTestServer(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var server *slacktest.Server
	var client *slack.Client

	maxWait := 1 * time.Second

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		provideTestRoundTripper(),
		fx.Populate(&conf, &server, &client),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test slacktest.Server")
	assert.NotNil(t, server)

	server.Start()
	defer server.Stop()

	rtm := client.NewRTM()
	go rtm.ManageConnection()
	rtm.SendMessage(&slack.OutgoingMessage{
		Channel: "foo",
		Text:    "should see this inbound message",
	})

	time.Sleep(maxWait)
	seenInbound := server.GetSeenInboundMessages()
	assert.True(t, len(seenInbound) > 0)
	for _, msg := range seenInbound {
		var m = slack.Message{}
		jerr := json.Unmarshal([]byte(msg), &m)
		assert.NoError(t, jerr, "messages should decode as slack.Message")
		if m.Text == "should see this inbound message" {
			break
		}
	}

	assert.True(t, server.SawMessage("should see this inbound message"))

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test slacktest.Server")
}

func provideTestRoundTripper() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() http.RoundTripper {
				return &http.Transport{}
			},
			fx.As(new(http.RoundTripper)),
		),
	)
}
