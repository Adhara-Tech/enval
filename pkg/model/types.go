package model

import (
	"fmt"
)

type Tool struct {
	Name               string          `yaml:"name"`
	Description        string          `yaml:"description"`
	Command            string          `yaml:"command"`
	VersionCommandArgs []string        `yaml:"version_command_args"`
	VersionChecker     *VersionChecker `yaml:"version_checker"`
	Flavors            []Flavor        `yaml:"flavors"`
}

func (t Tool) findFlavor(flavor string) (*Flavor, error) {
	for index := range t.Flavors {
		currentFlavor := t.Flavors[index]
		if currentFlavor.Name == flavor {
			return &currentFlavor, nil

		}
	}

	return nil, fmt.Errorf("flavor [%s] not found in tool [%s]", flavor, t.Name)
}

func (t Tool) ConsolidateVersionCommandArgsFor(flavor *string) ([]string, error) {
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

func (t Tool) ConsolidateVersionChecker(flavor *string) (*VersionChecker, error) {
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

type VersionChecker struct {
	Parser VersionParser     `yaml:"parser"`
	Fields map[string]string `yaml:"fields"`
}

type Flavor struct {
	Name               string          `yaml:"name"`
	VersionCommandArgs []string        `yaml:"version_command_args"`
	VersionChecker     *VersionChecker `yaml:"version_checker"`
}

type RegexpChecker struct {
}

type CodeSnippet struct {
}

type Manifest struct {
	Name  string         `yaml:"name"`
	Tools []ManifestTool `yaml:"tools"`
}

type VersionParser struct {
	Regexp      *string
	CodeSnippet *string
}

type ManifestTool struct {
	Name   string            `yaml:"name"`
	Flavor *string           `yaml:"flavor"`
	Checks map[string]string `yaml:"checks"`
}
