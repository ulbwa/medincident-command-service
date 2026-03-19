package model

import (
	"time"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func validateTimestampNotBefore(actual, minimum time.Time) error {
	if actual.Before(minimum) {
		return errs.NewTimestampBeforeMinimumError(actual, minimum)
	}
	return nil
}

func validateTimestampNotAfter(actual, maximum time.Time) error {
	if actual.After(maximum) {
		return errs.NewTimestampAfterMaximumError(actual, maximum)
	}
	return nil
}
