package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func ToInt(v interface{}) (res int) {
	switch s := v.(type) {
	case int:
		res = s
	case string:
		num, err := strconv.ParseFloat(s, 64)
		if err == nil {
			res = int(num)
		} else {
			num2, err2 := strconv.ParseInt(s, 10, 64)
			if err2 == nil {
				res = int(num2)
			}
		}
	case float64:
		res = int(s)
	default:
		res = 0
	}
	return res
}

func ToInt64(v interface{}) (res int64) {
	switch s := v.(type) {
	case int64:
		res = s
	case int:
		res = v.(int64)
	case string:
		num, err := strconv.ParseFloat(s, 64)
		if err == nil {
			res = int64(num)
		} else {
			num2, err2 := strconv.ParseInt(s, 10, 64)
			if err2 == nil {
				res = int64(num2)
			}
		}
	case float64:
		res = int64(s)
	default:
		res = 0
	}
	return res
}

func ToFloat(v interface{}) (res float64) {
	switch s := v.(type) {
	case int64:
		res = v.(float64)
	case int:
		res = v.(float64)
	case string:
		num, err := strconv.ParseFloat(s, 64)
		if err == nil {
			res = num
		}
	case float64:
		res = s
	default:
		res = 0
	}
	return res
}

func ToBool(v interface{}) (res bool) {
	switch s := v.(type) {
	case int64:
		res = v.(int64) > 0
	case int:
		res = v.(int) > 0
	case string:
		res, _ = strconv.ParseBool(v.(string))
	case float64:
		res = v.(float64) > 0
	case bool:
		res = s
	default:
		res = false
	}
	return res
}

func ToString(i interface{}) (res string) {
	switch s := i.(type) {
	case string:
		res = s
	case bool:
		res = strconv.FormatBool(s)
	case float64:
		res = strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		res = strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		res = strconv.Itoa(s)
	case int64:
		res = strconv.FormatInt(s, 10)
	case int32:
		res = strconv.Itoa(int(s))
	case int16:
		res = strconv.FormatInt(int64(s), 10)
	case int8:
		res = strconv.FormatInt(int64(s), 10)
	case uint:
		res = strconv.FormatUint(uint64(s), 10)
	case uint64:
		res = strconv.FormatUint(uint64(s), 10)
	case uint32:
		res = strconv.FormatUint(uint64(s), 10)
	case uint16:
		res = strconv.FormatUint(uint64(s), 10)
	case uint8:
		res = strconv.FormatUint(uint64(s), 10)
	case []byte:
		res = string(s)
	case nil:
		res = ""
	case fmt.Stringer:
		res = s.String()
	case error:
		res = s.Error()
	default:
		res = ""
	}
	return res
}

func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

func ToStringMapE(i interface{}) (map[string]interface{}, error) {
	var m = map[string]interface{}{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToString(k)] = val
		}
		return m, nil
	case map[string]interface{}:
		return v, nil
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]interface{}", i, i)
	}
}

func ToStringSlice(i interface{}) (res []string) {
	switch v := i.(type) {
	case []string:
		res = v
	case string:
		res = strings.Fields(v)
	case []error:
		for _, err := range i.([]error) {
			res = append(res, err.Error())
		}
	case []interface{}:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	case []int8:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	case []int:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	case []int32:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	case []int64:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	case []float32:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	case []float64:
		for _, u := range v {
			res = append(res, ToString(u))
		}
	default:
		res = []string{}
	}
	return res
}

func Empty(param interface{}) bool {
	return reflect.ValueOf(param).IsZero()
}

func FileExists(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}

func jsonStringToObject(s string, v interface{}) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}
