package fxjsonapi

import (
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

func MarshallResponse(c echo.Context, code int, params MarshallParams) error {
	mp, err := Marshall(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(code, jsonapi.MediaType, mp)
}
