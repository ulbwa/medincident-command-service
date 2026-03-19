package tests

import (
	stderrors "errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ulbwa/medincident-command-service/internal/common/errors"
	"github.com/ulbwa/medincident-command-service/internal/model"
)

func TestGeoPoint_NewGeoPoint(t *testing.T) {
	t.Parallel()

	t.Run("ValidCoordinates", func(t *testing.T) {
		t.Parallel()
		gp, err := model.NewGeoPoint(55.7558, 37.6173) // Moscow
		require.NoError(t, err)
		assert.Equal(t, 55.7558, gp.Latitude)
		assert.Equal(t, 37.6173, gp.Longitude)
	})

	t.Run("EdgeLatitudeValues", func(t *testing.T) {
		t.Parallel()
		// Min lat
		gp, err := model.NewGeoPoint(-90, 0)
		require.NoError(t, err)
		assert.Equal(t, -90.0, gp.Latitude)

		// Max lat
		gp, err = model.NewGeoPoint(90, 0)
		require.NoError(t, err)
		assert.Equal(t, 90.0, gp.Latitude)
	})

	t.Run("EdgeLongitudeValues", func(t *testing.T) {
		t.Parallel()
		// Min long
		gp, err := model.NewGeoPoint(0, -180)
		require.NoError(t, err)
		assert.Equal(t, -180.0, gp.Longitude)

		// Max long
		gp, err = model.NewGeoPoint(0, 180)
		require.NoError(t, err)
		assert.Equal(t, 180.0, gp.Longitude)
	})

	t.Run("InvalidLatitude", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewGeoPoint(-91, 0)
		var geoErr *errors.InvalidGeoPointError
		require.True(t, stderrors.As(err, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLatitude, geoErr.Field)

		_, err = model.NewGeoPoint(91, 0)
		require.True(t, stderrors.As(err, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLatitude, geoErr.Field)
	})

	t.Run("InvalidLongitude", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewGeoPoint(0, -181)
		var geoErr *errors.InvalidGeoPointError
		require.True(t, stderrors.As(err, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLongitude, geoErr.Field)

		_, err = model.NewGeoPoint(0, 181)
		require.True(t, stderrors.As(err, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLongitude, geoErr.Field)
	})
}

func TestGeoPoint_NewGeoPoint_AdditionalCases(t *testing.T) {
	t.Parallel()

	t.Run("ValidCoordinates", func(t *testing.T) {
		t.Parallel()
		gp, err := model.NewGeoPoint(55.7558, 37.6173)
		require.NoError(t, err)
		assert.Equal(t, 55.7558, gp.Latitude)
		assert.Equal(t, 37.6173, gp.Longitude)
	})

	t.Run("InvalidLatitude", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewGeoPoint(-91, 0)
		var geoErr *errors.InvalidGeoPointError
		require.True(t, stderrors.As(err, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLatitude, geoErr.Field)
	})

	t.Run("InvalidLongitude", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewGeoPoint(0, 181)
		var geoErr *errors.InvalidGeoPointError
		require.True(t, stderrors.As(err, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLongitude, geoErr.Field)
	})
}

func TestGeoPoint_Equals(t *testing.T) {
	t.Parallel()

	gp1, _ := model.NewGeoPoint(55.7558, 37.6173)
	gp2, _ := model.NewGeoPoint(55.7558, 37.6173)
	gp3, _ := model.NewGeoPoint(59.9343, 30.3351) // SPB

	assert.True(t, gp1.Equals(gp2))
	assert.False(t, gp1.Equals(gp3))
}

func TestAddress_NewAddress(t *testing.T) {
	t.Parallel()

	t.Run("ValidAddressWithoutPoint", func(t *testing.T) {
		t.Parallel()
		addr, err := model.NewAddress("Moscow, Red Square, 1", nil)
		require.NoError(t, err)
		assert.Equal(t, "Moscow, Red Square, 1", addr.Value)
		assert.Nil(t, addr.Point)
	})

	t.Run("ValidAddressWithPoint", func(t *testing.T) {
		t.Parallel()
		gp, _ := model.NewGeoPoint(55.7539, 37.6208)
		addr, err := model.NewAddress("Moscow, Red Square, 1", &gp)
		require.NoError(t, err)
		assert.Equal(t, "Moscow, Red Square, 1", addr.Value)
		assert.NotNil(t, addr.Point)
		assert.Equal(t, 55.7539, addr.Point.Latitude)
	})

	t.Run("EmptyAddress", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewAddress("", nil)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})

	t.Run("TooShortAddress", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewAddress("ABC", nil) // 3 chars < 5
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})

	t.Run("TooLongAddress", func(t *testing.T) {
		t.Parallel()
		longAddr := strings.Repeat("A", 501)
		_, err := model.NewAddress(longAddr, nil)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})

	t.Run("LeadingWhitespace", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewAddress(" Moscow, Red Square, 1", nil)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})

	t.Run("TrailingWhitespace", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewAddress("Moscow, Red Square, 1 ", nil)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})

	t.Run("InvalidPoint", func(t *testing.T) {
		t.Parallel()
		point := model.GeoPoint{Latitude: 100, Longitude: 37.6208}
		_, err := model.NewAddress("Moscow, Red Square, 1", &point)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldPoint, addressErr.Field)
		var geoErr *errors.InvalidGeoPointError
		require.True(t, stderrors.As(addressErr.Reason, &geoErr))
		assert.Equal(t, errors.GeoPointFieldLatitude, geoErr.Field)
	})
}

func TestAddress_NewAddress_AdditionalCases(t *testing.T) {
	t.Parallel()

	t.Run("ValidAddressWithoutPoint", func(t *testing.T) {
		t.Parallel()
		addr, err := model.NewAddress("Moscow, Red Square, 1", nil)
		require.NoError(t, err)
		assert.Equal(t, "Moscow, Red Square, 1", addr.Value)
		assert.Nil(t, addr.Point)
	})

	t.Run("InvalidAddressValue", func(t *testing.T) {
		t.Parallel()
		_, err := model.NewAddress("", nil)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldValue, addressErr.Field)
	})

	t.Run("InvalidPoint", func(t *testing.T) {
		t.Parallel()
		point := model.GeoPoint{Latitude: 100, Longitude: 37.6208}
		_, err := model.NewAddress("Moscow, Red Square, 1", &point)
		var addressErr *errors.InvalidAddressError
		require.True(t, stderrors.As(err, &addressErr))
		assert.Equal(t, errors.AddressFieldPoint, addressErr.Field)
	})

	t.Run("CopyPoint", func(t *testing.T) {
		t.Parallel()
		point := model.GeoPoint{Latitude: 55.7539, Longitude: 37.6208}
		addr, err := model.NewAddress("Moscow, Red Square, 1", &point)
		require.NoError(t, err)
		require.NotNil(t, addr.Point)

		point.Latitude = 10
		assert.Equal(t, 55.7539, addr.Point.Latitude)
		assert.NotSame(t, &point, addr.Point)
	})
}

func TestAddress_NewAddress_CopyPoint(t *testing.T) {
	t.Parallel()

	point := model.GeoPoint{Latitude: 55.7539, Longitude: 37.6208}
	addr, err := model.NewAddress("Moscow, Red Square, 1", &point)
	require.NoError(t, err)
	require.NotNil(t, addr.Point)

	point.Latitude = 10
	assert.Equal(t, 55.7539, addr.Point.Latitude)
	assert.NotSame(t, &point, addr.Point)
}

func TestAddress_Equals(t *testing.T) {
	t.Parallel()

	gp, _ := model.NewGeoPoint(55.7539, 37.6208)
	addr1, _ := model.NewAddress("Moscow, Red Square, 1", &gp)
	addr2, _ := model.NewAddress("Moscow, Red Square, 1", &gp)
	addr3, _ := model.NewAddress("Moscow, Red Square, 1", nil)
	addr4, _ := model.NewAddress("Saint Petersburg, Palace Square, 1", &gp)

	assert.True(t, addr1.Equals(addr2))
	assert.False(t, addr1.Equals(addr3)) // different point (nil vs non-nil)
	assert.False(t, addr1.Equals(addr4)) // different value

	addr5, _ := model.NewAddress("Moscow, Red Square, 1", nil)
	assert.True(t, addr3.Equals(addr5)) // both nil points
}
