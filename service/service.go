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
	GetKey(string) (string, error)
	SetKey(string, string) error
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

func (s *service) GetKey(key string) (string, error) {
	keyValueData, _ := s.keyValueRepository.GetKey(key)
	if keyValueData == "" {
		return "", errors.ErrorValueNotFound
	}
	return keyValueData, nil
}

func (s *service) SetKey(key, value string) error {
	keyValue, _ := s.GetKey(key)
	if keyValue != "" {
		return errors.ErrorKeyValueAlreadyExist
	}
	err := s.keyValueRepository.SetKey(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SaveKeyValuesToFile() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		err := s.keyValueRepository.SaveKeyValuesToFile()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

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

	files, _ := ioutil.ReadDir("tmp/")
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
