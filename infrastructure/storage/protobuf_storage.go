package storage

import (
	"os"

	"plex-tvtime-sync/dto"

	"plex-tvtime-sync/domain/interfaces"

	"google.golang.org/protobuf/proto"
)

func NewProtobufStorage(psp ProtobufStorageParams) interfaces.IStorage {
	return &protobufStorage{
		logger: psp.Logger,
		config: psp.Config,
	}
}

func (ps *protobufStorage) Save(data *dto.IntMap) error {
	bytes, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(ps.config.Storage.Filename, bytes, 0644); err != nil {
		return err
	}

	ps.mu.Lock()
	ps.loadedData = data
	ps.mu.Unlock()

	return nil
}

func (ps *protobufStorage) Load() (*dto.IntMap, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.loadedData != nil {
		return ps.loadedData, nil
	}

	bytes, err := os.ReadFile(ps.config.Storage.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			ps.loadedData = &dto.IntMap{}
			return ps.loadedData, nil
		}
		return nil, err
	}

	var data dto.IntMap
	if err := proto.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	ps.loadedData = &data
	return ps.loadedData, nil
}

func (ps *protobufStorage) AddSpecificPair(key, value *int64) error {
	intMap, err := ps.Load()
	if err != nil {
		return err
	}

	if intMap.Map == nil {
		intMap.Map = make(map[int64]int64)
	}
	intMap.Map[*key] = *value

	if err := ps.Save(intMap); err != nil {
		return err
	}

	return nil
}

func (ps *protobufStorage) GetValue(key int64) (*int64, error) {
	intMap, err := ps.Load()
	if err != nil {
		return nil, err
	}

	value, exists := intMap.Map[key]
	if !exists {
		return nil, nil
	}

	return &value, nil
}
