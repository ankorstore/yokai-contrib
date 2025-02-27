package handler

import (
	"net/http"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/model"
	"github.com/labstack/echo/v4"
)

type FooBarHandler struct {
	processor fxjsonapi.Processor
}

func NewFooBarHandler(processor fxjsonapi.Processor) *FooBarHandler {
	return &FooBarHandler{
		processor: processor,
	}
}

func (h *FooBarHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := model.Foo{}

		err := h.processor.ProcessRequest(c, &foo)
		if err != nil {
			return err
		}

		return h.processor.ProcessResponse(c, http.StatusOK, &foo)
	}
}
