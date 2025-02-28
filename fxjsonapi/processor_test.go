package fxjsonapi_test

import (
	"bytes"
	"fmt"
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

//nolint:maintidx
func TestProcessor(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(
		tb testing.TB,
		dynamicHandlerFunc handler.DynamicHandlerFunc,
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
			fxhttpserver.AsHandler("POST", "/test", handler.NewDynamicHandler),
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
		req := httptest.NewRequest(http.MethodPost, "/test", nil)

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

	t.Run("test request processing error with invalid json api data", func(t *testing.T) {
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
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"message": "JSON API request processing error",
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
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(mFoo))
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

	t.Run("test request processing success without logging and tracing", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.Foo{}

			err := p.ProcessRequest(
				c,
				&foo,
				fxjsonapi.WithLog(false),
				fxjsonapi.WithTrace(false),
			)
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
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(mFoo))
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())

		logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API request processing success",
		})

		tracetest.AssertHasNotTraceSpan(t, traceExporter, "JSON API request processing")
	})

	t.Run("test response processing error with invalid data", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			type invalid struct {
				Name string `jsonapi:"invalid"`
			}

			return p.ProcessResponse(c, http.StatusOK, &invalid{})
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "error",
			"error":   "Bad jsonapi struct tag format",
			"message": "JSON API response processing error",
		})

		span, err := traceExporter.Span("JSON API response processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Error, span.Snapshot().Status().Code)
	})

	t.Run("test response processing success with defaults", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.CreateTestFoo()

			return p.ProcessResponse(c, http.StatusOK, &foo)
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}},"included":[{"type":"bar","id":"456","attributes":{"name":"bar"},"meta":{"meta":"bar"}}]}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API response processing success",
		})

		span, err := traceExporter.Span("JSON API response processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Ok, span.Snapshot().Status().Code)
	})

	t.Run("test response processing success without included", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.CreateTestFoo()

			return p.ProcessResponse(
				c,
				http.StatusOK,
				&foo,
				fxjsonapi.WithIncluded(false),
			)
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}}}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API response processing success",
		})

		span, err := traceExporter.Span("JSON API response processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Ok, span.Snapshot().Status().Code)
	})

	t.Run("test response processing success with metadata", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.CreateTestFoo()

			return p.ProcessResponse(
				c,
				http.StatusOK,
				&foo,
				fxjsonapi.WithMetadata(map[string]interface{}{
					"baz": "buz",
				}),
			)
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}},"included":[{"type":"bar","id":"456","attributes":{"name":"bar"},"meta":{"meta":"bar"}}],"meta":{"baz":"buz"}}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API response processing success",
		})

		span, err := traceExporter.Span("JSON API response processing")
		assert.NoError(t, err)
		assert.Equal(t, codes.Ok, span.Snapshot().Status().Code)
	})

	t.Run("test response processing success without logging and tracing", func(t *testing.T) {
		fn := func(p fxjsonapi.Processor, c echo.Context) error {
			foo := model.CreateTestFoo()

			return p.ProcessResponse(
				c,
				http.StatusOK,
				&foo,
				fxjsonapi.WithLog(false),
				fxjsonapi.WithTrace(false),
			)
		}

		httpServer, logBuffer, traceExporter := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}},"included":[{"type":"bar","id":"456","attributes":{"name":"bar"},"meta":{"meta":"bar"}}]}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "debug",
			"message": "JSON API response processing success",
		})

		tracetest.AssertHasNotTraceSpan(t, traceExporter, "JSON API response processing")
	})
}
