# In Memory Key Value Store Project with Golang

## Project Details

This project keeps key-values saved in memory. It keeps these key-values in files named 
TIMESTAMP-data.json and rewrites it to the file every one minute. When the application is 
stopped and reopened, it writes the data in the last saved file back to the memory.

## API docs

| Method | Route | Request body | Description
| --- | --- | --- | --- |
| /api/key/:key | `GET` |     |  Gets the desired key-value from memory. |
| /api/keyValues | `POST` | {"key1":"value1", "key2":"value2"} | Creates new key-value in memory. |

## Installation
If golang is not installed on your computer, first install go from this link: `https://golang.org/doc/install`

Then install the dependencies
```bash
go mod download
```

## Usage 

Run the project
```bash
make run
```

### For create a new key-value to in memory:
Use this command on terminal: `curl -X POST 'http://localhost:8080/api/keyValues?key=keyData&value=valueData'`
If you want to use the browser or Postman, you can go to the following address: `http://localhost:8080/api/keyValues?key=keyData&value=valueData`

### For get a specific data from memory:
Use this command on terminal: `curl 'http://localhost:8080/api/key/keyData'`
If you want to use the browser or Postman, you can go to the following address: `http://localhost:8080/api/key/keyData`

Run tests
```bash
make test
```

Build with docker
```bash
sudo docker build -t ys-keyvalue-store .
```

For lint run this file on project root directory:
`lint.sh`

## Layout

The project tree looks like this: 

```
├── errors
│ └── errors.go
├── go.mod
├── go.sum
├── handler
│ ├── handler.go
│ └── handler_test.go
├── httpRequests.log
├── lint.sh
├── logger
│ └── request_logger.go
├── main.go
├── Makefile
├── README.md
├── repository
│ ├── repository.go
│ ├── repository_test.go
│ └── tmp
│     └── 1636616810-data.json
├── service
│ ├── service.go
│ ├── service_test.go
│ └── tmp
│     └── 1636616810-data.json
└── tmp
    ├── 1636543351-data.json
    └── 1636543411-data.json
```
