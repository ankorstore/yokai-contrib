package fxjsonapi

import (
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

func UnmarshallRequest(c echo.Context, data any) error {
	err := jsonapi.UnmarshalPayload(c.Request().Body, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
