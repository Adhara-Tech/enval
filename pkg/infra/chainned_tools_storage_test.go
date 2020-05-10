package infra_test

import (
	"strings"
	"testing"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/infra"
	"github.com/stretchr/testify/require"
)

func buildToolsStorageChain() *infra.ToolsStorageChain {
	storage := infra.NewToolsStorageChain()

	fileSystemStorage := infra.NewFileSystemToolsStorage(customToolsSpecDirectory)
	storage.Add(fileSystemStorage)

	boxedStorage := infra.NewPackrBoxedToolsStorage()
	storage.Add(boxedStorage)

	return storage
}

func TestToolsStorageChain_Find_OK_Override(t *testing.T) {

	storage := buildToolsStorageChain()

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: envalTool})
	require.Nil(t, err)
	require.NotNil(t, toolSpec)
	require.Equal(t, envalTool, toolSpec.Name)
	// Check that enval custom definition has precedence over the packed tool spec
	require.True(t, strings.Contains(toolSpec.Description, "override"))
}

func TestToolsStorageChain_Find_OK(t *testing.T) {

	storage := buildToolsStorageChain()

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: javaTool})
	require.Nil(t, err)
	require.NotNil(t, toolSpec)
	require.Equal(t, javaTool, toolSpec.Name)
}

func TestToolsStorageChain_Find_NotFound(t *testing.T) {
	storage := buildToolsStorageChain()

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: "random"})
	require.NotNil(t, err)
	require.Nil(t, toolSpec)
	require.True(t, adapters.IsToolNotFoundExError(err))
}
