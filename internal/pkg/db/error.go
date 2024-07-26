package db

const ErrNotFound = "no rows in result set"

func IsNotFound(err error) bool {
	return err.Error() == ErrNotFound
}
