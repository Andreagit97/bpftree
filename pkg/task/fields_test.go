package task

import (
	"testing"
)

// This test should avoid collision if we try to introduce new fields with the same name.
func TestFieldNameCollision(t *testing.T) {
	fieldNames := make(map[string]struct{})
	var member struct{}
	for _, info := range allowedFieldsSlice {
		for _, name := range info.allowedNames {
			if _, ok := fieldNames[name]; ok {
				t.Errorf("There are 2 fields with the same name '%s'\n", name)
			}
			fieldNames[name] = member
		}
	}
}
