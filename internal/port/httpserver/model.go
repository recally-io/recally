package httpserver

import "github.com/labstack/echo/v4"

// @Description JSONResult represents the structure of the JSON response.
type JSONResult struct {
	// Success is a boolean value that indicates whether the request was successful.
	Success bool `json:"success"`
	// Code is an integer value that represents the HTTP status code.
	Code int `json:"code" example:"200"`
	// Message is a string value that represents the message of the response.
	Message string `json:"message" example:"OK"`
	// Data is an interface value that represents the data of the response.
	Data interface{} `json:"data"`
	// Error is an error value that represents the error of the response.
	Error error `json:"error"`
}

func JsonResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, JSONResult{
		Success: true,
		Code:    code,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, err error) error {
	return c.JSON(code, JSONResult{
		Success: false,
		Code:    code,
		Message: err.Error(),
	})
}
