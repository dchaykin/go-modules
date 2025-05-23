package datamodel

import (
	"encoding/json"
	"reflect"
	"testing"

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

func TestFindJsonArray_FindsArray(t *testing.T) {
	rec := Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	var found []any
	rec.FindJsonArray([]string{"foo", "bar"}, func(arr []any) {
		found = arr
	})

	expected := []any{
		map[string]any{"baz": 1.0},
		map[string]any{"baz": 2.0},
		nil,
	}
	if !reflect.DeepEqual(found, expected) {
		require.EqualValues(t, expected, found)
	}
}

func TestFindJsonArray_NotFound(t *testing.T) {
	rec := Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	called := false
	rec.FindJsonArray([]string{"foo", "notfound"}, func(arr []any) {
		called = true
	})

	require.EqualValues(t, called, false)
}

func TestFindJsonField_FindsField(t *testing.T) {
	rec := Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	var found map[string]any
	var fieldName string
	rec.FindJsonField([]string{"foo", "foz"}, func(field map[string]any, name string) {
		found = field
		fieldName = name
	})

	require.EqualValues(t, 123, found[fieldName])
}

func TestFindJsonField_FindsFieldInArray(t *testing.T) {
	rec := Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	var founds []map[string]any
	rec.FindJsonField([]string{"foo", "bar", "baz"}, func(field map[string]any, name string) {
		founds = append(founds, field)
	})

	expected := []map[string]any{
		{"baz": 1.0},
		{"baz": 2.0},
	}
	if !reflect.DeepEqual(founds, expected) {
		require.EqualValues(t, expected, founds)
	}
}

func TestFindJsonField_NotFound(t *testing.T) {
	rec := Record{}
	err := json.Unmarshal([]byte(testJsonRecord), &rec)
	require.NoError(t, err)

	called := false
	rec.FindJsonField([]string{"foo", "notfound"}, func(field map[string]any, name string) {
		called = true
	})

	if called {
		require.EqualValues(t, called, false)
	}
}
