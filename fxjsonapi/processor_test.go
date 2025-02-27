package fxjsonapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/model"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type DynamicHandlerFunc func(p fxjsonapi.Processor, c echo.Context) error

type DynamicHandler struct {
	processor          fxjsonapi.Processor
	dynamicHandlerFunc DynamicHandlerFunc
}

func NewDynamicHandler(processor fxjsonapi.Processor, dynamicHandlerFunc DynamicHandlerFunc) *DynamicHandler {
	return &DynamicHandler{
		processor:          processor,
		dynamicHandlerFunc: dynamicHandlerFunc,
	}
}

func (h *DynamicHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		return h.dynamicHandlerFunc(h.processor, c)
	}
}

func TestProcessor(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(
		tb testing.TB,
		dynamicHandlerFunc DynamicHandlerFunc,
	) (
		*echo.Echo,
		logtest.TestLogBuffer,
		tracetest.TestTraceExporter,
	) {
		tb.Helper()

		var httpServer *echo.Echo
		var logBuffer logtest.TestLogBuffer
		var traceExporter tracetest.TestTraceExporter

		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxtrace.FxTraceModule,
			fxmetrics.FxMetricsModule,
			fxgenerate.FxGenerateModule,
			fxhttpserver.FxHttpServerModule,
			fxjsonapi.FxJSONAPIModule,
			fx.Supply(dynamicHandlerFunc),
			fxhttpserver.AsHandler("POST", "/dynamic", NewDynamicHandler),
			fx.Populate(&httpServer, &logBuffer, &traceExporter),
		).RequireStart().RequireStop()

		return httpServer, logBuffer, traceExporter
	}

	t.Run("test request processing error with invalid content type", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.Foo{}

			err := p.ProcessRequest(c, &foo)
			if err != nil {
				return err
			}

			return c.NoContent(http.StatusNoContent)
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/dynamic", nil)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code, rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"message": "JSON API request invalid content type",
		})

		span, err := traceExporter.Span("JSON API request processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Error, span.Snapshot().Status().Code)
	})

	t.Run("test request processing success with defaults", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.Foo{}

			err := p.ProcessRequest(c, &foo)
			if err != nil {
				return err
			}

			return c.NoContent(http.StatusNoContent)
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		foo := model.CreateTestFoo()

		mFoo, err := fxjsonapi.Marshall(&foo, fxjsonapi.MarshallParams{})
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/dynamic", bytes.NewBuffer(mFoo))
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API request processing success",
		})

		span, err := traceExporter.Span("JSON API request processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Ok, span.Snapshot().Status().Code)
	})
}
