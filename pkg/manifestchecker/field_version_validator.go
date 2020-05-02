package manifestchecker

import "github.com/Masterminds/semver/v3"

type FieldVersionValidator interface {
	Validate(version string, expectedVersion string) (bool, error)
}

type FieldVersionValidatorManager struct {
	versionValidators map[string]FieldVersionValidator
}

func NewFieldVersionValidatorManager(versionValidators map[string]FieldVersionValidator) FieldVersionValidatorManager {
	return FieldVersionValidatorManager{
		versionValidators: versionValidators,
	}
}

func (versionValidatorManager FieldVersionValidatorManager) FieldVersionValidator(fieldSpec FieldSpec) (FieldVersionValidator, bool) {
	validator, ok := versionValidatorManager.versionValidators[fieldSpec.Type]

	return validator, ok
}

type SemverFieldVersionValidator struct {
}

func (validator SemverFieldVersionValidator) Validate(fieldValue string, expectedVersion string) (bool, error) {
	version, err := semver.NewVersion(fieldValue)
	if err != nil {
		return false, err
	}

	versionConstraint, err := semver.NewConstraint(expectedVersion)
	if err != nil {
		return false, err
	}

	return versionConstraint.Check(version), nil
}

type ExactMatchFieldVersionValidator struct {
}

func (validator ExactMatchFieldVersionValidator) Validate(fieldValue string, expectedVersion string) (bool, error) {
	return fieldValue == expectedVersion, nil
}
