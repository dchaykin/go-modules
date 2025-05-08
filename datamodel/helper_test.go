package datamodel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testDomainEntity struct {
	ID   string
	Data map[string]interface{}
}

func (de testDomainEntity) UUID() string {
	return de.ID
}

func (de *testDomainEntity) SetUUID(uuid string) {
	de.ID = uuid
}

func (de testDomainEntity) DatabaseName() string {
	return "test"
}

func (de testDomainEntity) CollectionName() string {
	return "test"
}

func TestEnsureUUID(t *testing.T) {
	doc := testDomainEntity{}

	err := EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, 32, len(doc.UUID()))

	doc.SetUUID("invalid-uuid")
	err = EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, 32, len(doc.UUID()))

	uuid, err := GenerateUUID()
	doc.SetUUID(uuid)
	require.NoError(t, err)
	err = EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, uuid, doc.UUID())
}
