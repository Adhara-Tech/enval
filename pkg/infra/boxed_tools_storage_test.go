package infra_test

import (
	"testing"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/infra"

	"github.com/stretchr/testify/require"
)

const (
	javaTool = "java"
)

func TestPackrBoxedToolsStorage_Find_OK(t *testing.T) {
	storage := infra.NewPackrBoxedToolsStorage()

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: javaTool})
	require.Nil(t, err)
	require.NotNil(t, toolSpec)
	require.Equal(t, javaTool, toolSpec.Name)
}

func TestPackrBoxedToolsStorage_Find_NotFound(t *testing.T) {

	storage := infra.NewPackrBoxedToolsStorage()

	toolSpec, err := storage.Find(adapters.ToolFindOptions{Name: "random"})
	require.NotNil(t, err)
	require.Nil(t, toolSpec)
	require.True(t, adapters.IsToolNotFoundExError(err))
}
