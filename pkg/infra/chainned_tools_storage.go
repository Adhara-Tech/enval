package infra

import (
	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
)

var _ adapters.ToolsStorage = (*ToolsStorageChain)(nil)

type ToolsStorageChain struct {
	chain []adapters.ToolsStorage
}

func NewToolsStorageChain() *ToolsStorageChain {
	return &ToolsStorageChain{
		chain: make([]adapters.ToolsStorage, 0),
	}
}

func (storage *ToolsStorageChain) Add(toolStorage adapters.ToolsStorage) {
	storage.chain = append(storage.chain, toolStorage)
}

func (storage ToolsStorageChain) Find(toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {

	for _, toolStorage := range storage.chain {
		toolSpec, err := toolStorage.Find(toolsFindOptions)
		if err != nil {
			if adapters.IsToolNotFoundExError(err) {
				continue
			} else {
				return nil, err
			}
		}

		return toolSpec, nil
	}

	return nil, adapters.NewToolNotFoundExError(toolsFindOptions.Name)
}
