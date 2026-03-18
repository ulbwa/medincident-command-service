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

// Equals compares two GeoPoint instances.
func (g GeoPoint) Equals(other GeoPoint) bool {
	return g.Latitude == other.Latitude && g.Longitude == other.Longitude
}

// Address represents a physical or legal address with optional geographic coordinates.
type Address struct {
	Value string
	Point *GeoPoint
}

func (a Address) copy() Address {
	cloned := a
	if a.Point != nil {
		cloned.Point = a.Point.copy()
	}
	return cloned
}

// NewAddress creates a validated Address.
func NewAddress(value string, point *GeoPoint) (Address, error) {
	address := Address{Value: value}
	if point != nil {
		address.Point = point.copy()
	}
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
