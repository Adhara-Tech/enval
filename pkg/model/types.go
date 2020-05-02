package model

type Tool struct {
	Name               string          `yaml:"name"`
	Description        string          `yaml:"description"`
	Command            string          `yaml:"command"`
	VersionCommandArgs []string        `yaml:"version_command_args"`
	VersionChecker     *VersionChecker `yaml:"version_checker"`
	Flavors            []Flavor        `yaml:"flavors"`
}

//
//func (t Tool) findFlavor(flavor string) (*Flavor, error) {
//	for index := range t.Flavors {
//		currentFlavor := t.Flavors[index]
//		if currentFlavor.Name == flavor {
//			return &currentFlavor, nil
//
//		}
//	}
//
//	return nil, fmt.Errorf("flavor [%s] not found in tool [%s]", flavor, t.Name)
//}
//
//func (t Tool) ConsolidateVersionCommandArgsFor(flavor *string) ([]string, error) {
//	if flavor == nil {
//		return t.VersionCommandArgs, nil
//	}
//
//	currentFlavor, err := t.findFlavor(*flavor)
//	if err != nil {
//		return nil, err
//	}
//
//	if currentFlavor.VersionCommandArgs != nil {
//		return currentFlavor.VersionCommandArgs, nil
//	} else {
//		return t.VersionCommandArgs, nil
//	}
//}
//
//func (t Tool) ConsolidateVersionChecker(flavor *string) ([]*VersionChecker, error) {
//	versionCheckerArr := make([]*VersionChecker, 0)
//	if flavor == nil {
//		// Default version checker as first option
//		addDefaultChecker := false
//
//		if len(t.Flavors) > 0 {
//			for _, flavor := range t.Flavors {
//				if flavor.VersionChecker != nil {
//					versionCheckerArr = append(versionCheckerArr, flavor.VersionChecker)
//				} else {
//					addDefaultChecker = true
//				}
//			}
//			if addDefaultChecker {
//				if t.VersionChecker != nil {
//					versionCheckerArr = append(versionCheckerArr, t.VersionChecker)
//				} else {
//					return nil, fmt.Errorf("checker not found for tool [%s] and nil flavor. Default version checker needed because a flavor does not define its own one. Check tool definition", t.Name)
//				}
//			}
//			return versionCheckerArr, nil
//		} else if t.VersionChecker != nil {
//			versionCheckerArr = append(versionCheckerArr, t.VersionChecker)
//			return versionCheckerArr, nil
//		} else {
//			return nil, fmt.Errorf("checker not found for tool [%s] and nil flavor. Check tool definition", t.Name)
//		}
//
//	} else {
//
//		currentFlavor, err := t.findFlavor(*flavor)
//		if err != nil {
//			return nil, err
//		}
//
//		if currentFlavor.VersionChecker != nil {
//			versionCheckerArr = append(versionCheckerArr, currentFlavor.VersionChecker)
//		} else {
//			versionCheckerArr = append(versionCheckerArr, t.VersionChecker)
//		}
//
//		return versionCheckerArr, nil
//	}
//}

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
	Type        string
	Regexp      *string
	CodeSnippet *string
}

type ManifestTool struct {
	Name   string            `yaml:"name"`
	Flavor *string           `yaml:"flavor"`
	Checks map[string]string `yaml:"checks"`
}

func (mt ManifestTool) IsFlavoredCheck() bool {
	return mt.Flavor != nil
}
