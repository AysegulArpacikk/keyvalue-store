package handler_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"ys-keyvalue-store/handler"
	"ys-keyvalue-store/repository"
	"ys-keyvalue-store/service"

	"github.com/stretchr/testify/assert"
)

const testKey = "testKey"
const testValue = "testValue"

func Test_ShouldSetNewKeyValue(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{})
	keyValueService := service.NewService(keyValueRepository)
	keyValueHandler := handler.NewHandler(keyValueService)

	server := httptest.NewServer(keyValueHandler)

	request, err := http.NewRequest("POST", "/api/keyValues", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	setKeyError := keyValueService.SetKey(testKey, testValue)
	if setKeyError != nil {
		log.Fatalf("Error: '%s'", setKeyError)
	}

	handle := http.HandlerFunc(keyValueHandler.SetKeyValueHandler)
	handle.ServeHTTP(recorder, request)
	server.Close()

	assert.Nil(t, setKeyError)
}

func Test_ShouldNotSetNewKeyValueBecauseKeyValueAlreadyExist(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	keyValueService := service.NewService(keyValueRepository)
	keyValueHandler := handler.NewHandler(keyValueService)

	server := httptest.NewServer(keyValueHandler)
	defer server.Close()

	request, err := http.NewRequest("POST", "/api/keyValues", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	setKeyError := keyValueService.SetKey(testKey, testValue)

	handle := http.HandlerFunc(keyValueHandler.SetKeyValueHandler)
	handle.ServeHTTP(recorder, request)

	assert.NotNil(t, setKeyError)
	assert.EqualError(t, setKeyError, "key/value already exist ")
}

func Test_ShouldGetKeyValue(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{testKey: testValue})
	keyValueService := service.NewService(keyValueRepository)
	keyValueHandler := handler.NewHandler(keyValueService)

	server := httptest.NewServer(keyValueHandler)

	request, err := http.NewRequest("GET", "/api/key/", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	value, getKeyError := keyValueService.GetKey(testKey)
	if getKeyError != nil {
		log.Fatalln("Error")
	}

	handle := http.HandlerFunc(keyValueHandler.GetKeyValueHandler)
	handle.ServeHTTP(recorder, request)
	server.Close()

	assert.Equal(t, value, testValue)
	assert.Nil(t, getKeyError)
}

func Test_ShouldNotGetKeyValueBecauseKeyValueNotFound(t *testing.T) {
	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{})
	keyValueService := service.NewService(keyValueRepository)
	keyValueHandler := handler.NewHandler(keyValueService)

	server := httptest.NewServer(keyValueHandler)
	defer server.Close()

	request, err := http.NewRequest("GET", "/api/key/", io.Reader(nil))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	_, getKeyError := keyValueService.GetKey(testKey)

	handle := http.HandlerFunc(keyValueHandler.GetKeyValueHandler)
	handle.ServeHTTP(recorder, request)

	assert.NotNil(t, getKeyError)
	assert.EqualError(t, getKeyError, "value not found ")
}
