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
	Entity() map[string]interface{}
	SetMetaData(userIdentity auth.UserIdentity, userRole string)

	ValueString(fieldName string) string
	ValueInt(fieldName string) int
	ValueFloat(fieldName string) float32
	ValueDate(fieldName string) *time.Time
	ValueBool(fieldName string) bool
}
