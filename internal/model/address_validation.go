package model

import (
	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

const (
	minAddressLength = 5
	maxAddressLength = 500
)

func validateAddressValue(value string) error {
	if err := validateStringNoLeadingOrTrailingWhitespace(value); err != nil {
		return err
	}
	if err := validateStringNoConsecutiveSpaces(value); err != nil {
		return err
	}
	if err := validateStringMinLength(value, minAddressLength); err != nil {
		return err
	}
	if err := validateStringMaxLength(value, maxAddressLength); err != nil {
		return err
	}

	return nil
}

func validateGeoPoint(point GeoPoint) error {
	if point.Latitude < -90 || point.Latitude > 90 {
		return errs.NewInvalidGeoPointError(
			errs.GeoPointFieldLatitude,
			errs.NewNumberOutOfRangeError(-90, 90, point.Latitude),
		)
	}

	if point.Longitude < -180 || point.Longitude > 180 {
		return errs.NewInvalidGeoPointError(
			errs.GeoPointFieldLongitude,
			errs.NewNumberOutOfRangeError(-180, 180, point.Longitude),
		)
	}

	return nil
}

func validateAddress(address Address) error {
	if err := validateAddressValue(address.Value); err != nil {
		return errs.NewInvalidAddressError(errs.AddressFieldValue, err)
	}

	if address.Point != nil {
		if err := validateGeoPoint(*address.Point); err != nil {
			return errs.NewInvalidAddressError(errs.AddressFieldPoint, err)
		}
	}

	return nil
}
