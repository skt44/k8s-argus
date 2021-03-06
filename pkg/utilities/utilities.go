package utilities

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/logicmonitor/k8s-argus/pkg/metrics"

	"github.com/logicmonitor/lm-sdk-go"
)

// BuildSystemCategoriesFromLabels formats a system.categories string.
func BuildSystemCategoriesFromLabels(categories string, labels map[string]string) string {
	for k, v := range labels {
		categories += "," + k + "=" + v
	}
	return categories
}

// GetLabelByPrefix takes a list of labels returns the first label matching the specified prefix
func GetLabelByPrefix(prefix string, labels map[string]string) (string, string) {
	for k, v := range labels {
		if match, err := regexp.MatchString("^"+prefix, k); match {
			if err != nil {
				continue
			}
			return k, v
		}
	}
	return "", ""
}

// CheckAllErrors is a helper function to deal with the number of possible places that an API call can fail.
func CheckAllErrors(restResponse interface{}, apiResponse *logicmonitor.APIResponse, err error) error {
	var restResponseMessage string
	var restResponseStatus int64

	// Get the underlying concrete type.
	t := reflect.ValueOf(restResponse)

	// Check it the interface is a pointer and get the underlying value.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Ensure that it is a struct, and get the necessary fields if they are available.
	if t.Kind() == reflect.Struct {
		field := t.FieldByName("Status")
		if field.IsValid() {
			restResponseStatus = field.Int()
		}
		field = t.FieldByName("Errmsg")
		if field.IsValid() {
			restResponseMessage = field.String()
		}
	}

	if restResponseStatus != http.StatusOK {
		metrics.RESTError()
		return fmt.Errorf("[REST] [%d] %s", restResponseStatus, restResponseMessage)
	}

	if apiResponse.StatusCode != http.StatusOK {
		metrics.APIError()
		return fmt.Errorf("[API] [%d] %s", apiResponse.StatusCode, restResponseMessage)
	}

	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}

	return nil
}
