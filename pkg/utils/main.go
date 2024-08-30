package utils

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/poligonoio/vega-core/pkg/logger"
)

// isZeroType checks if the value from the struct is the zero value of its type
func IsZeroType(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()

	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return value.Len() == 0
	default:
		return reflect.DeepEqual(zero, value.Interface())
	}
}

func MapToStruct(mapObject map[string]interface{}, structOBject *interface{}) error {
	jsonStr, err := json.Marshal(mapObject)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonStr, structOBject); err != nil {
		return err
	}

	return nil
}

func SanitizeQuery(query string) string {
	q := strings.ReplaceAll(query, "`", "")
	q = strings.ReplaceAll(q, "sql", "")
	q = strings.ReplaceAll(q, ";", "")
	q = strings.ReplaceAll(q, "\"", "")

	return q
}
