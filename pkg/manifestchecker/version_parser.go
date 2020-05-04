package manifestchecker

import (
	"fmt"
	"strings"

	"github.com/Adhara-Tech/enval/pkg/exerrors"
)

// VersionParser interface wraps the Parse method
//
// Parse parses the rawVersion string typically obtained by executing the version flag of a command
// and returns a map of key values obtained as result of parsing the input or error
// When the parser does not match the input a UnsupportedInputRawVersionError error must be returned
type VersionParser interface {
	Parse(rawVersion string) (map[string]string, error)
}

const unsupportedInputRawVersion = "unsupported input raw version"

// NewUnsupportedInputRawVersionError creates a new UnsupportedInputRawVersionError
func NewUnsupportedInputRawVersionError(msg string) error {
	return exerrors.New(fmt.Sprint(unsupportedInputRawVersion, ". ", msg))
}

// IsUnsupportedInputRawVersionError checks if the given error is of type UnsupportedInputRawVersionError
func IsUnsupportedInputRawVersionError(err error) bool {
	return strings.HasPrefix(err.Error(), unsupportedInputRawVersion)
}
