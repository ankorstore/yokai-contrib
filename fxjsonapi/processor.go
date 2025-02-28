package fxjsonapi

import (
	"net/http"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ Processor = (*DefaultProcessor)(nil)

// Processor is the interface for json api processors implementations.
type Processor interface {
	ProcessRequest(c echo.Context, data any, options ...ProcessorOption) error
	ProcessResponse(c echo.Context, code int, data any, options ...ProcessorOption) error
}

// DefaultProcessor is the default [Processor] implementation.
type DefaultProcessor struct {
	config *config.Config
}

// NewDefaultProcessor returns a new [DefaultProcessor] instance.
func NewDefaultProcessor(config *config.Config) *DefaultProcessor {
	return &DefaultProcessor{
		config: config,
	}
}

// ProcessRequest processes a json api request.
//
//nolint:cyclop
func (p *DefaultProcessor) ProcessRequest(c echo.Context, data any, options ...ProcessorOption) error {
	processorOptions := DefaultProcessorOptions(p.config)
	for _, processorOption := range options {
		processorOption(&processorOptions)
	}

	req := c.Request()
	ctx := req.Context()

	var span oteltrace.Span

	if processorOptions.Trace {
		ctx, span = trace.CtxTracer(ctx).Start(ctx, "JSON API request processing")
	}

	defer func() {
		if processorOptions.Trace && span != nil {
			span.End()
		}
	}()

	logger := log.CtxLogger(ctx)

	contentTypeHeader, _, _ := strings.Cut(req.Header.Get(echo.HeaderContentType), ";")
	contentType := strings.ToLower(strings.TrimSpace(contentTypeHeader))

	if contentType != jsonapi.MediaType {
		errMsg := "JSON API request invalid content type"

		if processorOptions.Log {
			logger.Error().Msg(errMsg)
		}

		if processorOptions.Trace && span != nil {
			span.SetStatus(codes.Error, errMsg)
		}

		return echo.NewHTTPError(http.StatusUnsupportedMediaType, errMsg)
	}

	err := jsonapi.UnmarshalPayload(c.Request().Body, data)
	if err != nil {
		errMsg := "JSON API request processing error"

		if processorOptions.Log {
			logger.Error().Err(err).Msg(errMsg)
		}

		if processorOptions.Trace && span != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, errMsg)
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	okMsg := "JSON API request processing success"

	if processorOptions.Log {
		logger.Debug().Msg(okMsg)
	}

	if processorOptions.Trace && span != nil {
		span.SetStatus(codes.Ok, okMsg)
	}

	return nil
}

// ProcessResponse processes a json api response.
//
//nolint:cyclop
func (p *DefaultProcessor) ProcessResponse(c echo.Context, code int, data any, options ...ProcessorOption) error {
	processorOptions := DefaultProcessorOptions(p.config)
	for _, processorOption := range options {
		processorOption(&processorOptions)
	}

	ctx := c.Request().Context()

	var span oteltrace.Span

	if processorOptions.Trace {
		ctx, span = trace.CtxTracer(ctx).Start(ctx, "JSON API response processing")
	}

	defer func() {
		if processorOptions.Trace && span != nil {
			span.End()
		}
	}()

	logger := log.CtxLogger(ctx)

	marshalledData, err := Marshall(data, MarshallParams{
		WithoutIncluded: !processorOptions.Included,
		Metadata:        processorOptions.Metadata,
	})
	if err != nil {
		errMsg := "JSON API response processing error"

		if processorOptions.Log {
			logger.Error().Err(err).Msg(errMsg)
		}

		if processorOptions.Trace && span != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, errMsg)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	okMsg := "JSON API response processing success"

	if processorOptions.Log {
		logger.Debug().Msg("JSON API response processing success")
	}

	if processorOptions.Trace && span != nil {
		span.SetStatus(codes.Ok, okMsg)
	}

	return c.Blob(code, jsonapi.MediaType, marshalledData)
}
