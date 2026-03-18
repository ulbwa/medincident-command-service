package model

import (
	"fmt"
	"strings"
	"unicode/utf8"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minAddressLength = 5
	maxAddressLength = 500
)

func validateAddressValue(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed != value {
		return fmt.Errorf("%w: must not have leading or trailing whitespace", errs.ErrInvalidAddressValue)
	}

	length := utf8.RuneCountInString(value)
	if length < minAddressLength {
		return fmt.Errorf("%w: too short (min %d)", errs.ErrInvalidAddressValue, minAddressLength)
	}
	if length > maxAddressLength {
		return fmt.Errorf("%w: too long (max %d)", errs.ErrInvalidAddressValue, maxAddressLength)
	}

	return nil
}

func validateGeoPoint(point GeoPoint) error {
	if point.Latitude < -90 || point.Latitude > 90 {
		return fmt.Errorf("%w: must be between -90 and 90", errs.ErrInvalidLatitude)
	}

	if point.Longitude < -180 || point.Longitude > 180 {
		return fmt.Errorf("%w: must be between -180 and 180", errs.ErrInvalidLongitude)
	}

	return nil
}

func validateAddress(address Address) error {
	if err := validateAddressValue(address.Value); err != nil {
		return err
	}

	if address.Point != nil {
		if err := validateGeoPoint(*address.Point); err != nil {
			return err
		}
	}

	return nil
}
