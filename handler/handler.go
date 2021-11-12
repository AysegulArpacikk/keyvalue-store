package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"ys-keyvalue-store/errors"
	"ys-keyvalue-store/service"
)

type KeyValueHandler interface {
	GetKeyValueHandler(http.ResponseWriter, *http.Request)
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
func (h *handler) GetKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/key/")
	status := http.StatusOK
	value, getKeyValueErr := h.keyValueService.GetKeyValue(key)
	if getKeyValueErr != nil {
		if getKeyValueErr != errors.ErrorValueNotFound {
			status = http.StatusInternalServerError //nolint
		} else {
			getKeyValueErr = errors.ErrorValueNotFound
			status = http.StatusNotFound
			fmt.Fprintf(w, "Error = %s, Status = %d", getKeyValueErr, status)
		}
	} else {
		fmt.Fprintf(w, "Key= %q, Value = %q, Status = %d", key, value, status)
	}
}

// SetKeyValueHandler creates a new key-value pair.
// Key-values are taken and status is defined as StatusOK by default.
// Error codes are printed on the screen according to the response from the service layer.
func (h *handler) SetKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //nolint

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	status := http.StatusOK
	setErr := h.keyValueService.SetKeyValue(key, value)
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
