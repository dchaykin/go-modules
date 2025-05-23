package datamodel

import (
	"time"

	"github.com/dchaykin/go-modules/auth"
)

type DomainEntity interface {
	UUID() string
	SetUUID(uuid string)
	DatabaseName() string
	CollectionName() string
	Entity() map[string]any
	OverviewRow() map[string]any
	SetMetaData(userIdentity auth.UserIdentity, userRole string)
	GetAccessConfig() []AccessConfig

	ValueString(fieldName string) string
	ValueInt(fieldName string) int
	ValueFloat(fieldName string) float32
	ValueDate(fieldName string) *time.Time
	ValueBool(fieldName string) bool
}

type DomainItem struct {
	record map[string]any
}

func (item *DomainItem) Set(data map[string]any) {
	item.record = data
}

func (item DomainItem) AsBool(fieldName string, defaultValue bool) bool {
	value, ok := item.record[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return value.(bool)
}

func (item DomainItem) AsString(fieldName string, defaultValue string) string {
	value, ok := item.record[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return value.(string)
}

func (item DomainItem) AsInt(fieldName string, defaultValue int) int {
	value, ok := item.record[fieldName]
	if !ok || value == nil {
		return defaultValue
	}
	return value.(int)
}
