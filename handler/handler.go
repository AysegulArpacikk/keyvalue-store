package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"ys-keyvalue-store/errors"
	"ys-keyvalue-store/service"
)

type KeyValueHandler interface {
	GetKeyValueHandler(http.ResponseWriter, *http.Request) error
	SetKeyValueHandler(http.ResponseWriter, *http.Request)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type handler struct {
	keyValueService service.Service
}

func NewHandler(keyValueService service.Service) KeyValueHandler {
	return &handler{
		keyValueService: keyValueService,
	}
}

// GetKeyValueHandler gets the key/value pair from in memory.
// The status was defined as 200 by default.
// Then it goes to the service layer and returns the response.
// If no error is returned, it returns the desired key-value and StatusOK.
func (h *handler) GetKeyValueHandler(w http.ResponseWriter, r *http.Request) error {
	key := strings.TrimPrefix(r.URL.Path, "/api/key/")
	status := http.StatusOK

	var response map[string]string

	value, getKeyValueErr := h.keyValueService.GetKeyValue(key)
	if getKeyValueErr != nil {
		if getKeyValueErr != errors.ErrorValueNotFound {
			status = http.StatusInternalServerError //nolint
			response = map[string]string{
				"": "",
			}
		} else {
			getKeyValueErr = errors.ErrorValueNotFound
			status = http.StatusNotFound
			response = map[string]string{
				"": "",
			}
			fmt.Fprintf(w, "Error = %s, Status = %d", getKeyValueErr, status)
		}
	} else {
		response = map[string]string{
			key: value,
		}
		//fmt.Fprintf(w, "Key= %q, Value = %q, Status = %d", key, value, status)
	}

	return json.NewEncoder(w).Encode(response)
}

// SetKeyValueHandler creates a new key-value pair.
// Key-values are taken and status is defined as StatusOK by default.
// Error codes are printed on the screen according to the response from the service layer.
func (h *handler) SetKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //nolint

	keyValue := map[string]string{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	unmarshalErr := json.Unmarshal(reqBody, &keyValue)
	if unmarshalErr != nil {
		return
	}

	status := http.StatusOK
	setErr := h.keyValueService.SetKeyValue(&keyValue)
	if setErr != nil {
		if setErr != errors.ErrorKeyValueAlreadyExist {
			status = http.StatusInternalServerError
		} else {
			setErr = errors.ErrorKeyValueAlreadyExist
			status = http.StatusConflict
			fmt.Fprintf(w, "Error = %s", setErr)
		}
	}

	fmt.Fprintf(w, "Status = %v", status)
	w.WriteHeader(status)
}

// ServeHTTP was created for tests and http request logger. Handler satisfy
// http.Handler, so we can call their ServeHTTP method directly and pass in our
// Request and ResponseRecorder.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.URL.Path == "/api/keyValues" && r.Method == http.MethodPost {
		h.SetKeyValueHandler(w, r)
	} else if regexp.MustCompile(`/api/key/[a-zA-Z]+`).MatchString(r.URL.Path) && r.Method == http.MethodGet {
		h.GetKeyValueHandler(w, r)
	} else {
		http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /key/, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
