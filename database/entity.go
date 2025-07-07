package database

import (
	"fmt"

	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/log"
	"go.mongodb.org/mongo-driver/bson"
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

func ReadDomainEntities(session DatabaseSession, coll Collection, offset, limit int64) ([]datamodel.Record, error) {
	dataList := []any{}
	sortOpt := bson.D{{Key: "uuid", Value: 1}}
	count, err := session.Extract(coll, nil, &dataList, sortOpt, offset, limit)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	resultList, err := convertToDomainEntities(dataList)
	return resultList, log.WrapError(err)
}

func convertToDomainEntities(sourceList []any) (resultList []datamodel.Record, err error) {
	for i, item := range sourceList {
		o, err := bson.Marshal(item)
		if err != nil {
			return nil, err
		}

		var entity datamodel.Record
		err = bson.Unmarshal(o, &entity)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item %d: %w", i, err)
		}
		resultList = append(resultList, entity)
	}

	return resultList, nil
}
