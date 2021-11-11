package repository

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKey = "testKey"
const testValue = "testValue"

func Test_ShouldServiceSetNewKeyValueToMemory(t *testing.T) {
	keyValueRepository := NewKeyValueStoreRepository(map[string]string{})
	setKeyError := keyValueRepository.SetKey(testKey, testValue)
	if setKeyError != nil {
		log.Fatalf("An error occurred when set key/value to memory: '%s'", setKeyError)
	}

	assert.Nil(t, setKeyError)
}

func Test_ShouldServiceGetKeyValueToMemory(t *testing.T) {
	keyValueRepository := NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	value, getKeyError := keyValueRepository.GetKey(testKey)
	if getKeyError != nil {
		log.Fatalf("An error occurred when get key/value from memory: '%s'", getKeyError)
	}

	assert.Equal(t, value, testValue)
	assert.Nil(t, getKeyError)
}

func Test_SaveKeyValuesToFile(t *testing.T) {
	keyValueRepository := NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	err := keyValueRepository.SaveKeyValuesToFile()
	if err != nil {
		log.Fatalf("An error occurred when writing keystore data to file: '%s'", err.Error())
	}
}
