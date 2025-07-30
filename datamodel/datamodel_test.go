package datamodel

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dchaykin/go-modules/database"
	"github.com/stretchr/testify/require"
)

func TestTenantConfig(t *testing.T) {
	tc, err := LoadDataModelByRole("testdata-001", "customer")
	require.NoError(t, err)

	require.Nil(t, tc.Roles)
	require.NotNil(t, tc.Cmbs)

	require.Exactly(t, 2, len(*tc.Cmbs))
	require.Greater(t, len((*tc.Cmbs)["roles"]), 0)

	userPartner, ok := (*tc.Cmbs)["user"]["partner"]
	require.Equal(t, true, ok)
	require.NotNil(t, userPartner)
	require.Equal(t, "/app-config/api/partner/cmbs/userPartner", *userPartner.Source)
	require.Equal(t, "userPartner", userPartner.Name)

	roleNames, ok := (*tc.Cmbs)["roles"]["name"]
	require.Equal(t, true, ok)
	require.NotNil(t, roleNames)
	require.Equal(t, false, *roleNames.Translate)
	require.Equal(t, 3, len(roleNames.Content))
}

func TestCustomFieldValueByType(t *testing.T) {
	field := CustomField{
		"type": FieldTypeDate,
	}

	require.Equal(t, "2025-06-16", field.ValueByType("2025-06-16T00:00:00.000+02:00"))

	field["type"] = FieldTypeDateTime
	require.Equal(t, "2025-06-16T00:00:00.000+02:00", field.ValueByType("2025-06-16T00:00:00.000+02:00"))

	field["type"] = FieldTypeDate
	require.Equal(t, "2025", field.ValueByType("2025"))

	field = CustomField{
		"type": FieldTypeString,
		"size": 1024,
	}
	require.Equal(t, int64(1024), field.Size())

}

type testDomainEntity struct {
	Record
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

func TestCleanNil(t *testing.T) {
	rec := Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	rec.CleanNil()

	expected := map[string]any{}
	err = json.Unmarshal([]byte(`{"foo":{"boz":234,"bar":[{"baz":1},{"baz":2}],"foz":123}}`), &expected)
	require.NoError(t, err)

	require.EqualValues(t, true, reflect.DeepEqual(expected, rec.Fields))
}
