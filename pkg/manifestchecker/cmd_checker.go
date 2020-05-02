package manifestchecker

//
//type Checker interface {
//	Check(manifest model.Manifest, notifier CheckNotifier) (*CheckResult, error)
//}
//type Notification struct {
//	Tool               model.ManifestTool
//	IsToolAvailable    bool
//	VersionsFound      map[string]string
//	VersionValidations map[string]bool
//	IsVersionValid     bool
//	Error              string
//}
//type CheckNotifier func(notification Notification)
//
//var _ Checker = (*DefaultManifestChecker)(nil)
//
//type DefaultManifestChecker struct {
//	toolsStorageAdapter ToolsStorageAdapter
//	systemAdapter       SystemAdapter
//}
//
//func NewDefaultManifestChecker(toolsStorageAdapter ToolsStorageAdapter, systemAdapter SystemAdapter) *DefaultManifestChecker {
//	return &DefaultManifestChecker{
//		toolsStorageAdapter: toolsStorageAdapter,
//		systemAdapter:       systemAdapter,
//	}
//}
//
//type ToolsStorageAdapter interface {
//	Find(toolName string) (*model.Tool, error)
//}
//
//type SystemAdapter interface {
//	CheckCommandAvailable(command string) (bool, error)
//	ExecuteCommand(commandName string, params []string) (string, error)
//}
//
//type CheckResult struct {
//	Ok      bool
//	Message string
//}
//
//func ok() *CheckResult {
//	return &CheckResult{
//		Ok: true,
//	}
//}
//
//func ko(message string) *CheckResult {
//	return &CheckResult{
//		Ok:      false,
//		Message: message,
//	}
//}
//
//func (checker DefaultManifestChecker) Check(manifest model.Manifest, notifier CheckNotifier) (*CheckResult, error) {
//	//for _, currentToolConfig := range manifest.Tools {
//	//
//	//	tool, err := checker.toolsStorageAdapter.Find(currentToolConfig.Name)
//	//
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//
//	//	available, err := checker.systemAdapter.CheckCommandAvailable(tool.Command)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//
//	//	if !available {
//	//		return ko("command not found. Command must be in path"), nil
//	//	}
//	//
//	//	versionCommandArgs, err := tool.ConsolidateVersionCommandArgsFor(currentToolConfig.Flavor)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//
//	//	versionCommandOutputStr, err := checker.systemAdapter.ExecuteCommand(tool.Command, versionCommandArgs)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//
//	//	//fmt.Println(versionCommandOutputStr)
//	//
//	//	versionCheckerArr, err := tool.ConsolidateVersionChecker(currentToolConfig.Flavor)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//
//	//	for _, versionChecker := range versionCheckerArr {
//	//		var versionFieldValues map[string]string
//	//		if versionChecker.Parser.Regexp != nil {
//	//			parser := NewRegexVersionParser(*versionChecker.Parser.Regexp, keySliceFrom(versionChecker.Fields))
//	//			versionFieldValues, err = parser.Parse(versionCommandOutputStr)
//	//			if err != nil {
//	//				return nil, err
//	//			}
//	//
//	//			//fmt.Println(versionFieldValues)
//	//		}
//	//	}
//	//	notification := Notification{
//	//		Tool:               currentToolConfig,
//	//		IsToolAvailable:    true,
//	//		VersionsFound:      make(map[string]string),
//	//		VersionValidations: make(map[string]bool),
//	//		IsVersionValid:     true,
//	//		Error:              "",
//	//	}
//	//
//	//	for fieldName, expectedVersion := range currentToolConfig.Checks {
//	//
//	//		versionCheckType, ok := versionChecker.Fields[fieldName]
//	//		if !ok {
//	//
//	//			return ko(fmt.Sprintf("Check configured for unknown version field [%s] in tool [%s]", fieldName, currentToolConfig.Name)), nil
//	//		}
//	//		fieldValue, ok := versionFieldValues[fieldName]
//	//		notification.VersionsFound[fieldName] = fieldValue
//	//		if !ok {
//	//
//	//			return ko(fmt.Sprintf("Value not found for field with name [%s] in tool [%s]", fieldName, currentToolConfig.Name)), nil
//	//		}
//	//
//	//		if versionCheckType == "semver" {
//	//			version, err := semver.NewVersion(fieldValue)
//	//			if err != nil {
//	//				return nil, err
//	//			}
//	//
//	//			versionContraint, err := semver.NewConstraint(expectedVersion)
//	//			if err != nil {
//	//				return nil, err
//	//			}
//	//
//	//			validVersion := versionContraint.Check(version)
//	//			notification.VersionValidations[fieldName] = validVersion
//	//			if !validVersion {
//	//				notification.IsVersionValid = false
//	//			}
//	//
//	//		}
//	//	}
//	//	notifier(notification)
//	//
//	//}
//	return ok(), nil
//}

//func keySliceFrom(keyMap map[string]interface{}) []string {
//	keyArr := make([]string, len(keyMap))
//	counter := 0
//	for key := range keyMap {
//		keyArr[counter] = key
//		counter++
//	}
//
//	return keyArr
//}
