package manifestchecker

import (
	"Adhara-Tech/check-my-setup/pkg/model"
	"errors"
	"fmt"
	"os/exec"

	"github.com/Masterminds/semver/v3"
)

type Checker interface {
	Check(manifest model.Manifest, notifier CheckNotifier) error
}
type Notification struct {
	Tool               model.ManifestTool
	IsToolAvailable    bool
	VersionsFound      map[string]string
	VersionValidations map[string]bool
	IsVersionValid     bool
	Error              string
}
type CheckNotifier func(notification Notification)

var _ Checker = (*DefaultManifestChecker)(nil)

type DefaultManifestChecker struct {
	toolsStorageAdapter ToolsStorageAdapter
}

func NewDefaultManifestChecker(toolsStorageAdapter ToolsStorageAdapter) *DefaultManifestChecker {
	return &DefaultManifestChecker{toolsStorageAdapter: toolsStorageAdapter}
}

type ToolsStorageAdapter interface {
	Find(toolName string) (*model.Tool, error)
}

type VersionParser interface {
	Parse(rawVersion string) (string, error)
}

type SystemAdapter interface {
}

func (checker DefaultManifestChecker) Check(manifest model.Manifest, notifier CheckNotifier) error {
	for _, currentToolConfig := range manifest.Tools {

		tool, err := checker.toolsStorageAdapter.Find(currentToolConfig.Name)

		if err != nil {
			return err
		}

		available, err := CheckCommandAvailable(tool.Command)
		if err != nil {
			return err
		}

		if !available {
			return errors.New("command not found. Check if available and added to path")
		}

		versionCommandArgs, err := tool.ConsolidateVersionCommandArgsFor(currentToolConfig.Flavor)
		if err != nil {
			return err
		}

		versionCommandOutputStr, err := GetCommandVersionOutput(tool.Command, versionCommandArgs)
		if err != nil {
			return err
		}

		//fmt.Println(versionCommandOutputStr)

		versionChecker, err := tool.ConsolidateVersionChecker(currentToolConfig.Flavor)
		if err != nil {
			return err
		}

		var versionFieldValues map[string]string
		if versionChecker.Parser.Regexp != nil {
			parser := NewRegexVersionParser(*versionChecker.Parser.Regexp, keySliceFrom(versionChecker.Fields))
			versionFieldValues, err = parser.Parse(versionCommandOutputStr)
			if err != nil {
				return err
			}

			//fmt.Println(versionFieldValues)
		}
		notification := Notification{
			Tool:               currentToolConfig,
			IsToolAvailable:    true,
			VersionsFound:      make(map[string]string),
			VersionValidations: make(map[string]bool),
			IsVersionValid:     true,
			Error:              "",
		}

		for fieldName, expectedVersion := range currentToolConfig.Checks {

			versionCheckType, ok := versionChecker.Fields[fieldName]
			if !ok {
				//TODO
				return fmt.Errorf("TODO versionchecktype")
			}
			fieldValue, ok := versionFieldValues[fieldName]
			notification.VersionsFound[fieldName] = fieldValue
			if !ok {
				//TODO
				return fmt.Errorf("TODO fieldValue")
			}

			if versionCheckType == "semver" {
				version, err := semver.NewVersion(fieldValue)
				if err != nil {
					return err
				}

				versionContraint, err := semver.NewConstraint(expectedVersion)
				if err != nil {
					return err
				}

				validVersion := versionContraint.Check(version)
				notification.VersionValidations[fieldName] = validVersion
				if !validVersion {
					notification.IsVersionValid = false
				}
				//	//TODO
				//	//return fmt.Errorf("TODO invalid version")
				//
				//}
			}
		}
		notifier(notification)
	}
	return nil
}

func keySliceFrom(keyMap map[string]string) []string {
	keyArr := make([]string, len(keyMap))
	counter := 0
	for key := range keyMap {
		keyArr[counter] = key
		counter++
	}

	return keyArr
}

func CheckCommandAvailable(command string) (bool, error) {

	_, err := exec.LookPath(command)
	if err != nil {
		return false, err
	}

	return true, nil

}

func GetCommandVersionOutput(commandName string, params []string) (string, error) {

	cmd := exec.Command(commandName, params...)

	versionString, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(versionString), nil
}
