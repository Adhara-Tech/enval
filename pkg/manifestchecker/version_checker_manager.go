package manifestchecker

import (
	"fmt"

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

func (versionCheckerManager VersionCheckerManager) CheckVersion(versionCheckerSpec VersionCheckerSpec, versionCommandOutput string, manifestTool model.ManifestTool) (*CheckVersionResult, error) {
	versionParsers := make([]VersionParser, len(versionCheckerSpec.VersionParserArr))

	for index, parser := range versionCheckerSpec.VersionParserArr {
		if parser.Type == "regexp" {
			regexpVersionParser := NewRegexVersionParser(parser.Regexp, versionCheckerSpec.FieldNames())
			versionParsers[index] = regexpVersionParser
		} else {
			return nil, fmt.Errorf("unknown parser type [%s]", parser.Type)
		}
	}

	versionChecker := versionChecker{
		versionParserArr:   versionParsers,
		fields:             versionCheckerSpec.Fields,
		versionCheckerSpec: versionCheckerSpec,
	}

	return versionCheckerManager.doCheckVersion(versionChecker, versionCommandOutput, manifestTool)

}

func (versionCheckerManager VersionCheckerManager) doCheckVersion(versionChecker versionChecker, versionCommandOutput string, tool model.ManifestTool) (*CheckVersionResult, error) {

	for _, currentParser := range versionChecker.versionParserArr {

		parsedVersionFields, err := currentParser.Parse(versionCommandOutput)
		if err != nil {
			if !IsUnsupportedInputRawVersionError(err) {
				return nil, err
			} else {
				continue
			}
		}

		versionsFound := make(map[string]string)
		isVersionValid := true
		for fieldName, expectedVersion := range tool.Checks {
			fieldSpec, ok := versionChecker.versionCheckerSpec.GetFieldSpecBy(fieldName)
			if !ok {
				return &CheckVersionResult{
					IsVersionValid:           false,
					CheckVersionErrorMessage: fmt.Sprintf("unknown version field [%s]", fieldName),
				}, nil
			}

			fieldValue, ok := parsedVersionFields[fieldName]
			if !ok {
				return &CheckVersionResult{
					IsVersionValid:           false,
					CheckVersionErrorMessage: fmt.Sprintf(fmt.Sprintf("value not found for field [%s]", fieldName)),
				}, nil
			}

			fieldVersionValidator, ok := versionCheckerManager.fieldVersionValidatorManager.FieldVersionValidator(fieldSpec)
			if !ok {
				return &CheckVersionResult{
					IsVersionValid:           false,
					CheckVersionErrorMessage: fmt.Sprintf(fmt.Sprintf("unknown field validation type [%s] for field[%s]", fieldSpec.Type, fieldName)),
				}, nil
			}

			isFieldValid, err := fieldVersionValidator.Validate(fieldValue, expectedVersion)
			if err != nil {
				return nil, err
			}

			isVersionValid = isVersionValid && isFieldValid
			versionsFound[fieldName] = fieldValue
		}

		return &CheckVersionResult{
			IsVersionValid:           isVersionValid,
			VersionsFound:            versionsFound,
			CheckVersionErrorMessage: "",
		}, nil
	}

	return &CheckVersionResult{
		IsVersionValid:           false,
		CheckVersionErrorMessage: "parsers didn't match command version output",
	}, nil
}
