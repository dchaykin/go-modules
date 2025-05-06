package database

type DomainEntity interface {
	UUID() string
	DatabaseName() string
}
