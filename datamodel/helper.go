package datamodel

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dchaykin/go-modules/log"
)

func LoadAccessData(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	accessData := map[string]string{}
	if err = json.Unmarshal(data, &accessData); err != nil {
		panic(err)
	}

	for k, v := range accessData {
		os.Setenv(k, v)
	}
}

func CleanNil(data map[string]interface{}) map[string]interface{} {
	cleaned := make(map[string]interface{})
	for k, v := range data {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case map[string]interface{}:
			nested := CleanNil(val)
			if len(nested) > 0 {
				cleaned[k] = nested
			}
		case []interface{}:
			cleanedSlice := cleanSlice(val)
			if len(cleanedSlice) > 0 {
				cleaned[k] = cleanedSlice
			}
		default:
			cleaned[k] = v
		}
	}
	return cleaned
}

func cleanSlice(slice []interface{}) []interface{} {
	result := make([]interface{}, 0, len(slice))
	for _, v := range slice {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case map[string]interface{}:
			cleaned := CleanNil(val)
			if len(cleaned) > 0 {
				result = append(result, cleaned)
			}
		case []interface{}:
			nested := cleanSlice(val)
			if len(nested) > 0 {
				result = append(result, nested)
			}
		default:
			result = append(result, v)
		}
	}
	return result
}

func EnsureUUID(doc map[string]interface{}) error {
	val := doc["uuid"]
	uuid := ""
	if val != nil {
		uuid = fmt.Sprintf("%v", val)
		if len(uuid) > 0 && len(uuid) != 32 {
			log.Info("invalid uuid: %s. A new value will be generated", uuid)
		} else if len(uuid) == 32 {
			return nil
		}
	}
	uuid, err := GenerateUUID()
	if err != nil {
		return fmt.Errorf("could not generate a uuid: %v", err)
	}
	doc["uuid"] = uuid
	return nil
}
