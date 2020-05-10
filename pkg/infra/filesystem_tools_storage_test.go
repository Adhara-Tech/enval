package infra_test

import (
	"testing"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/infra"
	"github.com/stretchr/testify/require"
)

const (
	customToolsSpecDirectory = "testdata/custom-tool-specs"
	envalTool                = "enval"
)

func TestFileSystemToolsStorage_Find_OK(t *testing.T) {

	storage := infra.NewFileSystemToolsStorage(customToolsSpecDirectory)

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: envalTool})
	require.Nil(t, err)
	require.NotNil(t, toolSpec)
	require.Equal(t, envalTool, toolSpec.Name)
}

func TestFileSystemToolsStorage_Find_NotFound(t *testing.T) {
	storage := infra.NewFileSystemToolsStorage(customToolsSpecDirectory)

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: "random"})
	require.NotNil(t, err)
	require.Nil(t, toolSpec)
	require.True(t, adapters.IsToolNotFoundExError(err))
}
