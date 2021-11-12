package service

import (
	"io"
	"log"
	"net/http"
	"testing"
	"ys-keyvalue-store/repository"

	"github.com/stretchr/testify/assert"
)

const testKey = "testKey"
const testValue = "testValue"

func Test_ShouldServiceSetNewKeyValue(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{})
	keyValueService := NewService(keyValueRepository)

	_, err := http.NewRequest("POST", "/api/keyValues", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	setKeyError := keyValueService.SetKeyValue(testKey, testValue)
	if setKeyError != nil {
		log.Fatalf("Error: '%s'", setKeyError)
	}

	assert.Nil(t, setKeyError)
}

func Test_ShouldServiceNotSetNewKeyValue(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	keyValueService := NewService(keyValueRepository)

	_, err := http.NewRequest("POST", "/api/keyValues", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	setKeyError := keyValueService.SetKeyValue(testKey, testValue)

	assert.NotNil(t, setKeyError)
	assert.EqualError(t, setKeyError, "key/value already exist ")
}

func Test_ShouldServiceGetKeyValue(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	keyValueService := NewService(keyValueRepository)

	_, err := http.NewRequest("GET", "/api/key/", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	value, getKeyError := keyValueService.GetKeyValue(testKey)
	if getKeyError != nil {
		log.Fatalln("Error")
	}

	assert.Equal(t, value, testValue)
	assert.Nil(t, getKeyError)
}

func Test_ShouldServiceNotGetKeyValue(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{})
	keyValueService := NewService(keyValueRepository)

	_, err := http.NewRequest("GET", "/api/key/", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	_, getKeyError := keyValueService.GetKeyValue(testKey)

	assert.NotNil(t, getKeyError)
	assert.EqualError(t, getKeyError, "value not found ")
}

func Test_InitializeKeyStoreData(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	_ = keyValueRepository.SaveKeyValuesToFile()
	keyValueRepository = repository.NewKeyValueStoreRepository(map[string]string{})
	keyValueService := NewService(keyValueRepository)
	keyValueService.LoadKeyValueStoreToMemory()

	_, err := keyValueRepository.GetKey(testKey)
	if err != nil {
		log.Fatalf("An errors occurred when taking data from file %s", err)
	}
}
