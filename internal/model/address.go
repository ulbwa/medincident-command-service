package model

// GeoPoint represents geographic coordinates (latitude and longitude).
type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

func (g GeoPoint) copy() *GeoPoint {
	cloned := g
	return &cloned
}

// NewGeoPoint creates a validated GeoPoint.
func NewGeoPoint(latitude, longitude float64) (GeoPoint, error) {
	point := GeoPoint{Latitude: latitude, Longitude: longitude}
	if err := validateGeoPoint(point); err != nil {
		return GeoPoint{}, err
	}

	return point, nil
}

// RestoreGeoPoint restores an existing validated GeoPoint.
func RestoreGeoPoint(latitude, longitude float64) (GeoPoint, error) {
	return NewGeoPoint(latitude, longitude)
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
	var pointCopy *GeoPoint
	if point != nil {
		pointCopy = point.copy()
	}

	address := Address{Value: value, Point: pointCopy}
	if err := validateAddress(address); err != nil {
		return Address{}, err
	}

	return address, nil
}

// RestoreAddress restores an existing validated Address.
func RestoreAddress(value string, point *GeoPoint) (Address, error) {
	var pointCopy *GeoPoint
	if point != nil {
		pointCopy = point.copy()
	}

	address := Address{Value: value, Point: pointCopy}
	if err := validateAddress(address); err != nil {
		return Address{}, err
	}

	return address, nil
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
