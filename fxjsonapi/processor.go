package fxjsonapi

import (
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ Processor = (*DefaultProcessor)(nil)

type Processor interface {
	ProcessRequest(c echo.Context, data any, options ...ProcessorOption) error
	ProcessResponse(c echo.Context, code int, data any, options ...ProcessorOption) error
}

type DefaultProcessor struct {
	config *config.Config
}

func NewDefaultProcessor(config *config.Config) *DefaultProcessor {
	return &DefaultProcessor{
		config: config,
	}
}

func (p *DefaultProcessor) ProcessRequest(c echo.Context, data any, options ...ProcessorOption) error {
	ctx := c.Request().Context()

	processorOptions := DefaultProcessorOptions(p.config)
	for _, processorOption := range options {
		processorOption(&processorOptions)
	}

	var span oteltrace.Span

	if processorOptions.Trace {
		ctx, span = trace.CtxTracer(ctx).Start(ctx, "JSON API request processing")
	}

	err := jsonapi.UnmarshalPayload(c.Request().Body, data)
	if err != nil {
		errMsg := "JSON API request processing error"

		if processorOptions.Log {
			log.CtxLogger(ctx).Error().Err(err).Msg(errMsg)
		}

		if processorOptions.Trace && span != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, errMsg)
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if processorOptions.Log {
		log.CtxLogger(ctx).Debug().Msg("JSON API request processing success")
	}

	if processorOptions.Trace && span != nil {
		span.End()
	}

	return nil
}

func (p *DefaultProcessor) ProcessResponse(c echo.Context, code int, data any, options ...ProcessorOption) error {
	ctx := c.Request().Context()

	processorOptions := DefaultProcessorOptions(p.config)
	for _, processorOption := range options {
		processorOption(&processorOptions)
	}

	var span oteltrace.Span

	if processorOptions.Trace {
		ctx, span = trace.CtxTracer(ctx).Start(ctx, "JSON API response processing")
	}

	marshalledData, err := Marshall(data, MarshallParams{
		WithoutIncluded: !processorOptions.Included,
		Metadata:        processorOptions.Metadata,
	})
	if err != nil {
		errMsg := "JSON API response processing error"

		if processorOptions.Log {
			log.CtxLogger(ctx).Error().Err(err).Msg(errMsg)
		}

		if processorOptions.Trace && span != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, errMsg)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if processorOptions.Log {
		log.CtxLogger(ctx).Debug().Msg("JSON API response processing success")
	}

	if processorOptions.Trace && span != nil {
		span.End()
	}

	return c.Blob(code, jsonapi.MediaType, marshalledData)
}
