package adapters

import (
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

func (adapter DefaultToolsStorageAdapter) Find(toolName string) (*manifestchecker.ToolSpec, error) {
	return adapter.toolsStorage.Find(ToolFindOptions{
		Name: toolName,
	})
}
