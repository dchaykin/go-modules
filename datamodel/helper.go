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

func CleanNil(data map[string]any) map[string]any {
	cleaned := make(map[string]any)
	for k, v := range data {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case map[string]any:
			nested := CleanNil(val)
			if len(nested) > 0 {
				cleaned[k] = nested
			}
		case []any:
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

func cleanSlice(slice []any) []any {
	result := make([]any, 0, len(slice))
	for _, v := range slice {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case map[string]any:
			cleaned := CleanNil(val)
			if len(cleaned) > 0 {
				result = append(result, cleaned)
			}
		case []any:
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

func EnsureUUID(domainEntity DomainEntity) error {
	uuid := domainEntity.UUID()
	if len(uuid) > 0 && len(uuid) != 32 {
		log.Info("invalid uuid: %s. A new value will be generated", uuid)
	} else if len(uuid) == 32 {
		return nil
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return fmt.Errorf("could not generate a uuid: %v", err)
	}
	domainEntity.SetUUID(uuid)
	return nil
}
