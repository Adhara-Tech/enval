package adapters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Adhara-Tech/enval/pkg/adapters"
)

func TestDefaultSystemAdapter_CheckCommandAvailable(t *testing.T) {
	systemAdapter := adapters.NewDefaultSystemAdapter()

	ok, err := systemAdapter.CheckCommandAvailable("go")

	require.Nil(t, err)
	require.True(t, ok)
}

func TestDefaultSystemAdapter_GetCommandVersionOutput(t *testing.T) {
	systemAdapter := adapters.NewDefaultSystemAdapter()

	versionOutput, err := systemAdapter.ExecuteCommand("go", []string{"version"})
	require.Nil(t, err)
	require.NotEmpty(t, versionOutput)
}
