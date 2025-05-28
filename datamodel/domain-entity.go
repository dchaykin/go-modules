package datamodel

import (
	"slices"
	"strings"

	"github.com/dchaykin/go-modules/auth"
)

type DomainEntity interface {
	UUID() string
	SetUUID(uuid string)
	DatabaseName() string
	CollectionName() string
	Entity() map[string]any
	OverviewRow() map[string]any
	SetMetadata(userIdentity auth.UserIdentity, subject string)
	GetAccessConfig() []AccessConfig
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
	return value.(bool)
}

func (item DomainItem) AsString(fieldName string, defaultValue string) string {
	value, ok := item[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return value.(string)
}

func (item DomainItem) AsInt(fieldName string, defaultValue int) int {
	value, ok := item[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return value.(int)
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
