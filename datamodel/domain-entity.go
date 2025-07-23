package datamodel

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

type DomainEntity interface {
	UUID() string
	SetUUID(uuid string)
	CreateEmpty() DomainEntity
	SetValue(key string, value any)
	GetValue(key string) any
	DatabaseName() string
	CollectionName() string
	Entity() map[string]any
	OverviewRow() map[string]any
	SetMetadata(userIdentity auth.UserIdentity, subject string)
	GetAccessConfig() []AccessConfig
	CleanNil()
	BeforeSave() error

	NormalizePrimitives()
	ApplyMapper()
}

type DomainItemList []any

func (l DomainItemList) UniqueKeyList(fieldName string, separator string) string {
	return l.UniqueKeysList([]string{fieldName}, "", separator)
}

func (l DomainItemList) UniqueKeysList(fieldNames []string, fieldSeparator, itemSeparator string) string {
	result := []string{}
	for _, v := range l {
		if v == nil {
			continue
		}
		item := DomainItem(v.(map[string]any))
		result = append(result, item.Uniques(fieldNames, fieldSeparator))
	}
	return strings.Join(result, itemSeparator)
}

type DomainItem map[string]any

func (item DomainItem) AsBool(fieldName string, defaultValue bool) bool {
	value, ok := item[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	switch v := value.(type) {
	case bool:
		return v
	default:
		return defaultValue
	}
}

func (item DomainItem) AsString(fieldName string, defaultValue string) string {
	value, ok := item[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return fmt.Sprintf("%v", value)
}

func (item DomainItem) AsInt(fieldName string, defaultValue int) int {
	value, ok := item[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Warn("could not convert string to int: %s, return default value %d", v, defaultValue)
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func (item DomainItem) Uniques(fieldNames []string, separator string) string {
	result := []string{}
	for _, fieldName := range fieldNames {
		value := item.AsString(fieldName, "")
		if value == "" || slices.Contains(result, value) {
			continue
		}
		result = append(result, value)
	}
	return strings.Join(result, separator)
}

func (item DomainItem) Unique(fieldName string) string {
	return item.Uniques([]string{fieldName}, "")
}
