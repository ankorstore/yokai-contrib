package fxjsonapi

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
	config *config.Config
}

func NewErrorHandler(config *config.Config) *ErrorHandler {
	return &ErrorHandler{
		config: config,
	}
}

func (h *ErrorHandler) Handle() echo.HTTPErrorHandler {
	obfuscate := !h.config.AppDebug() || h.config.GetBool("modules.http.server.errors.obfuscate")

	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		logger := log.CtxLogger(c.Request().Context())

		var outErrors []*jsonapi.ErrorObject
		var outCode int

		var jsonErr *jsonapi.ErrorObject
		var httpErr *echo.HTTPError
		var valErr validator.ValidationErrors

		switch {
		case errors.As(err, &jsonErr):
			outErrors, outCode = h.handleJSONAPIError(c, jsonErr, obfuscate)
		case errors.As(err, &httpErr):
			outErrors, outCode = h.handleHTTPError(c, httpErr, obfuscate)
		case errors.As(err, &valErr):
			outErrors, outCode = h.handleValidationError(c, valErr, obfuscate)
		default:
			outErrors, outCode = h.handleGenericError(c, err, obfuscate)
		}

		logger.
			Error().
			Err(err).
			Any("errors", outErrors).
			Int("code", outCode).
			Msg("json api error handler")

		c.Set("Content-Type", jsonapi.MediaType)

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(outCode)

			return
		}

		buf := bytes.Buffer{}
		err = jsonapi.MarshalErrors(&buf, outErrors)
		if err != nil {
			logger.Error().Err(err).Msg("json api error handler marshall failure")
		}

		err = c.Blob(outCode, jsonapi.MediaType, buf.Bytes())
		if err != nil {
			logger.Error().Err(err).Msg("json api error handler blob failure")
		}
	}
}

func (h *ErrorHandler) handleJSONAPIError(c echo.Context, inErr *jsonapi.ErrorObject, obfuscate bool) ([]*jsonapi.ErrorObject, int) {
	outErr := &jsonapi.ErrorObject{
		ID:     httpserver.CtxRequestId(c),
		Title:  inErr.Title,
		Detail: inErr.Detail,
		Status: inErr.Status,
		Code:   inErr.Code,
		Meta:   inErr.Meta,
	}

	outCode, err := strconv.Atoi(inErr.Status)
	if err != nil {
		outCode = http.StatusInternalServerError
	}

	if obfuscate {
		outErr.Detail = http.StatusText(outCode)
	}

	return []*jsonapi.ErrorObject{outErr}, outCode
}

func (h *ErrorHandler) handleHTTPError(c echo.Context, inErr *echo.HTTPError, obfuscate bool) ([]*jsonapi.ErrorObject, int) {
	outErr := &jsonapi.ErrorObject{
		ID:     httpserver.CtxRequestId(c),
		Title:  http.StatusText(inErr.Code),
		Detail: inErr.Error(),
		Status: fmt.Sprintf("%d", inErr.Code),
		Code:   fmt.Sprintf("%d", inErr.Code),
	}

	if obfuscate {
		outErr.Detail = http.StatusText(inErr.Code)
	}

	return []*jsonapi.ErrorObject{outErr}, inErr.Code
}

func (h *ErrorHandler) handleValidationError(c echo.Context, inErr validator.ValidationErrors, obfuscate bool) ([]*jsonapi.ErrorObject, int) {
	var outErrs []*jsonapi.ErrorObject

	for k, iErr := range inErr {
		outErr := &jsonapi.ErrorObject{
			ID:     fmt.Sprintf("%s#%d", httpserver.CtxRequestId(c), k),
			Title:  iErr.Field(),
			Detail: iErr.Error(),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Code:   fmt.Sprintf("%d", http.StatusBadRequest),
		}

		if obfuscate {
			outErr.Detail = fmt.Sprintf("Validation error for %s", iErr.Field())
		}

		outErrs = append(outErrs, outErr)
	}

	return outErrs, http.StatusBadRequest
}

func (h *ErrorHandler) handleGenericError(c echo.Context, inErr error, obfuscate bool) ([]*jsonapi.ErrorObject, int) {
	outErr := &jsonapi.ErrorObject{
		ID:     httpserver.CtxRequestId(c),
		Title:  http.StatusText(http.StatusInternalServerError),
		Detail: inErr.Error(),
		Status: fmt.Sprintf("%d", http.StatusInternalServerError),
		Code:   fmt.Sprintf("%d", http.StatusInternalServerError),
	}

	if obfuscate {
		outErr.Detail = http.StatusText(http.StatusInternalServerError)
	}

	return []*jsonapi.ErrorObject{outErr}, http.StatusInternalServerError
}
