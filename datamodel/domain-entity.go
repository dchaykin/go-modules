package datamodel

import "time"

type DomainEntity interface {
	UUID() string
	SetUUID(uuid string)
	DatabaseName() string
	CollectionName() string
	Entity() map[string]interface{}

	ValueString(fieldName string) string
	ValueInt(fieldName string) int
	ValueFloat(fieldName string) float32
	ValueDate(fieldName string) time.Time
	ValueBool(fieldName string) bool
}
