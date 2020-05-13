package manifestchecker

import (
	"fmt"

	"github.com/Adhara-Tech/enval/pkg/exerrors"

	"github.com/Adhara-Tech/enval/pkg/model"
)

type versionChecker struct {
	versionParserArr   []VersionParser
	fields             []FieldSpec
	versionCheckerSpec VersionCheckerSpec
}

type VersionCheckerManager struct {
	fieldVersionValidatorManager FieldVersionValidatorManager
}

func NewVersionCheckerManager(fieldVersionValidatorManager FieldVersionValidatorManager) VersionCheckerManager {
	return VersionCheckerManager{
		fieldVersionValidatorManager: fieldVersionValidatorManager,
	}
}

type CheckVersionResult struct {
	IsVersionValid           bool
	VersionsFound            map[string]string
	CheckVersionErrorMessage string
}

type CheckVersionRequest struct {
	VersionCheckerSpec   VersionCheckerSpec
	VersionCommandOutput string
	ManifestTool         model.ManifestTool
}

func (versionCheckerManager VersionCheckerManager) buildVersionChecker(checkVersionRequest CheckVersionRequest) (*versionChecker, error) {
	versionCheckerSpec := checkVersionRequest.VersionCheckerSpec

	versionParsers := make([]VersionParser, len(versionCheckerSpec.VersionParserArr))

	for index, parser := range versionCheckerSpec.VersionParserArr {
		if parser.Type == "regexp" {
			regexpVersionParser := NewRegexVersionParser(parser.Regexp, versionCheckerSpec.FieldNames())
			versionParsers[index] = regexpVersionParser
		} else {
			return nil, exerrors.New(fmt.Sprintf("unknown parser type [%s]", parser.Type), exerrors.UnknownParserEnvalErrorKind)
		}
	}

	versionChecker := &versionChecker{
		versionParserArr:   versionParsers,
		fields:             versionCheckerSpec.Fields,
		versionCheckerSpec: versionCheckerSpec,
	}

	return versionChecker, nil
}

func (versionCheckerManager VersionCheckerManager) CheckVersion(checkVersionRequest CheckVersionRequest) (*ToolValidationResult, error) {
	versionChecker, err := versionCheckerManager.buildVersionChecker(checkVersionRequest)
	if err != nil {
		return nil, err
	}
	versionCommandOutput := checkVersionRequest.VersionCommandOutput
	manifestTool := checkVersionRequest.ManifestTool

	for _, currentParser := range versionChecker.versionParserArr {

		parsedVersionFields, err := currentParser.Parse(versionCommandOutput)
		if err != nil {
			if !IsUnsupportedInputRawVersionError(err) {
				return nil, err
			} else {
				continue
			}
		}

		toolValidationResult := ToolValidationResultFor(manifestTool)
		for fieldName, expectedVersion := range manifestTool.Checks {
			fieldSpec, ok := versionChecker.versionCheckerSpec.GetFieldSpecBy(fieldName)
			if !ok {
				toolValidationResult.InvalidField(fieldName, "", fmt.Sprintf("missing field spec for [%s]", fieldName))
				continue
			}

			fieldValue, ok := parsedVersionFields[fieldName]
			if !ok {
				toolValidationResult.InvalidField(fieldName, "", fmt.Sprintf("value not found for field [%s]", fieldName))
				continue
			}

			fieldVersionValidator, ok := versionCheckerManager.fieldVersionValidatorManager.FieldVersionValidator(fieldSpec)
			if !ok {
				toolValidationResult.InvalidField(fieldName, fieldValue, fmt.Sprintf("unknown field validation type [%s] for field[%s]", fieldSpec.Type, fieldName))
				continue
			}

			isFieldValid, err := fieldVersionValidator.Validate(fieldValue, expectedVersion)
			if err != nil {

				return nil, err
			}

			if !isFieldValid {
				toolValidationResult.InvalidField(fieldName, fieldValue, "")
			} else {
				toolValidationResult.ValidField(fieldName, fieldValue)
			}

		}

		return toolValidationResult, nil
	}

	return ToolValidationResultFor(manifestTool).NotParseableVersionOutputCommand(fmt.Sprintf("parsers configured for tool [%s] couldn't match command version output", manifestTool.Name)), nil
}
