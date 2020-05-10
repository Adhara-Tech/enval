package infra

import (
	"github.com/Adhara-Tech/enval/pkg/exerrors"

	"github.com/Adhara-Tech/enval/pkg/manifestchecker"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	_ "github.com/gobuffalo/packr"
)

var _ adapters.ToolsStorage = (*PackrBoxedToolsStorage)(nil)

type ListFinder interface {
	packd.Lister
	packd.Finder
}

type PackrBoxedToolsStorage struct {
	innerBox ListFinder
}

func NewPackrBoxedToolsStorage() *PackrBoxedToolsStorage {
	return &PackrBoxedToolsStorage{
		innerBox: packr.NewBox("../../tool-specs"),
	}
}

func findInBox(box ListFinder, toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {
	for _, currentToolSpec := range box.List() {
		toolSpecBytes, err := box.Find(currentToolSpec)
		if err != nil {
			return nil, exerrors.Wrap(err)
		}

		tool := &manifestchecker.ToolSpec{}
		err = yaml.Unmarshal(toolSpecBytes, tool)
		if err != nil {
			return nil, err
		}

		if tool.Name == toolsFindOptions.Name {
			return tool, nil

		}
	}

	return nil, nil
}

func (storage PackrBoxedToolsStorage) Find(toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {
	tool, err := findInBox(storage.innerBox, toolsFindOptions)
	if err != nil {
		return nil, err
	}
	if tool != nil {
		return tool, nil
	}

	return nil, adapters.NewToolNotFoundExError(toolsFindOptions.Name)
}
