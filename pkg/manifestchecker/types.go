package manifestchecker

import (
	"fmt"

	"github.com/Adhara-Tech/enval/pkg/exerrors"
	"github.com/Adhara-Tech/enval/pkg/model"
)

type ToolSpec struct {
	Name               string              `yaml:"name"`
	Description        string              `yaml:"description"`
	Command            string              `yaml:"command"`
	VersionCommandArgs []string            `yaml:"version_command_args"`
	VersionChecker     *VersionCheckerSpec `yaml:"version_checker"`
	Flavors            []FlavorSpec        `yaml:"flavors"`
}

func (t ToolSpec) HasFlavors() bool {
	return len(t.Flavors) > 0
}

func (t ToolSpec) findFlavor(flavor string) (*FlavorSpec, error) {
	for index := range t.Flavors {
		currentFlavor := t.Flavors[index]
		if currentFlavor.Name == flavor {
			return &currentFlavor, nil

		}
	}

	return nil, exerrors.New(fmt.Sprintf("flavor [%s] not found for tool [%s]", flavor, t.Name), exerrors.FlavorDefinitionNotFoundEnvalErrorKind)
}

func (t ToolSpec) consolidateVersionCommandArgsForFlavor(flavor *string) ([]string, error) {
	if flavor == nil {
		return t.VersionCommandArgs, nil
	}

	currentFlavor, err := t.findFlavor(*flavor)
	if err != nil {
		return nil, err
	}

	if currentFlavor.VersionCommandArgs != nil {
		return currentFlavor.VersionCommandArgs, nil
	} else {
		return t.VersionCommandArgs, nil
	}
}

func (t ToolSpec) consolidateVersionCheckerForFlavor(flavor *string) (*VersionCheckerSpec, error) {
	if flavor == nil {
		return t.VersionChecker, nil
	}

	currentFlavor, err := t.findFlavor(*flavor)
	if err != nil {
		return nil, err
	}

	if currentFlavor.VersionChecker != nil {
		return currentFlavor.VersionChecker, nil
	} else {
		return t.VersionChecker, nil
	}
}

type FlavorSpec struct {
	Name               string              `yaml:"name"`
	VersionCommandArgs []string            `yaml:"version_command_args"`
	VersionChecker     *VersionCheckerSpec `yaml:"version_checker"`
}

type VersionCheckerSpec struct {
	VersionParserArr []VersionParserSpec `yaml:"parsers"`
	Fields           []FieldSpec         `yaml:"fields"`
}

func (v VersionCheckerSpec) FieldNames() []string {
	keyArr := make([]string, len(v.Fields))
	counter := 0
	for _, fieldSpec := range v.Fields {
		keyArr[counter] = fieldSpec.Name
		counter++
	}

	return keyArr
}

func (v VersionCheckerSpec) GetFieldSpecBy(name string) (FieldSpec, bool) {
	for _, fieldSpec := range v.Fields {
		if name == fieldSpec.Name {
			return fieldSpec, true
		}

	}
	return FieldSpec{}, false
}

type FieldSpec struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

type VersionParserSpec struct {
	Type   string `yaml:"type"`
	Regexp string `yaml:"regexp"`
}

type ToolValidationResult struct {
	Tool                            model.ManifestTool
	IsToolAvailable                 bool
	FieldValidations                map[string]FieldValidationResult
	IsVersionValid                  bool
	IsCommandVersionOutputParseable bool
	ResultDescription               string
}

func ToolValidationResultFor(tool model.ManifestTool) *ToolValidationResult {
	return &ToolValidationResult{
		Tool:                            tool,
		IsToolAvailable:                 true,
		FieldValidations:                make(map[string]FieldValidationResult),
		IsVersionValid:                  true,
		IsCommandVersionOutputParseable: true,
	}
}

type FieldValidationResult struct {
	FieldName         string
	ValueFound        string
	IsValid           bool
	ResultDescription string
}

func (result *ToolValidationResult) WithToolAvailable(available bool) *ToolValidationResult {
	result.IsVersionValid = available

	return result
}

func (result *ToolValidationResult) InvalidField(field string, valueFound string, resultDescription string) *ToolValidationResult {
	result.FieldValidations[field] = FieldValidationResult{
		FieldName:         field,
		ValueFound:        valueFound,
		IsValid:           false,
		ResultDescription: resultDescription,
	}

	result.IsVersionValid = false
	return result
}

func (result *ToolValidationResult) ValidField(field string, valueFound string) *ToolValidationResult {
	result.FieldValidations[field] = FieldValidationResult{
		FieldName:  field,
		ValueFound: valueFound,
		IsValid:    true,
	}

	return result
}

func (result *ToolValidationResult) ToolNotAvailable() *ToolValidationResult {
	result.IsToolAvailable = false
	return result
}

func (result *ToolValidationResult) NotParseableVersionOutputCommand(msg string) *ToolValidationResult {
	result.IsCommandVersionOutputParseable = false
	result.IsVersionValid = false
	result.ResultDescription = msg
	return result
}

func (result *ToolValidationResult) IsValid() bool {
	return result.IsToolAvailable && result.IsVersionValid && result.IsCommandVersionOutputParseable
}

type executionVersionCommandResult struct {
	IsToolAvailable         bool
	rawVersionCommandOutput string
}
