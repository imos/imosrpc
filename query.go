package imosrpc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func getJsonKeys(reflectType reflect.Type) map[string][]int {
	fields := map[string][]int{}
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	if reflectType.Kind() != reflect.Struct {
		return fields
	}
	for fieldIndex := 0; fieldIndex < reflectType.NumField(); fieldIndex++ {
		field := reflectType.Field(fieldIndex)
		if field.Anonymous == true {
			for key, childIndex := range getJsonKeys(field.Type) {
				fields[key] = append([]int{fieldIndex}, childIndex...)
			}
			continue
		}
		jsonTag := field.Tag.Get("json")
		key := field.Name
		if jsonTag == "-" {
			continue
		}
		jsonTags := strings.Split(jsonTag, ",")
		if jsonTags[0] != "" {
			key = jsonTags[0]
		}
		fields[key] = []int{fieldIndex}
	}
	return fields
}

func buildQuery(input interface{}) (querySet map[string]string, err error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		err = fmt.Errorf("failed to build query: %s", err)
		return
	}
	rawQuerySet := map[string]interface{}{}
	err = json.Unmarshal(jsonData, &rawQuerySet)
	if err != nil {
		err = fmt.Errorf("failed to build query: %s", err)
		return
	}

	querySet = map[string]string{}
	for key, value := range rawQuerySet {
		// fmt.Sprintf may cause panic, so this catches it and converts it to error.
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("failed to parse: %s=%s", key, value)
			}
		}()
		if value != nil {
			querySet[key] = fmt.Sprintf("%v", value)
		}
	}
	return
}

func parseQuery(output interface{}, input map[string]string) error {
	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return fmt.Errorf(
			"output must be a pointer, but %s.",
			reflect.ValueOf(output).Kind().String())
	}
	outputValue := reflect.ValueOf(output).Elem()
	keyToFieldIndex := getJsonKeys(reflect.ValueOf(output).Elem().Type())
	for key, value := range input {
		if fieldIndex, ok := keyToFieldIndex[key]; ok {
			for sliceSize := 1; sliceSize < len(fieldIndex); sliceSize++ {
				subField := outputValue.FieldByIndex(fieldIndex[:sliceSize])
				if subField.Kind() == reflect.Ptr && subField.IsNil() {
					subField.Set(reflect.New(subField.Type().Elem()))
				}
			}
			field := outputValue.FieldByIndex(fieldIndex)
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				field = field.Elem()
			}
			switch field.Interface().(type) {
			case bool, int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64, float32, float64:
				err := json.Unmarshal([]byte(value), field.Addr().Interface())
				if err != nil {
					return fmt.Errorf("failed to parse: %s=%s", key, value)
				}
			default:
				err := json.Unmarshal([]byte("\""+value+"\""), field.Addr().Interface())
				if err != nil {
					return fmt.Errorf("failed to parse: %s=%s", key, value)
				}
			}
		} else {
			return fmt.Errorf("unknown key: %s", key)
		}
	}
	return nil
}

func BuildForm(request interface{}) (form url.Values, err error) {
	querySet, err := buildQuery(request)
	if err != nil {
		return
	}
	form = url.Values{}
	for key, value := range querySet {
		form[key] = []string{value}
	}
	return
}

func ParseForm(responsePtr interface{}, form url.Values) error {
	input := map[string]string{}
	for key, values := range form {
		if len(values) != 1 {
			return fmt.Errorf("key '%s' has multiple values.", key)
		}
		input[key] = values[0]
	}
	return parseQuery(responsePtr, input)
}
