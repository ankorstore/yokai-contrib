package fxslack_test

import (
	"context"
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
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var client *slack.Client

	var roundTripperProvide = func() http.RoundTripper {
		return http.DefaultTransport
	}

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		fx.Populate(&conf, &client),
		fx.Provide(roundTripperProvide),
	)

	err := app.Start(context.Background())
	assert.NoError(t, err, "failed to create slack.Client")
	assert.NotNil(t, client)

	err = app.Stop(context.Background())
	assert.NoError(t, err, "failed to close slack.Client")
}

func TestFxSlackTestClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var client *slack.Client

	var roundTripperProvide = func() http.RoundTripper {
		return http.DefaultTransport
	}

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		fx.Populate(&conf, &client),
		fx.Provide(roundTripperProvide),
	)

	err := app.Start(context.Background())
	assert.NoError(t, err, "failed to create test slack.Client")
	assert.NotNil(t, client)

	err = app.Stop(context.Background())
	assert.NoError(t, err, "failed to close test slack.Client")
}

func TestFxSlackTestServer(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var server *slacktest.Server

	var roundTripperProvide = func() http.RoundTripper {
		return http.DefaultTransport
	}

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		fx.Populate(&conf, &server),
		fx.Provide(roundTripperProvide),
	)

	err := app.Start(context.Background())
	assert.NoError(t, err, "failed to create test slacktest.Server")
	assert.NotNil(t, server)

	err = app.Stop(context.Background())
	assert.NoError(t, err, "failed to close test slacktest.Server")
}

func TestSlackClientWithTestServer(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var server *slacktest.Server
	var client *slack.Client

	var roundTripperProvide = func() http.RoundTripper {
		return http.DefaultTransport
	}

	maxWait := 5 * time.Second

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		fx.Populate(&conf, &server, &client),
		fx.Provide(roundTripperProvide),
	)

	err := app.Start(context.Background())
	assert.NoError(t, err, "failed to create test slacktest.Server")
	assert.NotNil(t, server)

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

	err = app.Stop(context.Background())
	assert.NoError(t, err, "failed to close test slacktest.Server")
}
