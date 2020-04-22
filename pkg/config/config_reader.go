package config

import (
	"io/ioutil"

	"github.com/Adhara-Tech/enval/pkg/model"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ReadManifestFrom(path string) (*model.Manifest, error) {

	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "configuration file [%s] could not be read", path)
	}

	manifest := &model.Manifest{}
	err = yaml.Unmarshal(configBytes, manifest)
	if err != nil {
		return nil, errors.Wrapf(err, "configuration file [%s] could not be parsed", path)
	}

	return manifest, nil
}

func ReadManifest() (*model.Manifest, error) {

	return ReadManifestFrom(DefaultManifestFile)
}
