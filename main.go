package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"ys-keyvalue-store/handler"
	"ys-keyvalue-store/logger"
	"ys-keyvalue-store/repository"
	"ys-keyvalue-store/service"
)

var (
	httpAddr = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	perm     = 0666
)

func main() {
	fileName := "httpRequests.log"

	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(perm))
	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)
	defer logFile.Close()

	keyValueRepository := repository.NewKeyValueStoreRepository(map[string]string{})
	keyValueService := service.NewService(keyValueRepository)
	keyValueHandler := handler.NewHandler(keyValueService)
	keyValueService.LoadKeyValueStoreToMemory()
	go keyValueService.SaveKeyValuesToFile()

	mux := http.NewServeMux()
	mux.Handle("/api/key/", keyValueHandler)
	mux.Handle("/api/keyValues", keyValueHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, logger.RequestLogger(mux))) //nolint
}
