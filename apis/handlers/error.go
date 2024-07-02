package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"vibrain/apis"
)

func newHttpError(code int, message string) error {
	return fmt.Errorf("%d:%s", code, message)
}

func jsonError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(apis.ErrorResponse{
		Errors: []apis.Error{
			{
				Code:    code,
				Message: message,
			},
		},
		Success: false,
	})
	if err != nil {
		slog.Error("Error encoding error response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleRequestError(w http.ResponseWriter, r *http.Request, err error) {
	jsonError(w, http.StatusBadRequest, err.Error())
}

func handleResponseError(w http.ResponseWriter, r *http.Request, err error) {
	message := err.Error()
	msgs := strings.Split(message, ":")
	if len(msgs) > 1 {
		message = strings.Join(msgs[1:], ":")
	}
	code := strings.Trim(msgs[0], " ")
	codeNum, err := strconv.Atoi(code)
	if err != nil {
		codeNum = http.StatusInternalServerError
	}

	jsonError(w, codeNum, message)
}
