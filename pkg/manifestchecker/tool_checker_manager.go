package manifestchecker

import (
	"fmt"

	"github.com/Adhara-Tech/enval/pkg/exerrors"
	"github.com/Adhara-Tech/enval/pkg/model"
)

type ToolsManager struct {
	toolsStorageAdapter   ToolsStorageAdapter
	systemAdapter         SystemAdapter
	versionCheckerManager VersionCheckerManager
}

type SystemAdapter interface {
	CheckCommandAvailable(command string) (bool, error)
	// Returns true if the directory exists. Returns false if the path does not exist or if it is a file
	CheckDirExist(path string) (bool, error)
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
	return tm.ValidateManifestAndNotify(manifest, func(_ []ToolValidationResult) {

	})
}

func (tm ToolsManager) ValidateManifestAndNotify(manifest model.Manifest, notifyFunc func(toolValidationResult []ToolValidationResult)) ([]ToolValidationResult, error) {
	result := make([]ToolValidationResult, 0)
	for _, manifestTool := range manifest.Tools {
		validationResultArr, err := tm.ValidateTool(manifestTool)
		if err != nil {
			return nil, err
		}
		result = append(result, validationResultArr...)
		notifyFunc(validationResultArr)
	}

	return result, nil
}

func (tm ToolsManager) IsManifestCompliant(manifest model.Manifest) (bool, error) {
	if manifest.CustomSpecs != "" {
		ok, err := tm.systemAdapter.CheckDirExist(manifest.CustomSpecs)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, exerrors.New(fmt.Sprintf("Custom specs directory [%s] defined in manifest does not exist or it is not a directory", manifest.CustomSpecs), exerrors.InvalidCustomSpecsDirEnvalErrorKind)
		}

	}

	for _, manifestTool := range manifest.Tools {
		_, err := tm.toolsStorageAdapter.Find(manifestTool.Name)
		if err != nil {
			return false, err
		}
	}

	return true, nil
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
func (tm ToolsManager) ValidateTool(manifestTool model.ManifestTool) ([]ToolValidationResult, error) {

	toolToCheck, err := tm.toolsStorageAdapter.Find(manifestTool.Name)
	if err != nil {
		return nil, err
	}

	result := make([]ToolValidationResult, 0)

	if manifestTool.IsFlavoredCheck() { // * Manifest has flavor configured
		// Get combined version checker spec and use it
		versionChecker, err := toolToCheck.consolidateVersionCheckerForFlavor(manifestTool.Flavor)
		if err != nil {
			return nil, err
		}

		toolValidationResult, err := tm.doValidate(manifestTool, toolToCheck, versionChecker, manifestTool.Flavor)
		if err != nil {
			return nil, err
		}

		result = append(result, *toolValidationResult)
		return result, nil

	} else { // * Manifest DOESN'T have flavor configured

		if toolToCheck.HasFlavors() { // ** Has the tool flavors
			// *** yes: Iterate all flavors
			for _, flavor := range toolToCheck.Flavors {
				// **** if flavor Version Checker match return
				flavorVersionCheckerSpec, err := toolToCheck.consolidateVersionCheckerForFlavor(&flavor.Name)
				if err != nil {
					return nil, err
				}
				toolValidationResult, err := tm.doValidate(manifestTool, toolToCheck, flavorVersionCheckerSpec, &flavor.Name)
				if err != nil {
					return nil, err
				}

				result = append(result, *toolValidationResult)

				if toolValidationResult.IsVersionValid {

					return result, nil
				}
				// **** Version checker that matches not found: Return not match result with all the checks done
			}

			return result, nil

		} else {
			// *** no: use main tool version checker to Check version
			toolValidationResult, err := tm.doValidate(manifestTool, toolToCheck, toolToCheck.VersionChecker, nil)
			if err != nil {
				return nil, err
			}

			result = append(result, *toolValidationResult)
			return result, nil
		}

	}
}
func (tm ToolsManager) doValidate(manifestTool model.ManifestTool, toolToCheck *ToolSpec, versionChecker *VersionCheckerSpec, flavor *string) (*ToolValidationResult, error) {
	versionCommandOutput, err := tm.executeVersionCommand(*toolToCheck, flavor)
	if err != nil {
		return nil, err
	}

	if !versionCommandOutput.IsToolAvailable {
		return ToolValidationResultFor(manifestTool).ToolNotAvailable(), nil
	}

	toolValidationResult, err := tm.versionCheckerManager.CheckVersion(CheckVersionRequest{
		VersionCheckerSpec:   *versionChecker,
		VersionCommandOutput: versionCommandOutput.rawVersionCommandOutput,
		ManifestTool:         manifestTool,
	})

	if err != nil {
		return nil, err
	}

	return toolValidationResult, nil
}

func (tm ToolsManager) executeVersionCommand(tool ToolSpec, flavor *string) (*executionVersionCommandResult, error) {

	commandAvailable, err := tm.systemAdapter.CheckCommandAvailable(tool.Command)
	if err != nil {
		return nil, err
	}

	if !commandAvailable {
		return &executionVersionCommandResult{
			IsToolAvailable:         false,
			rawVersionCommandOutput: "",
		}, nil
	}
	commandArgs, err := tool.consolidateVersionCommandArgsForFlavor(flavor)
	if err != nil {
		return nil, err
	}

	rawOutput, err := tm.systemAdapter.ExecuteCommand(tool.Command, commandArgs)

	if err != nil {
		return nil, err
	}

	return &executionVersionCommandResult{
		IsToolAvailable:         true,
		rawVersionCommandOutput: rawOutput,
	}, nil

}
