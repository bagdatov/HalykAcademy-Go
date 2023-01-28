package flatjson

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	errorNil        = errors.New("nil pointer as argument")
	errorNotPointer = errors.New("received argument is not pointer to struct")
)

func Marshal(src interface{}) ([]byte, error) {
	if src == nil {
		return nil, errorNil
	}
	// TODO: Write code here
	return nil, nil
}

// Unmarshal is ensuring destination is a pointer to struct.
// Decodes data to map[string]interface{}.
// Recursivly prepares new map using nested values from decoded map.
// Marshals this new map and unmarshals it to struct.
// Very slow and very costly.
func Unmarshal(data []byte, dst interface{}) error {
	// validation
	if data == nil || dst == nil {
		return errorNil
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr {
		return errorNotPointer
	}

	vl := reflect.Indirect(rv)
	if vl.Kind() != reflect.Struct {
		return errorNotPointer
	}

	// decodes data to map
	decodedMap := make(map[string]interface{})
	err := json.Unmarshal(data, &decodedMap)
	if err != nil {
		return err
	}

	// create new map
	// this is simply required to convert json with nested data to flat json
	newMap := make(map[string]interface{})
	// recursivly create data with key from previous level of nesting + current key
	recursiveMapping(decodedMap, newMap, "")

	m, err := json.Marshal(newMap)
	if err != nil {
		return err
	}
	// and then unmarshal flat json to flat struct
	return json.Unmarshal(m, dst)
}

func recursiveMapping(decodedMap, newMap map[string]interface{}, prevK string) {
	for k, v := range decodedMap {
		key := prevK + k
		switch v.(type) {
		case []map[string]interface{}:
			continue
		case map[string]interface{}:
			recursiveMapping(v.(map[string]interface{}), newMap, key)
		default:
			newMap[key] = v
		}
	}
}

// UnmarshalAlternative is variant with use of reflect values.
func UnmarshalAlternative(data []byte, dst interface{}) error {
	// validation
	if data == nil || dst == nil {
		return errorNil
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr {
		return errorNotPointer
	}

	vl := reflect.Indirect(rv)
	if vl.Kind() != reflect.Struct {
		return errorNotPointer
	}

	// decode with unmarshal to map
	var res map[string]interface{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	// recursivly check if our struct has field with name matching from decoded map
	// and add data to this field
	recursiveUnmarshal(res, &vl, "")

	return nil
}

func recursiveUnmarshal(res map[string]interface{}, vl *reflect.Value, lastKey string) {
	for k, v := range res {
		// check before each insertion field existence in our struct
		switch v.(type) {
		case string:
			field := vl.FieldByName(lastKey + k)
			if !field.IsValid() || field.Kind() == reflect.Ptr {
				continue
			}
			field.SetString(v.(string))

		case bool:
			field := vl.FieldByName(lastKey + k)
			if !field.IsValid() || field.Kind() == reflect.Ptr {
				continue
			}
			field.SetBool(v.(bool))

		// note: json.Unmarshal converted all numeric types from map to float64
		case float64:
			field := vl.FieldByName(lastKey + k)
			if !field.IsValid() || field.Kind() == reflect.Ptr {
				continue
			}
			fieldKind := field.Kind()

			if fieldKind == reflect.Float32 || fieldKind == reflect.Float64 {
				field.SetFloat(v.(float64))

			} else if fieldKind >= reflect.Int && fieldKind <= reflect.Int64 {
				field.SetInt(int64(v.(float64)))

			} else if fieldKind >= reflect.Uint && fieldKind <= reflect.Uint64 {
				field.SetUint(uint64(v.(float64)))
			}

		// go deeper to nested value:
		case map[string]interface{}:
			recursiveUnmarshal(v.(map[string]interface{}), vl, lastKey+k)
		}
	}
}
