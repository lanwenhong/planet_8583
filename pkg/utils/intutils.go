package utils

import (
	"fmt"
	"strconv"
)

func ToInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case []byte:
		return strconv.Atoi(string(v))
	default:
		return 0, fmt.Errorf("unsupported type for conversion to int")
	}
}

func SafeToInt(i interface{}) int {
	v, _ := ToInt(i)
	return v
}

func ToInt64(i interface{}) (int64, error) {
	switch v := i.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported type for conversion to int64")
	}
}

func SafeToInt64(i interface{}) int64 {
	v, _ := ToInt64(i)
	return v
}

func SafeToBool(i interface{}) bool {
	switch v := i.(type) {
	case bool:
		return v
	case string:
		return v == "1" || v == "true" || v == "True" || v == "TRUE"
	case []byte:
		s := string(v)
		return s == "1" || s == "true" || s == "True" || s == "TRUE"
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0
	default:
		return false
	}
}
