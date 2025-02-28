package handler

import (
	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/labstack/echo/v4"
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
