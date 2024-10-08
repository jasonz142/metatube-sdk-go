package testkit

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValidateFunc func(*testing.T, any)

func logJSONContent() ValidateFunc {
	return func(t *testing.T, a any) {
		data, err := json.MarshalIndent(a, "", "\t")
		require.NoError(t, err)
		// Print out JSON data.
		t.Logf("%s", data)
	}
}

func assertIsValid() ValidateFunc {
	return func(t *testing.T, a any) {
		if v, ok := a.(interface{ Valid() bool }); ok {
			assert.True(t, v.Valid())
			return
		}
		// must be a slice.
		require.Equal(t, reflect.Slice, reflect.TypeOf(a).Kind())
		s := reflect.ValueOf(a)
		for i := 0; i < s.Len(); i++ {
			x := s.Index(i).Interface()
			require.Implements(t, (*interface{ Valid() bool })(nil), x)
			assert.True(t, x.(interface{ Valid() bool }).Valid())
		}
	}
}

func FieldsNotEmpty(fields ...string) ValidateFunc {
	return func(t *testing.T, a any) {
		for _, field := range fields {
			f, err := getStructFieldByName(a, field)
			require.NoError(t, err)
			assert.NotEmpty(t, f)
		}
	}
}

func FieldsNotEmptyAny(fields ...string) ValidateFunc {
	return func(t *testing.T, a any) {
		ok := false
		for _, field := range fields {
			f, err := getStructFieldByName(a, field)
			if err != nil {
				continue
			}
			z := reflect.Zero(reflect.ValueOf(f).Type())
			if !reflect.DeepEqual(f, z.Interface()) {
				ok = true
			}
		}
		assert.Truef(t, ok, "fields are all empty: %v", fields)
	}
}
