package fxjsonapitest

import (
	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

var _ fxjsonapi.Processor = (*ProcessorMock)(nil)

// ProcessorMock is a [Processor] mock.
type ProcessorMock struct {
	mock.Mock
}

// ProcessRequest is a mocked ProcessRequest implementation.
func (m *ProcessorMock) ProcessRequest(c echo.Context, data any, options ...fxjsonapi.ProcessorOption) error {
	args := m.Called(c, data, options)

	return args.Error(0)
}

// ProcessResponse is a mocked ProcessResponse implementation.
func (m *ProcessorMock) ProcessResponse(c echo.Context, code int, data any, options ...fxjsonapi.ProcessorOption) error {
	args := m.Called(c, code, data, options)

	return args.Error(0)
}
