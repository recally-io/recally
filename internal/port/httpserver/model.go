package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Description	JSONResult represents the structure of the JSON response.
type JSONResult struct {
	// Success is a boolean value that indicates whether the request was successful.
	Success bool `json:"success"`
	// Code is an integer value that represents the HTTP status code.
	Code int `json:"code" example:"200"`
	// Message is a string value that represents the message of the response.
	Message string `json:"message" example:"OK"`
	// Data is an interface value that represents the data of the response.
	Data any `json:"data"`
	// Error is an error value that represents the error of the response.
	Error error `json:"error"`
}

func JsonResponse(c echo.Context, code int, data any) error {
	return c.JSON(code, JSONResult{
		Success: true,
		Code:    code,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, err error) error {
	return echo.NewHTTPError(code).SetInternal(err)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	msg := err.Error()

	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code

		if he.Internal != nil {
			msg = he.Internal.Error()
		}
	}

	_ = c.JSON(code, JSONResult{
		Success: false,
		Code:    code,
		Message: msg,
	})
}

// bindAndValidate binds the request data to the provided struct and validates it.
// It returns an error if either binding or validation fails.
//
// Parameters:
//   - c: The echo.Context object representing the current HTTP request context.
//   - req: A pointer to the struct where the request data will be bound.
//
// Returns:
//   - An error if binding or validation fails, nil otherwise.
func bindAndValidate(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	if err := c.Validate(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	return nil
}
