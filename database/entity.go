package database

import (
	"fmt"

	"github.com/dchaykin/mygolib/log"
	"go.mongodb.org/mongo-driver/bson"
)

func FindDomainEntityByUUID(uuid string, domainEntity DomainEntity) (bool, error) {
	session, err := OpenSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	return session.GetEntityByUUID(uuid, domainEntity)
}

func GetDomainEntityByUUID(uuid string, domainEntity DomainEntity) error {
	bFound, err := FindDomainEntityByUUID(uuid, domainEntity)
	if err != nil {
		return err
	}
	if !bFound {
		return fmt.Errorf("no record with UUID %s found", uuid)
	}
	return nil
}

func ReadDomainEntities(session DatabaseSession, domainEntity DomainEntity, offset, limit int64) ([]DomainEntity, error) {
	coll := session.GetCollection(domainEntity.DatabaseName(), domainEntity.CollectionName())
	dataList := []any{}
	sortOpt := bson.D{{Key: "uuid", Value: 1}}
	count, err := session.Extract(coll, nil, &dataList, sortOpt, offset, limit)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	resultList, err := convertToDomainEntities(dataList, domainEntity)
	return resultList, log.WrapError(err)
}

func convertToDomainEntities(sourceList []any, domainEntity DomainEntity) (resultList []DomainEntity, err error) {
	for i, item := range sourceList {
		o, err := bson.Marshal(item)
		if err != nil {
			return nil, log.WrapError(err)
		}

		entity := domainEntity.CreateEmpty()
		err = bson.Unmarshal(o, entity)
		if err != nil {
			return nil, log.WrapError(fmt.Errorf("failed to unmarshal item %d: %w", i, err))
		}
		resultList = append(resultList, entity)
	}

	return resultList, nil
}
