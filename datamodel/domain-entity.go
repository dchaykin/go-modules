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

type EntityNodeList []EntityNode

func (l *EntityNodeList) Create(node any) {
	if node == nil {
		return
	}
	switch v := node.(type) {
	case []any:
		for _, item := range v {
			switch i := item.(type) {
			case map[string]any:
				{
					*l = append(*l, i)
				}
			default:
				log.Warn("unexpected type of DomainItemList item (expected map): %T, %v", item, item)
			}
		}
	default:
		log.Warn("unexpected type of DomainItemList (expected slice): %T, %v", node, node)
	}
}

func (l EntityNodeList) UniqueKeyList(fieldName string, separator string) string {
	return l.UniqueKeysList([]string{fieldName}, "", separator)
}

func (l EntityNodeList) UniqueKeysList(fieldNames []string, fieldSeparator, itemSeparator string) string {
	result := []string{}
	for _, node := range l {
		if node == nil {
			continue
		}
		result = append(result, node.Uniques(fieldNames, fieldSeparator))
	}
	return strings.Join(result, itemSeparator)
}

type EntityNode map[string]any

func (node EntityNode) UUID() string {
	return node.AsString("uuid", "")
}

func (node EntityNode) AsBool(fieldName string, defaultValue bool) bool {
	value, ok := node[fieldName]
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

func (node EntityNode) AsString(fieldName string, defaultValue string) string {
	value, ok := node[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return fmt.Sprintf("%v", value)
}

func (node EntityNode) AsInt(fieldName string, defaultValue int) int {
	value, ok := node[fieldName]
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

func (node EntityNode) Uniques(fieldNames []string, separator string) string {
	result := []string{}
	for _, fieldName := range fieldNames {
		value := node.AsString(fieldName, "")
		if value == "" || slices.Contains(result, value) {
			continue
		}
		result = append(result, value)
	}
	return strings.Join(result, separator)
}

func (node EntityNode) Unique(fieldName string) string {
	return node.Uniques([]string{fieldName}, "")
}
