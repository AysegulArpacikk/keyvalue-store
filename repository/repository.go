package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type KeyValueRepository interface {
	GetKey(string) (string, error)
	SetKey(*map[string]string) error
	SaveKeyValuesToFile() error
	LoadKeyValueStoreToMemory(string) error
}

// KeyValueStore structure was created with map.
type KeyValueStore struct {
	store []*map[string]string
}

func NewKeyValueStoreRepository(keyValue []*map[string]string) *KeyValueStore {
	store := &KeyValueStore{store: keyValue}
	return store
}

// GetKey gets value by key from in memory
func (k *KeyValueStore) GetKey(key string) (string, error) {
	for key1, _ := range k.store {
		for key2, value := range *k.store[key1] {
			if key2 == key {
				return value, nil
			}
		}
	}
	return "", errors.New("key-value not found")
}

// SetKey sets key and value to in memory
func (k *KeyValueStore) SetKey(keyValue *map[string]string) error {
	k.store = k.addKeyValueData(keyValue)
	return nil
}

func (k *KeyValueStore) addKeyValueData(keyValue *map[string]string) []*map[string]string {
	return append(k.store, keyValue)
}

// LoadKeyValueStoreToMemory loads the key-values saved in the file into memory.
func (k *KeyValueStore) LoadKeyValueStoreToMemory(file string) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("Open the failed: '%s'\n", err)
	}
	fmt.Printf("Successfully Opened '%s'", file)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile) //nolint

	err = json.Unmarshal(byteValue, &k.store)
	if err != nil {
		return err
	}
	fmt.Println(k.store)
	return nil
}

// SaveKeyValuesToFile saves the entered key-values to the file.
// Under the tmp directory, it converts the recorded time to timestamp
// and creates a json file with that name. Then the key-values saved
// in the memory are encoded, converted to json format and written to the file.
func (k *KeyValueStore) SaveKeyValuesToFile() error {
	var stringBuilder strings.Builder
	var base = 10
	if _, err := os.Stat("tmp/"); os.IsNotExist(err) {
		err := os.Mkdir("tmp/", os.ModePerm)
		if err != nil {
			log.Println(err.Error())
		}
	}
	stringBuilder.WriteString("tmp/")
	stringBuilder.WriteString(strconv.FormatInt(time.Now().Unix(), base))
	stringBuilder.WriteString("-data.json")
	fileName, err := os.Create(stringBuilder.String())
	if err != nil {
		return err
	}

	log.Printf("'%s' file created", fileName.Name())
	defer fileName.Close()

	dataBytes, err1 := json.Marshal(k.store)
	if err1 != nil {
		return err1
	}
	_, writeToFileError := fileName.Write(dataBytes)
	if writeToFileError != nil {
		log.Printf("An arror was occure when write key/value to json file")
		return writeToFileError
	}

	log.Printf("'%s' data writed to file '%s' file", string(dataBytes), fileName.Name())

	return nil
}
