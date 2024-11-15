package db

import "strings"

const ErrNotFound = "no rows in result set"

func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), ErrNotFound)
}
