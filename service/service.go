package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"ys-keyvalue-store/errors"
	"ys-keyvalue-store/repository"
)

type Service interface {
	GetKeyValue(string) (string, error)
	SetKeyValue(string, string) error
	SaveKeyValuesToFile()
	LoadKeyValueStoreToMemory()
}

type service struct {
	keyValueRepository repository.KeyValueRepository
}

func NewService(keyValueRepository repository.KeyValueRepository) Service {
	return &service{
		keyValueRepository: keyValueRepository,
	}
}

// GetKeyValue gets key and values from store.
func (s *service) GetKeyValue(key string) (string, error) {
	keyValueData, _ := s.keyValueRepository.GetKey(key)
	if keyValueData == "" {
		return "", errors.ErrorValueNotFound
	}
	return keyValueData, nil
}

// SetKeyValue sets key and values to the store.
// First, it checks the store for the same key-value.
// Otherwise, it saves the new key-value to store.
func (s *service) SetKeyValue(key, value string) error {
	keyValue, _ := s.GetKeyValue(key)
	if keyValue != "" {
		return errors.ErrorKeyValueAlreadyExist
	}
	err := s.keyValueRepository.SetKey(key, value)
	if err != nil {
		return err
	}
	return nil
}

// SaveKeyValuesToFile writes the last data to a new file every 1 minute.
func (s *service) SaveKeyValuesToFile() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		err := s.keyValueRepository.SaveKeyValuesToFile()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

// LoadKeyValueStoreToMemory saves the entered key-values to the memory.
// First, it checks if there is a folder named "tmp". Then it looks at
// every file in this folder and saves its contents to memory.
func (s *service) LoadKeyValueStoreToMemory() {
	dir := "tmp/"
	var lastTime int64
	var lastFile string
	var filePath strings.Builder
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
		}
	}

	files, _ := ioutil.ReadDir("tmp/") //nolint
	if len(files) > 0 {
		for _, file := range files {
			fileName, err := os.Stat("tmp/" + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			currentTime := fileName.ModTime().Unix()
			if currentTime > lastTime {
				lastTime = currentTime
				lastFile = file.Name()
			}
		}
		filePath.WriteString("tmp/")
		filePath.WriteString(lastFile)
		err := s.keyValueRepository.LoadKeyValueStoreToMemory(filePath.String())
		if err != nil {
			log.Println(err.Error())
		}
	}
}
