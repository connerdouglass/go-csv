package csv

import "reflect"

type structTag struct {
	ColumnName string
	FieldName  string
}

func parseStructTag(field reflect.StructField) *structTag {
	// Get the "csv" tag from the struct field
	tagStr, ok := field.Tag.Lookup("csv")
	if !ok {
		return nil
	}

	// Parse the tag
	var result structTag
	result.FieldName = field.Name
	parseStructTagString(&result, tagStr)
	return &result
}

func parseStructTagString(tag *structTag, str string) {
	tag.ColumnName = str
}
