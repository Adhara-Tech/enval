package infra

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/exerrors"
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
	yaml "gopkg.in/yaml.v2"
)

var _ adapters.ToolsStorage = (*FileSystemToolsStorage)(nil)

type FileSystemToolsStorage struct {
	path string
}

func NewFileSystemToolsStorage(path string) *FileSystemToolsStorage {
	return &FileSystemToolsStorage{
		path: path,
	}
}

func (storage FileSystemToolsStorage) findInFileSystem(toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {
	toolSpecFiles := make([]string, 0)
	err := filepath.Walk(storage.path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			toolSpecFiles = append(toolSpecFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, exerrors.Wrap(err, exerrors.InternalEnvalErrorKind)
	}

	for _, currentToolSpec := range toolSpecFiles {
		toolSpecBytes, err := ioutil.ReadFile(currentToolSpec)
		if err != nil {
			return nil, exerrors.Wrap(err, exerrors.InternalEnvalErrorKind)
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

func (storage FileSystemToolsStorage) Find(toolsFindOptions adapters.ToolFindOptions) (*manifestchecker.ToolSpec, error) {
	tool, err := storage.findInFileSystem(toolsFindOptions)
	if err != nil {
		return nil, err
	}
	if tool != nil {
		return tool, nil
	}

	return nil, adapters.NewToolNotFoundExError(toolsFindOptions.Name)
}
