package infra

import (
	"fmt"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/model"

	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/packr"
	_ "github.com/gobuffalo/packr"
)

var _ adapters.ToolsStorage = (*DefaultToolsStorage)(nil)

type DefaultToolsStorage struct {
	innerBox packr.Box
}

func NewDefaultToolsStorage(toolsSpecPath string) *DefaultToolsStorage {
	box := packr.NewBox(toolsSpecPath)

	return &DefaultToolsStorage{innerBox: box}
}

func (storage DefaultToolsStorage) Find(toolsFindOptions adapters.ToolFindOptions) (*model.Tool, error) {

	tool := &model.Tool{}
	for _, currentToolSpec := range storage.innerBox.List() {

		toolSpecBytes, err := storage.innerBox.Find(currentToolSpec)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(toolSpecBytes, tool)
		if err != nil {
			return nil, err
		}

		if tool.Name == toolsFindOptions.Name {
			return tool, nil

		}
	}

	return nil, fmt.Errorf("tool with name [%s] not found", toolsFindOptions.Name)
}
