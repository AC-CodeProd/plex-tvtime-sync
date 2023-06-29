package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/fx"
)

type Link map[int]int

type JsonStorage struct {
	config Config
	logger Logger
}

// Params defines the base objects for a storage.
type JsonStorageParams struct {
	fx.In
	Config Config
	Logger Logger
}

// Result defines the objects that the storage module provides.
type JsonStorageResult struct {
	fx.Out

	JsonStorage JsonStorage
}

func GetJsonStorage(jsP JsonStorageParams) JsonStorageResult {
	jsonStorage := JsonStorage{
		config: jsP.Config,
		logger: jsP.Logger,
	}
	return JsonStorageResult{
		JsonStorage: jsonStorage,
	}
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (jS *JsonStorage) GetLinks() (Link, error) {
	const names = "__storage.go__ : GetLinks"
	jS.logger.Info(names)
	links := make(Link)
	err := ensureDir(jS.config.FileStorage.Filename)
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	file, err := ioutil.ReadFile(jS.config.FileStorage.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return links, nil
		}

		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	err = json.Unmarshal(file, &links)
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	return links, nil
}

func (jS *JsonStorage) AddLink(idShowPlex int, idShowTvTime int) error {
	const names = "__storage.go__ : AddLink"
	jS.logger.Info(fmt.Sprintf("%s | %d - %d", names, idShowPlex, idShowTvTime))
	links, err := jS.GetLinks()
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return err
	}

	links[idShowPlex] = idShowTvTime

	file, err := json.Marshal(links)
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return err
	}

	err = ioutil.WriteFile(jS.config.FileStorage.Filename, file, 0644)
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return err
	}

	return nil
}

func (jS *JsonStorage) HasLink(idShowPlex int) (bool, error) {
	const names = "__storage.go__ : HasLink"
	jS.logger.Info(fmt.Sprintf("%s | %d", names, idShowPlex))
	links, err := jS.GetLinks()
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	_, exists := links[idShowPlex]

	return exists, nil
}

func (jS *JsonStorage) GetLink(idShowPlex int) (int, bool, error) {
	const names = "__storage.go__ : GetLink"
	jS.logger.Info(fmt.Sprintf("%s | %d", names, idShowPlex))
	links, err := jS.GetLinks()
	if err != nil {
		jS.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return 0, false, err
	}

	link, exists := links[idShowPlex]

	return link, exists, nil
}
