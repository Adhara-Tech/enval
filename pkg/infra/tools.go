package infra

import (
	"fmt"

	"github.com/Adhara-Tech/enval/pkg/exerrors"

	"github.com/Adhara-Tech/enval/pkg/manifestchecker"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/packr"
	_ "github.com/gobuffalo/packr"
)

var _ adapters.ToolsStorage = (*DefaultToolsStorage)(nil)

type DefaultToolsStorage struct {
	innerBox packr.Box
}

func NewDefaultToolsStorage() *DefaultToolsStorage {
	box := packr.NewBox("../../tool-specs")

	return &DefaultToolsStorage{innerBox: box}
}

func (storage DefaultToolsStorage) Find(toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {

	for _, currentToolSpec := range storage.innerBox.List() {
		tool := &manifestchecker.ToolSpec{}
		toolSpecBytes, err := storage.innerBox.Find(currentToolSpec)
		if err != nil {
			return nil, exerrors.Wrap(err)
		}
		err = yaml.Unmarshal(toolSpecBytes, tool)
		if err != nil {
			return nil, err
		}

		if tool.Name == toolsFindOptions.Name {
			return tool, nil

		}
	}

	return nil, exerrors.New(fmt.Sprintf("tool with name [%s] not found", toolsFindOptions.Name))
}
