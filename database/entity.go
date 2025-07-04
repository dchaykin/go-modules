package database

import (
	"fmt"

	"github.com/dchaykin/go-modules/datamodel"
)

func GetDomainEntityByUUID(uuid string, domainEntity datamodel.DomainEntity) error {
	session, err := OpenSession()
	if err != nil {
		return err
	}
	defer session.Close()

	bFound, err := session.GetEntityByUUID(uuid, domainEntity)
	if err != nil {
		return err
	}
	if !bFound {
		return fmt.Errorf("no record with UUID %s found", uuid)
	}
	return nil
}
