package manifestchecker

import (
	"errors"

	"github.com/Adhara-Tech/enval/pkg/model"
)

type ToolsManager struct {
	toolsStorageAdapter   ToolsStorageAdapter
	systemAdapter         SystemAdapter
	versionCheckerManager VersionCheckerManager
}

type SystemAdapter interface {
	CheckCommandAvailable(command string) (bool, error)
	ExecuteCommand(commandName string, params []string) (string, error)
}

func NewToolsManager(toolsStorageAdapter ToolsStorageAdapter,
	systemAdapter SystemAdapter,
	versionCheckerManager VersionCheckerManager) ToolsManager {

	return ToolsManager{
		toolsStorageAdapter:   toolsStorageAdapter,
		systemAdapter:         systemAdapter,
		versionCheckerManager: versionCheckerManager,
	}
}

type ToolsStorageAdapter interface {
	Find(toolName string) (*ToolSpec, error)
}

func (tm ToolsManager) ValidateManifest(manifest model.Manifest) ([]ToolValidationResult, error) {
	return tm.ValidateManifestAndNotify(manifest, func(_ *ToolValidationResult) {

	})
}

func (tm ToolsManager) ValidateManifestAndNotify(manifest model.Manifest, notifyFunc func(toolValidationResult *ToolValidationResult)) ([]ToolValidationResult, error) {
	result := make([]ToolValidationResult, 0)
	for _, manifestTool := range manifest.Tools {
		validationResult, err := tm.ValidateTool(manifestTool)
		if err != nil {
			return nil, err
		}
		result = append(result, *validationResult)
		notifyFunc(validationResult)
	}

	return result, nil
}

// * Manifest has flavor configured
// ** Tool has flavor configured?
// *** yes: Check version
// *** no: Error not configured error
// * Manifest DOESN'T have flavor configured
// ** Has the tool flavors
// *** no: use main tool version checker to Check version
// *** yes: Iterate all flavors
// **** if flavor Version Checker match return
// **** Version checker that matches not found: Return not match result with all the checks done
func (tm ToolsManager) ValidateTool(manifestTool model.ManifestTool) (*ToolValidationResult, error) {

	toolToCheck, err := tm.toolsStorageAdapter.Find(manifestTool.Name)
	if err != nil {
		//TODO add error of type not found and use it
		return &ToolValidationResult{
			Tool:               manifestTool,
			IsToolAvailable:    false,
			VersionsFound:      nil,
			VersionValidations: nil,
			IsVersionValid:     false,
			Error:              "tool not configured",
		}, err
	}

	if manifestTool.IsFlavoredCheck() { // * Manifest has flavor configured
		// ** Tool has flavor configured?
		versionChecker, ok := toolToCheck.consolidateVersionCheckerForFlavor(manifestTool.Flavor)
		if ok {
			// *** yes: Check version of flavor
			versionCommandOutput, err := tm.executeVersionCommand(*toolToCheck, nil)
			if err != nil {
				//TODO
			}

			checkVersionResult, err := tm.versionCheckerManager.CheckVersion(*versionChecker, versionCommandOutput, manifestTool)
			if err != nil {
				return nil, err
			}

			return &ToolValidationResult{
				Tool:               manifestTool,
				IsToolAvailable:    true,
				VersionsFound:      checkVersionResult.VersionsFound,
				VersionValidations: nil,
				IsVersionValid:     checkVersionResult.IsVersionValid,
				Error:              "",
			}, nil
		} else {
			// *** no: Error not configured error
			//TODO add error of type not found and use it
			return &ToolValidationResult{
				Tool:               manifestTool,
				IsToolAvailable:    false,
				VersionsFound:      nil,
				VersionValidations: nil,
				IsVersionValid:     false,
				Error:              "flavor not configured",
			}, err
		}
	} else { // * Manifest DOESN'T have flavor configured

		if toolToCheck.HasFlavors() { // ** Has the tool flavors
			// *** yes: Iterate all flavors
			// **** if flavor Version Checker match return
			// **** Version checker that matches not found: Return not match result with all the checks done
		} else {
			// *** no: use main tool version checker to Check version
			versionCommandOutput, err := tm.executeVersionCommand(*toolToCheck, nil)
			if err != nil {
				//TODO
			}

			checkVersionResult, err := tm.versionCheckerManager.CheckVersion(*toolToCheck.VersionChecker, versionCommandOutput, manifestTool)
			if err != nil {
				return nil, err
			}

			return &ToolValidationResult{
				Tool:               manifestTool,
				IsToolAvailable:    true,
				VersionsFound:      checkVersionResult.VersionsFound,
				VersionValidations: nil,
				IsVersionValid:     checkVersionResult.IsVersionValid,
				Error:              "",
			}, nil
		}

	}

	return nil, errors.New("abnormal validation ending")
}

func (tm ToolsManager) executeVersionCommand(tool ToolSpec, flavor *string) (string, error) {

	commandAvailable, err := tm.systemAdapter.CheckCommandAvailable(tool.Command)
	if err != nil {
		// TODO
	}
	if !commandAvailable {
		// TODO
	}
	commandArgs, ok := tool.consolidateVersionCommandArgsForFlavor(flavor)
	if !ok {
		// TODO
	}

	return tm.systemAdapter.ExecuteCommand(tool.Command, commandArgs)
}
