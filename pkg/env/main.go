package env

import (
	"os"
	"strconv"

	"github.com/poligonoio/vega-core/pkg/logger"
)

func GetBoolEnv(key string) bool {
	value, exists := os.LookupEnv(key)
	boolenv := false // Default value

	if exists {
		boolenv, err := strconv.ParseBool(value)
		if err != nil {
			logger.Error.Panicf("Error reading environment variable '%s': %v", key, err)
			return boolenv
		}

		return boolenv
	}

	return boolenv
}
