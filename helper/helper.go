package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/dchaykin/go-modules/datamodel"
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

func EnsureUUID(domainEntity datamodel.DomainEntity) error {
	uuid := domainEntity.UUID()
	if len(uuid) > 0 && len(uuid) != 32 {
		log.Info("invalid uuid: %s. A new value will be generated", uuid)
	} else if len(uuid) == 32 {
		return nil
	}

	uuid, err := datamodel.GenerateUUID()
	if err != nil {
		return fmt.Errorf("could not generate a uuid: %v", err)
	}
	domainEntity.SetUUID(uuid)
	return nil
}

func ValueString(fields map[string]any, fieldName string) string {
	value, ok := fields[fieldName]
	if !ok || value == nil {
		return ""
	}
	return fmt.Sprintf("%s", value)
}

func FloatFromString(value string) float64 {
	if value == "" {
		return 0
	}
	result, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return result
	}
	log.Errorf("Could not parse %s into float: %v", value, err)
	return 0
}
