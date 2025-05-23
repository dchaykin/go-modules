package helper

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/datamodel"
	"github.com/stretchr/testify/require"
)

const testJsonRecord = `{
			"metadata": {},
			"entity": {
				"foo": {
					"boz": 234,
					"bar": [
						{ "baz": 1 },
						{ "baz": 2 },
						null
					],
					"foz": 123,
					"emptyf": null,
					"emptyl": []
				}
			}
		}`

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

func (de testDomainEntity) Entity() map[string]any {
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

func (de *testDomainEntity) GetAccessConfig() []datamodel.AccessConfig {
	return nil
}

func (de testDomainEntity) OverviewRow() map[string]any {
	return nil
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

	uuid, err := datamodel.GenerateUUID()
	doc.SetUUID(uuid)
	require.NoError(t, err)
	err = EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, uuid, doc.UUID())
}

func TestCleanNil(t *testing.T) {
	rec := datamodel.Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	rec.Fields = CleanNil(rec.Fields)

	expected := map[string]any{}
	err = json.Unmarshal([]byte(`{"foo":{"boz":234,"bar":[{"baz":1},{"baz":2}],"foz":123}}`), &expected)
	require.NoError(t, err)

	require.EqualValues(t, true, reflect.DeepEqual(expected, rec.Fields))
}
