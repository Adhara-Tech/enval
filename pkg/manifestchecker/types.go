package manifestchecker

import "github.com/Adhara-Tech/enval/pkg/model"

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

func (t ToolSpec) findFlavor(flavor string) (*FlavorSpec, bool) {
	for index := range t.Flavors {
		currentFlavor := t.Flavors[index]
		if currentFlavor.Name == flavor {
			return &currentFlavor, true

		}
	}

	return nil, false
}

func (t ToolSpec) consolidateVersionCommandArgsForFlavor(flavor *string) ([]string, bool) {
	if flavor == nil {
		return t.VersionCommandArgs, true
	}

	currentFlavor, ok := t.findFlavor(*flavor)
	if !ok {
		return nil, ok
	}

	if currentFlavor.VersionCommandArgs != nil {
		return currentFlavor.VersionCommandArgs, true
	} else {
		return t.VersionCommandArgs, true
	}
}

func (t ToolSpec) consolidateVersionCheckerForFlavor(flavor *string) (*VersionCheckerSpec, bool) {
	if flavor == nil {
		return t.VersionChecker, true
	}

	currentFlavor, ok := t.findFlavor(*flavor)
	if !ok {
		return nil, ok
	}

	if currentFlavor.VersionChecker != nil {
		return currentFlavor.VersionChecker, true
	} else {
		return t.VersionChecker, true
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
	Tool               model.ManifestTool
	IsToolAvailable    bool
	VersionsFound      map[string]string
	VersionValidations map[string]bool
	IsVersionValid     bool
	IsError            bool
	Error              string
}

func ToolValidationResultFor(tool model.ManifestTool) *ToolValidationResult {
	return &ToolValidationResult{
		Tool:               tool,
		IsToolAvailable:    false,
		VersionsFound:      make(map[string]string),
		VersionValidations: make(map[string]bool),
		IsVersionValid:     false,
		IsError:            false,
		Error:              "",
	}
}

func (result *ToolValidationResult) WithToolAvailable(available bool) *ToolValidationResult {
	result.IsVersionValid = available

	return result
}

func (result *ToolValidationResult) WithError(msg string) *ToolValidationResult {
	result.IsError = true
	result.Error = msg
	return result
}

func (result *ToolValidationResult) AddVersionValidation(field string, valueFound string, valid bool) *ToolValidationResult {
	result.VersionsFound[field] = valueFound
	result.VersionValidations[field] = valid

	return result
}
