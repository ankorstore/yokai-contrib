package fxgcppubsub_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxGcpPubSubModule(t *testing.T) {
	ctx := context.Background()

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
	)

	app.RequireStart()
	assert.NoError(t, app.Err())

	app.RequireStop()
	assert.NoError(t, app.Err())
}
