package datamodel

import (
	"testing"

	"github.com/dchaykin/go-modules/database"
	"github.com/dchaykin/go-modules/helper"
	"github.com/stretchr/testify/require"
)

type location struct {
	Record
}

func (l location) CollectionName() string {
	return "location"
}

func (l location) DatabaseName() string {
	return "masterData"
}

func (l location) CreateEmpty() database.DomainEntity {
	return &location{}
}

func (l *location) GetAccessConfig() []database.AccessConfig {
	return nil
}

func (l *location) OverviewRow() map[string]any {
	return l.Fields
}

func TestReadDomainEntites(t *testing.T) {
	helper.LoadAccessData("../.do-not-commit/env.vars")

	session, err := database.OpenSession()
	require.NoError(t, err)
	defer session.Close()

	offset := int64(0)

	entities, err := database.ReadDomainEntities(session, &location{}, offset, 3)
	require.NoError(t, err)
	require.EqualValues(t, 3, len(entities), "Expected to read 3 entities from the collection.")

	t.Logf("Found %d entities in the first transaction.", len(entities))

	offset += int64(len(entities))

	entities, err = database.ReadDomainEntities(session, &location{}, offset, 3)
	require.NoError(t, err)
	require.Greater(t, len(entities), 0)
	t.Logf("Found %d entities in the second transaction.", len(entities))

	offset += int64(len(entities))

	entities, err = database.ReadDomainEntities(session, &location{}, offset, 3)
	require.NoError(t, err)
	require.EqualValues(t, 0, len(entities))
	t.Logf("No records to read in the third transaction.")
}
