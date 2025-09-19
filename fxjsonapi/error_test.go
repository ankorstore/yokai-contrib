package fxjsonapi_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/handler"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

//nolint:maintidx
func TestErrorHandler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(
		tb testing.TB,
		dynamicHandlerFunc handler.DynamicHandlerFunc,
	) (
		*echo.Echo,
		logtest.TestLogBuffer,
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

		return httpServer, logBuffer
	}

	t.Run("test json api error handling", func(t *testing.T) {
		fn := func(fxjsonapi.Processor, echo.Context) error {
			return &jsonapi.ErrorObject{
				ID:     "error-id",
				Title:  "error-title",
				Detail: "error-detail",
			}
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"error-title","detail":"error-detail"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      500,
			"error":     "Error: error-title error-detail",
			"errors":    "[map[detail:error-detail id:request-id title:error-title]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test json api error handling with invalid status", func(t *testing.T) {
		fn := func(fxjsonapi.Processor, echo.Context) error {
			return &jsonapi.ErrorObject{
				ID:     "error-id",
				Title:  "error-title",
				Detail: "error-detail",
				Status: "invalid-status",
			}
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"error-title","detail":"error-detail","status":"invalid-status"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      500,
			"error":     "Error: error-title error-detail",
			"errors":    "[map[detail:error-detail id:request-id status:invalid-status title:error-title]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test json api error handling with obfuscation", func(t *testing.T) {
		t.Setenv("MODULES_HTTP_SERVER_ERRORS_OBFUSCATE", "true")

		fn := func(fxjsonapi.Processor, echo.Context) error {
			return &jsonapi.ErrorObject{
				ID:     "error-id",
				Title:  "error-title",
				Detail: "error-detail",
			}
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"error-title","detail":"Internal Server Error"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      500,
			"error":     "Error: error-title error-detail",
			"errors":    "[map[detail:Internal Server Error id:request-id title:error-title]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test validator error handling", func(t *testing.T) {
		fn := func(fxjsonapi.Processor, echo.Context) error {
			return validator.ValidationErrors{}
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

		expected := `{"errors":[]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      400,
			"error":     "",
			"errors":    "[]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test validator error handling with obfuscation", func(t *testing.T) {
		t.Setenv("MODULES_HTTP_SERVER_ERRORS_OBFUSCATE", "true")

		fn := func(fxjsonapi.Processor, echo.Context) error {
			return validator.ValidationErrors{}
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())

		expected := `{"errors":[]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      400,
			"error":     "",
			"errors":    "[]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test http error handling", func(t *testing.T) {
		fn := func(fxjsonapi.Processor, echo.Context) error {
			return echo.NewHTTPError(http.StatusBadGateway, "test-error")
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadGateway, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"Bad Gateway","detail":"code=502, message=test-error","status":"502","code":"502"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      502,
			"error":     "code=502, message=test-error",
			"errors":    "[map[code:502 detail:code=502, message=test-error id:request-id status:502 title:Bad Gateway]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test http error handling with invalid error code", func(t *testing.T) {
		fn := func(fxjsonapi.Processor, echo.Context) error {
			return echo.NewHTTPError(0, "test-error")
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"Internal Server Error","detail":"code=0, message=test-error","status":"500","code":"500"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      500,
			"error":     "code=0, message=test-error",
			"errors":    "[map[code:500 detail:code=0, message=test-error id:request-id status:500 title:Internal Server Error]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test http error handling with obfuscation", func(t *testing.T) {
		t.Setenv("MODULES_HTTP_SERVER_ERRORS_OBFUSCATE", "true")

		fn := func(fxjsonapi.Processor, echo.Context) error {
			return echo.NewHTTPError(http.StatusBadGateway, "test-error")
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadGateway, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"Bad Gateway","detail":"Bad Gateway","status":"502","code":"502"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      502,
			"error":     "code=502, message=test-error",
			"errors":    "[map[code:502 detail:Bad Gateway id:request-id status:502 title:Bad Gateway]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test generic error handling", func(t *testing.T) {
		fn := func(fxjsonapi.Processor, echo.Context) error {
			return errors.New("test-error")
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"Internal Server Error","detail":"test-error","status":"500","code":"500"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      500,
			"error":     "test-error",
			"errors":    "[map[code:500 detail:test-error id:request-id status:500 title:Internal Server Error]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})

	t.Run("test generic error handling with obfuscation", func(t *testing.T) {
		t.Setenv("MODULES_HTTP_SERVER_ERRORS_OBFUSCATE", "true")

		fn := func(fxjsonapi.Processor, echo.Context) error {
			return errors.New("test-error")
		}

		httpServer, logBuffer := runTest(t, fn)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set(echo.HeaderContentType, jsonapi.MediaType)
		req.Header.Set(echo.HeaderXRequestID, "request-id")

		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())

		expected := `{"errors":[{"id":"request-id","title":"Internal Server Error","detail":"Internal Server Error","status":"500","code":"500"}]}`
		assert.Equal(t, fmt.Sprintf("%s\n", expected), rec.Body.String())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "error",
			"code":      500,
			"error":     "test-error",
			"errors":    "[map[code:500 detail:Internal Server Error id:request-id status:500 title:Internal Server Error]]",
			"requestID": "request-id",
			"message":   "json api error handler",
		})
	})
}
