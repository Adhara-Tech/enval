package adapters

import (
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
	"github.com/Adhara-Tech/enval/pkg/model"
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
	Find(toolsFindOptions ToolFindOptions) (*model.Tool, error)
}

func (adapter DefaultToolsStorageAdapter) Find(toolName string) (*model.Tool, error) {
	return adapter.toolsStorage.Find(ToolFindOptions{
		Name: toolName,
	})
}
