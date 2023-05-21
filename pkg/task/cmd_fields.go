package task

import "strings"

func RenderAllowedFields() {
	displayGraph("| Fields | Description |")
	displayGraph("| ------ | ----------- |")
	for _, info := range &allowedFieldsSlice {
		displayGraph("|", strings.Join(info.allowedNames, ","), "|", info.description, "|")
	}
}
