package datamodel

type DomainEntity interface {
	UUID() string
	SetUUID(uuid string)
	DatabaseName() string
	CollectionName() string
}
