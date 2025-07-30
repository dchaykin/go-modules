package datamodel

import (
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/database"
	"github.com/dchaykin/go-modules/httpcomm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metadata struct {
	Timestamp time.Time `bson:"timestamp"`
	User      string    `bson:"user"`
	Partner   string    `bson:"partner"`
	Role      string    `bson:"role"`
}

type Mapper struct {
	Cmbs     map[string]any `json:"cmbs" bson:"cmbs"`
	Richtext map[string]any `json:"richtext" bson:"richtext"`
}

type Record struct {
	Metadata Metadata       `bson:"metadata"`
	Fields   map[string]any `json:"entity" bson:"entity"`
	Mapper   Mapper         `json:"mapper" bson:"mapper"`

	userIdentity auth.UserIdentity `bson:"-"`
}

func (r Record) UserIdentity() auth.UserIdentity {
	return r.userIdentity
}

func (r *Record) SetUserIdentity(userIdentity auth.UserIdentity) {
	r.userIdentity = userIdentity
}

func (r *Record) CleanNil() {
	r.Fields = r.cleanNil(r.Fields)
}

func (r *Record) SetValue(key string, value any) {
	if r.Fields == nil {
		r.Fields = make(map[string]any)
	}
	r.Fields[key] = value
}

func (r *Record) GetValue(key string) any {
	if r.Fields == nil {
		return nil
	}
	value, ok := r.Fields[key]
	if !ok {
		return nil
	}
	return value
}

func (r *Record) AddRecord(key string, value any) {
	if r.Fields == nil {
		r.Fields = make(map[string]any)
	}
	_, ok := r.Fields[key]
	if !ok {
		r.Fields[key] = []any{value}
	} else {
		r.Fields[key] = append(r.Fields[key].([]any), value)
	}
}

func (r *Record) cleanNil(data map[string]any) map[string]any {
	cleaned := make(map[string]any)
	for k, v := range data {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case map[string]any:
			nested := r.cleanNil(val)
			if len(nested) > 0 {
				cleaned[k] = nested
			}
		case []any:
			cleanedSlice := r.cleanSlice(val)
			if len(cleanedSlice) > 0 {
				cleaned[k] = cleanedSlice
			}
		default:
			cleaned[k] = v
		}
	}
	return cleaned
}

func (r *Record) NormalizePrimitives() {
	r.Fields = r.normalizePrimitives(r.Fields)
}

func (r *Record) normalizePrimitives(m map[string]any) map[string]any {
	for k, v := range m {
		switch val := v.(type) {
		case primitive.A:
			arr := make([]any, len(val))
			copy(arr, val)
			m[k] = arr
		case map[string]any:
			m[k] = r.normalizePrimitives(val)
		}
	}
	return m
}

func (r *Record) cleanSlice(slice []any) []any {
	result := make([]any, 0, len(slice))
	for _, v := range slice {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case map[string]any:
			cleaned := r.cleanNil(val)
			if len(cleaned) > 0 {
				result = append(result, cleaned)
			}
		case []any:
			nested := r.cleanSlice(val)
			if len(nested) > 0 {
				result = append(result, nested)
			}
		default:
			result = append(result, v)
		}
	}
	return result
}

func (r *Record) SetMetadata(appName string) {
	r.Metadata.Timestamp = time.Now()
	r.Metadata.Partner = r.userIdentity.Partner()
	r.Metadata.Role = r.userIdentity.RoleByApp(appName)
	r.Metadata.User = r.userIdentity.Username()
}

func (r *Record) BeforeSave(session database.DatabaseSession) error {
	return nil
}

type OnJsonArrayFound func(array []any)

func (r *Record) FindJsonArray(jsonPath []string, f OnJsonArrayFound) {
	r.findJsonArray(r.Fields, jsonPath, f)
}

func (r *Record) findJsonArray(node map[string]any, jsonPath []string, f OnJsonArrayFound) {
	if len(jsonPath) == 0 {
		return
	}
	lookedNodeName := jsonPath[0]
	for k := range node {
		if k == lookedNodeName {
			switch vTyped := node[k].(type) {
			case []any:
				f(vTyped)
				return
			case map[string]any:
				r.findJsonArray(vTyped, jsonPath[1:], f)
			}
		}
	}
}

type OnJsonFieldFound func(field map[string]any, name string)

func (r *Record) FindJsonField(jsonPath []string, f OnJsonFieldFound) {
	r.findJsonField(r.Fields, jsonPath, f)
}

func (r *Record) findJsonField(node map[string]any, jsonPath []string, f OnJsonFieldFound) {
	if len(jsonPath) == 0 {
		return
	}
	lookedNodeName := jsonPath[0]
	for k, v := range node {
		if k == lookedNodeName {
			switch vTyped := v.(type) {
			case nil:
				if len(jsonPath) == 1 {
					f(node, k)
				}
			case []any:
				for i := range vTyped {
					if vTyped[i] == nil {
						continue
					}
					r.findJsonField(vTyped[i].(map[string]any), jsonPath[1:], f)
				}
				return
			case map[string]any:
				r.findJsonField(node[k].(map[string]any), jsonPath[1:], f)
			case any:
				if len(jsonPath) == 1 {
					f(node, k)
				}
			}
		}
	}
}

func (r Record) UUID() string {
	if uuid, ok := r.Fields["uuid"]; ok {
		return fmt.Sprintf("%s", uuid)
	}
	return ""
}

func (r *Record) SetUUID(UUID string) {
	r.SetValue("uuid", UUID)
}

func (r Record) Entity() map[string]any {
	return r.Fields
}

func GetErrorResponse(err error) *httpcomm.ServiceResponse {
	result := httpcomm.ServiceResponse{Error: new(string)}
	*result.Error = fmt.Sprintf("%v", err)
	return &result
}

func (r *Record) ApplyMapper() {
	r.applyMapper(0, r.Fields, r.Mapper.Cmbs, r.GetOnFoundNewMapping(FieldTypeCombobox))
	r.applyMapper(0, r.Fields, r.Mapper.Richtext, r.GetOnFoundNewMapping(FieldTypeRichtext))
}

func (r *Record) GetOnFoundNewMapping(fieldType string) OnFoundNewMapping {
	if fieldType == FieldTypeRichtext {
		return func(key, indexStr string, oldValue, newValue any) any {
			str := fmt.Sprintf("%v", oldValue)
			if len(str) > 1024 {
				return truncateString(fmt.Sprintf("%v", newValue), 1024)
			}
			return oldValue
		}
	}
	return nil
}

func truncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) > maxLen {
		runes := []rune(s)
		return string(runes[:maxLen-3]) + "..."
	}
	return s
}

type OnFoundNewMapping func(key string, indexStr string, oldValue any, newValue any) any

func (r *Record) applyMapper(index int, fields map[string]any, mapper map[string]any, callback OnFoundNewMapping) {
	for key, mapping := range mapper {
		subMap, ok := mapping.(map[string]any)
		if !ok {
			continue
		}

		fieldVal, fieldExists := fields[key]
		if !fieldExists {
			continue
		}

		switch fv := fieldVal.(type) {
		case map[string]any:
			// z. B. entity["foo"] → {"bar": ..., "x": ...}
			r.applyMapper(index, fv, subMap, callback)

		case []any:
			// z. B. entity["style"] → []map[string]any
			for x, item := range fv {
				if m, ok := item.(map[string]any); ok {
					r.applyMapper(x, m, subMap, callback)
				}
			}
		default:
			indexStr := fmt.Sprintf("%d", index)
			// Basiswert: versuche, passenden Mapper-Eintrag zu finden
			keyStr := toKeyString(fv)
			newVal, ok := subMap[keyStr]
			if !ok {
				newVal, ok = subMap[indexStr]
			}
			if ok {
				if callback != nil {
					newVal = callback(key, indexStr, fv, newVal)
				}
				fields[key] = newVal
			}
		}
	}
}

func toKeyString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
