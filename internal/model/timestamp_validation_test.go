package model

import (
	stderrors "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func TestValidateTimestampNotBefore(t *testing.T) {
	t.Parallel()

	minimum := time.Now().UTC()
	before := minimum.Add(-time.Second)

	err := validateTimestampNotBefore(before, minimum)
	require.Error(t, err)

	var beforeErr *errs.TimestampBeforeMinimumError
	require.True(t, stderrors.As(err, &beforeErr))
	require.Equal(t, before, beforeErr.Actual)
	require.Equal(t, minimum, beforeErr.ExpectedAfter)

	err = validateTimestampNotBefore(minimum, minimum)
	require.NoError(t, err)
}
