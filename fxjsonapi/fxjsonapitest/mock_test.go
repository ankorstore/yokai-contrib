package fxjsonapitest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/fxjsonapitest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type wrapper struct {
	processor *fxjsonapitest.ProcessorMock
}

func (w *wrapper) processRequest(c echo.Context, data any, options ...fxjsonapi.ProcessorOption) error {
	return w.processor.ProcessRequest(c, data, options...)
}

func (w *wrapper) processResponse(c echo.Context, code int, data any, options ...fxjsonapi.ProcessorOption) error {
	return w.processor.ProcessResponse(c, code, data, options...)
}

func TestProcessorMock(t *testing.T) {
	t.Parallel()

	t.Run("request processing", func(t *testing.T) {
		t.Parallel()

		c := echo.New().NewContext(httptest.NewRequest("GET", "/test", nil), nil)

		s := struct{}{}

		o := []fxjsonapi.ProcessorOption{
			fxjsonapi.WithLog(true),
		}

		m := new(fxjsonapitest.ProcessorMock)
		m.On("ProcessRequest", c, s, o).Return(nil).Once()

		w := &wrapper{m}

		err := w.processRequest(c, s, o...)
		require.NoError(t, err)

		m.AssertExpectations(t)
	})

	t.Run("response processing", func(t *testing.T) {
		t.Parallel()

		c := echo.New().NewContext(httptest.NewRequest("GET", "/test", nil), nil)

		s := struct{}{}

		o := []fxjsonapi.ProcessorOption{
			fxjsonapi.WithLog(true),
		}

		h := http.StatusInternalServerError

		m := new(fxjsonapitest.ProcessorMock)
		m.On("ProcessResponse", c, h, s, o).Return(errors.New("error")).Once()

		w := &wrapper{m}

		err := w.processResponse(c, h, s, o...)
		require.Error(t, err)
		require.Equal(t, "error", err.Error())

		m.AssertExpectations(t)
	})
}
