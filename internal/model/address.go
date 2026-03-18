package model

import (
	"fmt"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

// GeoPoint represents geographic coordinates (latitude and longitude).
type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

// NewGeoPoint creates a validated GeoPoint.
func NewGeoPoint(latitude, longitude float64) (GeoPoint, error) {
	if latitude < -90 || latitude > 90 {
		return GeoPoint{}, fmt.Errorf("%w: must be between -90 and 90", errs.ErrInvalidLatitude)
	}
	if longitude < -180 || longitude > 180 {
		return GeoPoint{}, fmt.Errorf("%w: must be between -180 and 180", errs.ErrInvalidLongitude)
	}
	return GeoPoint{Latitude: latitude, Longitude: longitude}, nil
}

// Equals compares two GeoPoint instances.
func (g GeoPoint) Equals(other GeoPoint) bool {
	return g.Latitude == other.Latitude && g.Longitude == other.Longitude
}

// Address represents a physical or legal address with optional geographic coordinates.
type Address struct {
	Value string
	Point *GeoPoint
}

// NewAddress creates a validated Address.
func NewAddress(value string, point *GeoPoint) (Address, error) {
	if err := validateAddressValue(value); err != nil {
		return Address{}, err
	}
	return Address{Value: value, Point: point}, nil
}

// Equals compares two Address instances.
func (a Address) Equals(other Address) bool {
	if a.Value != other.Value {
		return false
	}

	if a.Point == nil && other.Point == nil {
		return true
	}
	if a.Point == nil || other.Point == nil {
		return false
	}

	return a.Point.Equals(*other.Point)
}
