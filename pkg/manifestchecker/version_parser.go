package manifestchecker

import (
	"fmt"

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

// NewUnsupportedInputRawVersionError creates a new UnsupportedInputRawVersionError
func NewUnsupportedInputRawVersionError(msg string) error {
	return exerrors.New(fmt.Sprint("unsupported input raw version. ", msg), exerrors.UnsupportedInputRawVersionEnvalErrorKind)
}

// IsUnsupportedInputRawVersionError checks if the given error is of type UnsupportedInputRawVersionError
func IsUnsupportedInputRawVersionError(err error) bool {
	return exerrors.IsEnvalErrorWithKind(err, exerrors.UnsupportedInputRawVersionEnvalErrorKind)
}
