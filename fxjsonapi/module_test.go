package fxjsonapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/handler"
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

func TestFxJSONAPIModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(tb testing.TB) (*echo.Echo, logtest.TestLogBuffer, tracetest.TestTraceExporter) {
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
			fxhttpserver.AsHandler("POST", "/foobar", handler.NewFooBarHandler),
			fx.Populate(&httpServer, &logBuffer, &traceExporter),
		).RequireStart().RequireStop()

		return httpServer, logBuffer, traceExporter
	}

	t.Run("test failure with unsupported content type", func(t *testing.T) {
		httpServer, logBuffer, traceExporter := runTest(t)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/foobar", nil)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code, rec.Body.String())
		assert.Contains(t, rec.Body.String(), "Unsupported Media Type")

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"message": "JSON API request invalid content type",
		})

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"code":    http.StatusUnsupportedMediaType,
			"error":   "code=415, message=JSON API request invalid content type",
			"message": "json api error handler",
		})

		span, err := traceExporter.Span("JSON API request processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Error, span.Snapshot().Status().Code)
	})

	t.Run("test failure with invalid json api request body", func(t *testing.T) {
		httpServer, logBuffer, traceExporter := runTest(t)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/foobar", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
		assert.Contains(t, rec.Body.String(), "Bad Request")

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"message": "JSON API request processing error",
		})

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"code":    http.StatusBadRequest,
			"error":   "code=400, message=invalid character 'i' looking for beginning of value",
			"message": "json api error handler",
		})

		span, err := traceExporter.Span("JSON API request processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Error, span.Snapshot().Status().Code)
	})

	t.Run("test success with logging and tracing enabled", func(t *testing.T) {
		httpServer, logBuffer, traceExporter := runTest(t)

		foo := model.CreateTestFoo()

		mFoo, err := fxjsonapi.Marshall(&foo, fxjsonapi.MarshallParams{})
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/foobar", bytes.NewBuffer(mFoo))
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
		assert.Equal(t, rec.Body.String(), string(mFoo))

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API request processing success",
		})

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API response processing success",
		})

		reqSpan, err := traceExporter.Span("JSON API request processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Ok, reqSpan.Snapshot().Status().Code)

		respSpan, err := traceExporter.Span("JSON API response processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Ok, respSpan.Snapshot().Status().Code)
	})

	t.Run("test success with logging and tracing disabled", func(t *testing.T) {
		t.Setenv("MODULES_JSONAPI_LOG_ENABLED", "false")
		t.Setenv("MODULES_JSONAPI_TRACE_ENABLED", "false")

		httpServer, logBuffer, traceExporter := runTest(t)

		foo := model.CreateTestFoo()

		mFoo, err := fxjsonapi.Marshall(&foo, fxjsonapi.MarshallParams{})
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/foobar", bytes.NewBuffer(mFoo))
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
		assert.Equal(t, rec.Body.String(), string(mFoo))

		logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API request processing success",
		})

		logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API response processing success",
		})

		tracetest.AssertHasNotTraceSpan(t, traceExporter, "JSON API request processing")
		tracetest.AssertHasNotTraceSpan(t, traceExporter, "JSON API response processing")
	})
}
