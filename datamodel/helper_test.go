package datamodel

import (
	"testing"
	"time"

	"github.com/dchaykin/go-modules/auth"
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

func (de testDomainEntity) Entity() map[string]interface{} {
	return de.Data
}

func (de testDomainEntity) ValueString(fieldName string) string {
	return ""
}

func (de testDomainEntity) ValueInt(fieldName string) int {
	return 0
}

func (de testDomainEntity) ValueFloat(fieldName string) float32 {
	return 0
}

func (de testDomainEntity) ValueDate(fieldName string) *time.Time {
	return nil
}

func (de testDomainEntity) ValueBool(fieldName string) bool {
	return false
}

func (de *testDomainEntity) SetMetaData(userIdentity auth.UserIdentity, userRole string) {

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
