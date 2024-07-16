package handlers

import "github.com/labstack/echo/v4"

func JsonResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c echo.Context, code int, err error) error {
	return c.JSON(code, map[string]interface{}{
		"success": false,
		"errors":  err.Error(),
		"data":    nil,
	})
}
