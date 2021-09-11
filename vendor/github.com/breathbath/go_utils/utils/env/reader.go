package env

import (
	"fmt"
	"github.com/breathbath/go_utils/utils/errs"
	"os"
	"strconv"
)

func ReadEnv(name, defaultVal string) string {
	return readEnvironment(
		name,
		defaultVal,
		func(val string) interface{} {
			return val
		},
	).(string)
}

func readEnvironment(name string, defaultVal interface{}, conv func(val string) interface{}) interface{} {
	envVar, found := os.LookupEnv(name)
	if !found {
		return defaultVal
	}
	return conv(envVar)
}

func ReadEnvInt64(name string, defaultVal int64) int64 {
	return readEnvironment(
		name,
		defaultVal,
		func(val string) interface{} {
			intRes, err := strconv.ParseInt(val, 0, 64)
			if err != nil {
				return defaultVal
			}
			return intRes
		},
	).(int64)
}

func ReadEnvInt(name string, defaultVal int) int {
	return readEnvironment(
		name,
		defaultVal,
		func(val string) interface{} {
			intRes, err := strconv.Atoi(val)
			if err != nil {
				return defaultVal
			}
			return intRes
		},
	).(int)
}

func ReadEnvBool(name string, defaultVal bool) bool {
	return readEnvironment(
		name,
		defaultVal,
		func(val string) interface{} {
			if val == "true" {
				return true
			}
			return false
		},
	).(bool)
}

func ReadEnvFloat(name string, defaultVal float64) float64 {
	return readEnvironment(
		name,
		defaultVal,
		func(val string) interface{} {
			floatResult, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return defaultVal
			}
			return floatResult
		},
	).(float64)
}

func ReadEnvOrError(name string) (string, error) {
	val := ReadEnv(name, "")
	if val == "" {
		return "", fmt.Errorf("Required env variable '%s' is not set", name)
	}

	return val, nil
}

func ReadEnvOrFail(name string) string {
	val, err := ReadEnvOrError(name)
	errs.FailOnError(err)
	return val
}
