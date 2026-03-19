package model

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/ulbwa/medincident-command-service/internal/common/errors"
)

func TestValidateStringMinLength(t *testing.T) {
	t.Parallel()

	t.Run("TooShort", func(t *testing.T) {
		t.Parallel()

		err := validateStringMinLength("ab", 5)
		require.Error(t, err)

		var tooShortErr *errs.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
		assert.Equal(t, 5, tooShortErr.MinLength)
		assert.Equal(t, 2, tooShortErr.ActualLength)
		assert.Equal(t, "ab", tooShortErr.ActualValue)
	})

	t.Run("ExactMinLength", func(t *testing.T) {
		t.Parallel()

		err := validateStringMinLength("abcde", 5)
		require.NoError(t, err)
	})

	t.Run("AboveMinLength", func(t *testing.T) {
		t.Parallel()

		err := validateStringMinLength("abcdef", 5)
		require.NoError(t, err)
	})

	t.Run("UnicodeCharacters", func(t *testing.T) {
		t.Parallel()

		// "Ан" — 2 runes, not 4 bytes
		err := validateStringMinLength("Ан", 5)
		require.Error(t, err)

		var tooShortErr *errs.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
		assert.Equal(t, 5, tooShortErr.MinLength)
		assert.Equal(t, 2, tooShortErr.ActualLength)
	})

	t.Run("EmptyString", func(t *testing.T) {
		t.Parallel()

		err := validateStringMinLength("", 1)
		require.Error(t, err)

		var tooShortErr *errs.StringTooShortError
		require.True(t, stderrors.As(err, &tooShortErr))
		assert.Equal(t, 1, tooShortErr.MinLength)
		assert.Equal(t, 0, tooShortErr.ActualLength)
	})
}

func TestValidateStringMaxLength(t *testing.T) {
	t.Parallel()

	t.Run("TooLong", func(t *testing.T) {
		t.Parallel()

		err := validateStringMaxLength("abcdef", 5)
		require.Error(t, err)

		var tooLongErr *errs.StringTooLongError
		require.True(t, stderrors.As(err, &tooLongErr))
		assert.Equal(t, 5, tooLongErr.MaxLength)
		assert.Equal(t, 6, tooLongErr.ActualLength)
		assert.Equal(t, "abcdef", tooLongErr.ActualValue)
	})

	t.Run("ExactMaxLength", func(t *testing.T) {
		t.Parallel()

		err := validateStringMaxLength("abcde", 5)
		require.NoError(t, err)
	})

	t.Run("BelowMaxLength", func(t *testing.T) {
		t.Parallel()

		err := validateStringMaxLength("ab", 5)
		require.NoError(t, err)
	})

	t.Run("UnicodeCharacters", func(t *testing.T) {
		t.Parallel()

		// "Привет" — 6 runes, not 12 bytes
		err := validateStringMaxLength("Привет", 5)
		require.Error(t, err)

		var tooLongErr *errs.StringTooLongError
		require.True(t, stderrors.As(err, &tooLongErr))
		assert.Equal(t, 5, tooLongErr.MaxLength)
		assert.Equal(t, 6, tooLongErr.ActualLength)
	})
}

func TestValidateStringNoLeadingOrTrailingWhitespace(t *testing.T) {
	t.Parallel()

	t.Run("LeadingSpace", func(t *testing.T) {
		t.Parallel()

		err := validateStringNoLeadingOrTrailingWhitespace(" abc")
		require.Error(t, err)

		var wsErr *errs.StringLeadingOrTrailingWhitespaceError
		require.True(t, stderrors.As(err, &wsErr))
		assert.Equal(t, " abc", wsErr.ActualValue)
	})

	t.Run("TrailingSpace", func(t *testing.T) {
		t.Parallel()

		err := validateStringNoLeadingOrTrailingWhitespace("abc ")
		require.Error(t, err)

		var wsErr *errs.StringLeadingOrTrailingWhitespaceError
		require.True(t, stderrors.As(err, &wsErr))
		assert.Equal(t, "abc ", wsErr.ActualValue)
	})

	t.Run("LeadingAndTrailingSpace", func(t *testing.T) {
		t.Parallel()

		err := validateStringNoLeadingOrTrailingWhitespace(" abc ")
		require.Error(t, err)

		var wsErr *errs.StringLeadingOrTrailingWhitespaceError
		require.True(t, stderrors.As(err, &wsErr))
		assert.Equal(t, " abc ", wsErr.ActualValue)
	})

	t.Run("InternalSpaceAllowed", func(t *testing.T) {
		t.Parallel()

		err := validateStringNoLeadingOrTrailingWhitespace("hello world")
		require.NoError(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		err := validateStringNoLeadingOrTrailingWhitespace("hello")
		require.NoError(t, err)
	})

	t.Run("EmptyString", func(t *testing.T) {
		t.Parallel()

		err := validateStringNoLeadingOrTrailingWhitespace("")
		require.NoError(t, err)
	})
}
