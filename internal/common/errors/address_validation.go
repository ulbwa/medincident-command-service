package errors

import "fmt"

type AddressField string

const (
	AddressFieldValue AddressField = "value"
	AddressFieldPoint AddressField = "point"
)

type InvalidAddressError struct {
	Field  AddressField
	Reason error
}

func (e *InvalidAddressError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid address field %s: %s", e.Field, e.Reason)
}

func (e *InvalidAddressError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidAddressError(field AddressField, reason error) *InvalidAddressError {
	return &InvalidAddressError{Field: field, Reason: reason}
}

type GeoPointField string

const (
	GeoPointFieldLatitude  GeoPointField = "latitude"
	GeoPointFieldLongitude GeoPointField = "longitude"
)

type InvalidGeoPointError struct {
	Field  GeoPointField
	Reason error
}

func (e *InvalidGeoPointError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("invalid geo point field %s: %s", e.Field, e.Reason)
}

func (e *InvalidGeoPointError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Reason
}

func NewInvalidGeoPointError(field GeoPointField, reason error) *InvalidGeoPointError {
	return &InvalidGeoPointError{Field: field, Reason: reason}
}

type NumberOutOfRangeError struct {
	Min    float64
	Max    float64
	Actual float64
}

func (e *NumberOutOfRangeError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("number is out of range: expected [%v, %v], got %v", e.Min, e.Max, e.Actual)
}

func NewNumberOutOfRangeError(min, max, actual float64) *NumberOutOfRangeError {
	return &NumberOutOfRangeError{Min: min, Max: max, Actual: actual}
}
