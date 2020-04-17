package adapters

import (
	"Adhara-Tech/check-my-setup/pkg/manifestchecker"
	"Adhara-Tech/check-my-setup/pkg/model"
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
