package infra_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/infra"
	"github.com/stretchr/testify/require"
)

func createCustomSpec(t *testing.T, content string) (string, func()) {
	dir, err := ioutil.TempDir("", "custom-specs")
	require.NoError(t, err)

	tmpfn := filepath.Join(dir, "tmpfile")
	err = ioutil.WriteFile(tmpfn, []byte(content), 0666)
	require.NoError(t, err)

	return dir, func() { os.RemoveAll(dir) }
}

func TestWithCustomSpecs(t *testing.T) {
	expectedName := "some-name-for-testing"
	customSpecsPath, cleanup := createCustomSpec(t, "name: "+expectedName)
	defer cleanup()

	storage := infra.NewDefaultToolsStorage()
	storage = storage.WithCustomSpecs(customSpecsPath)

	tool, err := storage.Find(adapters.ToolFindOptions{
		Name: expectedName,
	})

	require.NoError(t, err)
	require.NotNil(t, tool)
	require.Equal(t, expectedName, tool.Name)
}
