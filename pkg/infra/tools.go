package infra

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/Adhara-Tech/enval/pkg/exerrors"

	"github.com/Adhara-Tech/enval/pkg/manifestchecker"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	_ "github.com/gobuffalo/packr"
)

var _ adapters.ToolsStorage = (*DefaultToolsStorage)(nil)

type ListFinder interface {
	packd.Lister
	packd.Finder
}

type DefaultToolsStorage struct {
	innerBox  ListFinder
	customBox ListFinder
}

func NewDefaultToolsStorage() *DefaultToolsStorage {
	return &DefaultToolsStorage{
		innerBox: packr.NewBox("../../tool-specs"),
	}
}

func (storage DefaultToolsStorage) WithCustomSpecs(customSpecsPath string) *DefaultToolsStorage {
	if !path.IsAbs(customSpecsPath) {
		customSpecsPath = filepath.Join("..", "..", customSpecsPath)
	}
	storage.customBox = packr.NewBox(customSpecsPath)
	return &storage
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

func (storage DefaultToolsStorage) Find(toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {
	if storage.customBox != nil {
		tool, err := findInBox(storage.customBox, toolsFindOptions)
		if err != nil {
			return nil, err
		}
		if tool != nil {
			return tool, nil
		}
	}

	tool, err := findInBox(storage.innerBox, toolsFindOptions)
	if err != nil {
		return nil, err
	}
	if tool != nil {
		return tool, nil
	}

	return nil, exerrors.New(fmt.Sprintf("tool with name [%s] not found", toolsFindOptions.Name))
}
