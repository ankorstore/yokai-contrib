package handler

import (
	"net/http"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/model"
	"github.com/labstack/echo/v4"
)

type JSONAPIHandler struct {
	processor fxjsonapi.Processor
}

func NewJSONAPIHandler(processor fxjsonapi.Processor) *JSONAPIHandler {
	return &JSONAPIHandler{
		processor: processor,
	}
}

func (h *JSONAPIHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := model.Foo{
			ID:   123,
			Name: "foo",
			Bar: &Bar{
				ID:   456,
				Name: "bar",
			},
		}

		return h.processor.ProcessResponse(
			c,
			http.StatusOK,
			&foo,
			fxjsonapi.WithMetadata(map[string]interface{}{
				"some": "meta",
			}),
			fxjsonapi.WithIncluded(true),
			fxjsonapi.WithLog(true),
			fxjsonapi.WithTrace(true),
		)

	}
}
