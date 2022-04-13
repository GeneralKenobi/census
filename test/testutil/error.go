package testutil

import "github.com/GeneralKenobi/census/internal/api"

// IsBadInputError checks if err is api.StatusError with status api.StatusBadInput.
func IsBadInputError(err error) bool {
	return isStatusError(err, api.StatusBadInput)
}

func isStatusError(err error, expectedStatus api.Status) bool {
	if err == nil {
		return false
	}

	statusErr, ok := err.(api.StatusError)
	if !ok {
		return false
	}

	return statusErr.Status() == expectedStatus
}
