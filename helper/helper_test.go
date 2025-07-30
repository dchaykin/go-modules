package helper

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dchaykin/go-modules/database"
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
	datamodel.Record
}

func (de testDomainEntity) DatabaseName() string {
	return "test"
}

func (de testDomainEntity) CollectionName() string {
	return "test"
}

func (de *testDomainEntity) GetAccessConfig() []database.AccessConfig {
	return nil
}

func (de testDomainEntity) OverviewRow() map[string]any {
	return nil
}

func (de testDomainEntity) BeforeSave(session database.DatabaseSession) error {
	return nil
}

func (de testDomainEntity) CreateEmpty() database.DomainEntity {
	return &testDomainEntity{}
}

func TestEnsureUUID(t *testing.T) {
	doc := testDomainEntity{}

	err := datamodel.EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, 32, len(doc.UUID()))

	doc.SetUUID("invalid-uuid")
	err = datamodel.EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, 32, len(doc.UUID()))

	uuid, err := datamodel.GenerateUUID()
	doc.SetUUID(uuid)
	require.NoError(t, err)
	err = datamodel.EnsureUUID(&doc)
	require.NoError(t, err)
	require.EqualValues(t, uuid, doc.UUID())
}

func TestCleanNil(t *testing.T) {
	rec := datamodel.Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	rec.CleanNil()

	expected := map[string]any{}
	err = json.Unmarshal([]byte(`{"foo":{"boz":234,"bar":[{"baz":1},{"baz":2}],"foz":123}}`), &expected)
	require.NoError(t, err)

	require.EqualValues(t, true, reflect.DeepEqual(expected, rec.Fields))
}

func TestFloatFromString(t *testing.T) {
	require.EqualValues(t, 0, FloatFromString(""))
	require.EqualValues(t, 6.4, FloatFromString("6.4"))
	require.EqualValues(t, 1.0, FloatFromString("1"))
	require.EqualValues(t, 0, FloatFromString("6,4"))
}
