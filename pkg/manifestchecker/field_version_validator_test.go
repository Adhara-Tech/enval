package manifestchecker_test

import (
	"testing"

	"github.com/Adhara-Tech/enval/pkg/manifestchecker"

	"github.com/stretchr/testify/require"
)

var semverValidator = manifestchecker.SemverFieldVersionValidator{}
var exactMatchValidator = manifestchecker.ExactMatchFieldVersionValidator{}

func TestSemverFieldVersionValidator_Validate_ValidVersion(t *testing.T) {
	result, err := semverValidator.Validate("12.0.1", ">12.0.0")
	require.Nil(t, err)
	require.True(t, result)
}

func TestSemverFieldVersionValidator_Validate_InvalidVersion(t *testing.T) {
	result, err := semverValidator.Validate("12.0.1", "<12.0.0")
	require.Nil(t, err)
	require.False(t, result)
}

func TestSemverFieldVersionValidator_Validate_Error_ExpectedVersionNotSemver(t *testing.T) {
	result, err := semverValidator.Validate("12.0.1", "abcd")
	require.NotNil(t, err)
	require.False(t, result)
}

func TestSemverFieldVersionValidator_Validate_Error_FieldValueNotSemver(t *testing.T) {
	result, err := semverValidator.Validate("abcdd", ">=12.0.0")
	require.NotNil(t, err)
	require.False(t, result)
}

func TestExactMatchFieldVersionValidator_Validate_ValidValue(t *testing.T) {
	exactMatchValidator.Validate("abcd", "abcd")
}

func TestExactMatchFieldVersionValidator_Validate_InvalidValue(t *testing.T) {
	exactMatchValidator.Validate("abcde", "abcd")
}
