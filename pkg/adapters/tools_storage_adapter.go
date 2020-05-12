package adapters

import (
	"fmt"

	"github.com/Adhara-Tech/enval/pkg/exerrors"
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
)

var _ manifestchecker.ToolsStorageAdapter = (*DefaultToolsStorageAdapter)(nil)

type DefaultToolsStorageAdapter struct {
	toolsStorage ToolsStorage
}

func NewDefaultStorageAdapter(storage ToolsStorage) *DefaultToolsStorageAdapter {
	return &DefaultToolsStorageAdapter{toolsStorage: storage}
}

type ToolFindOptions struct {
	Name string
}

type ToolsStorage interface {
	Find(toolsFindOptions ToolFindOptions) (*manifestchecker.ToolSpec, error)
}

func NewToolNotFoundExError(toolName string) error {
	return exerrors.New(fmt.Sprintf("tool with name [%s] not found", toolName), exerrors.ToolDefinitionNotFoundEnvalErrorKind)
}

func IsToolNotFoundExError(err error) bool {
	return exerrors.IsEnvalErrorWithKind(err, exerrors.ToolDefinitionNotFoundEnvalErrorKind)
}

func (adapter DefaultToolsStorageAdapter) Find(toolName string) (*manifestchecker.ToolSpec, error) {
	return adapter.toolsStorage.Find(ToolFindOptions{
		Name: toolName,
	})
}
