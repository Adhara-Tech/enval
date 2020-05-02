package manifestchecker_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
	"github.com/Adhara-Tech/enval/pkg/model"
)

func TestVersionCheckerManager_CheckVersion_RegexpVersionParser(t *testing.T) {
	versionValidators := map[string]manifestchecker.FieldVersionValidator{
		"semver": manifestchecker.SemverFieldVersionValidator{},
	}
	fieldVersionValidatorManager := manifestchecker.NewFieldVersionValidatorManager(versionValidators)
	versionCheckerManager := manifestchecker.NewVersionCheckerManager(fieldVersionValidatorManager)
	versionFields := []manifestchecker.FieldSpec{
		manifestchecker.FieldSpec{
			Name:     "version",
			Type:     "semver",
			Required: true,
		},
	}

	versionCommandOutput := `openjdk 11.0.7 2020-04-14
OpenJDK Runtime Environment 18.9 (build 11.0.7+10)
OpenJDK 64-Bit Server VM 18.9 (build 11.0.7+10, mixed mode)`

	testData := []struct {
		TestName             string
		VersionParserSpec    manifestchecker.VersionParserSpec
		VersionCommandOutput string
		MustBeValid          bool
		MustBeError          bool
		ExpectedVersions     map[string]string
		VersionsFound        map[string]string
	}{
		{
			TestName: "Single field search found and valid",
			VersionParserSpec: manifestchecker.VersionParserSpec{
				Type:   "regexp",
				Regexp: `openjdk (?P<version>\d+\.\d+\.\d+)`,
			},
			VersionCommandOutput: versionCommandOutput,
			ExpectedVersions:     map[string]string{"version": ">=11.0.7"},
			VersionsFound:        map[string]string{"version": "11.0.7"},
			MustBeValid:          true,
			MustBeError:          false,
		},
		{
			TestName: "Single field search found but invalid",
			VersionParserSpec: manifestchecker.VersionParserSpec{
				Type:   "regexp",
				Regexp: `openjdk (?P<version>\d+\.\d+\.\d+) `,
			},
			VersionCommandOutput: versionCommandOutput,
			ExpectedVersions:     map[string]string{"version": "<11.0.7"},
			VersionsFound:        map[string]string{"version": "11.0.7"},
			MustBeValid:          false,
			MustBeError:          false,
		},
		{
			TestName: "Single field search not found",
			VersionParserSpec: manifestchecker.VersionParserSpec{
				Type:   "regexp",
				Regexp: `java (?P<version>\d+\.\d+\.\d+)`,
			},
			VersionCommandOutput: versionCommandOutput,
			ExpectedVersions:     map[string]string{"version": ">=11.0.7"},
			VersionsFound:        map[string]string{},
			MustBeValid:          false,
			MustBeError:          true,
		},
	}

	for _, currentTestData := range testData {
		t.Run(currentTestData.TestName, func(t *testing.T) {

			versionCheckerSpec := manifestchecker.VersionCheckerSpec{
				VersionParserArr: []manifestchecker.VersionParserSpec{
					currentTestData.VersionParserSpec,
				},
				Fields: versionFields,
			}

			manifestTool := model.ManifestTool{
				Name:   "test",
				Flavor: nil,
				Checks: currentTestData.ExpectedVersions,
			}

			checkVersionResult, err := versionCheckerManager.CheckVersion(versionCheckerSpec, versionCommandOutput, manifestTool)
			require.Nil(t, err)
			require.NotNil(t, checkVersionResult)

			for key, value := range currentTestData.VersionsFound {
				versionFound, ok := checkVersionResult.VersionsFound[key]
				require.True(t, ok)
				require.Equal(t, value, versionFound)
				require.Equal(t, currentTestData.MustBeValid, checkVersionResult.IsVersionValid)
				require.Empty(t, checkVersionResult.CheckVersionErrorMessage)
			}
		})
	}

}
