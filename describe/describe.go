package describe

import (
	"fmt"
	"reflect"
	"strings"
)

// StructDescription provides metadata about a struct
type StructDescription struct {
	Name   string
	Fields []FieldDescription
}

// Names returns the names of the fields
func (desc StructDescription) Names() []string {
	names := []string{}
	for _, fld := range desc.Fields {
		names = append(names, fld.Name)
	}

	return names
}

// Columns returns the database column names
func (desc StructDescription) Columns() []string {
	cols := []string{}
	for _, fld := range desc.Fields {
		cols = append(cols, fld.Column)
	}

	return cols
}

func (desc StructDescription) ColumnsOmitFields(omitList ...string) []string {
	cols := []string{}
	for _, fld := range desc.Fields {
		if !fieldInList(fld, omitList...) {
			cols = append(cols, fld.Column)
		}
	}

	return cols
}

func (desc StructDescription) ColumnsOf(includeList ...string) []string {
	cols := []string{}
	for _, fld := range desc.Fields {
		if fieldInList(fld, includeList...) {
			cols = append(cols, fld.Column)
		}
	}

	return cols
}

func (desc StructDescription) IndexesOf(includeList ...string) []int {
	indexes := []int{}
	for i, fld := range desc.Fields {
		if fieldInList(fld, includeList...) {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

// fieldInList returns true if the field matches anything in the skiplist.
func fieldInList(fld FieldDescription, matchList ...string) bool {
	for _, match := range matchList {
		if strings.EqualFold(fld.Name, match) || strings.EqualFold(fld.Column, match) {
			return true
		}
	}
	return false
}

func (desc StructDescription) ColumnNameOfField(idField string) string {
	for _, fld := range desc.Fields {
		if strings.EqualFold(idField, fld.Name) {
			return fld.Column
		}
	}

	return ""
}

func (desc StructDescription) FieldIndex(field string) int {
	for i, fld := range desc.Fields {
		if strings.EqualFold(field, fld.Name) {
			return i
		}
	}

	return -1
}

func (desc StructDescription) FieldIndexesOmitField(idField string) []int {
	indx := []int{}
	for i, fld := range desc.Fields {
		if strings.EqualFold(idField, fld.Name) {
			continue
		}
		indx = append(indx, i)
	}

	return indx
}

// Labels returns the labels
func (desc StructDescription) Labels() []string {
	labels := []string{}
	for _, fld := range desc.Fields {
		labels = append(labels, fld.Label)
	}

	return labels
}

// Types returns the types
func (desc StructDescription) Types() []string {
	types := []string{}
	for _, fld := range desc.Fields {
		types = append(types, fld.Type)
	}

	return types
}

// FieldDescription provides metadata about struct fields
type FieldDescription struct {
	Name   string
	Column string
	Label  string
	Type   string
}

// Describe will introspect a struct and return a StructDescription
func Describe(thing interface{}) (StructDescription, error) {

	str := reflect.TypeOf(thing)
	if str.Kind() != reflect.Struct {
		return StructDescription{}, fmt.Errorf("is not a struct %v", thing)
	}

	result := StructDescription{
		Name: str.Name(),
	}

	for i := 0; i < str.NumField(); i++ {
		fld := str.Field(i)

		result.Fields = append(result.Fields, FieldDescription{
			Name:   fld.Name,
			Column: columnName(fld),
			Label:  label(fld),
			Type:   fld.Type.Name(),
		})
	}

	return result, nil
}

func columnName(fld reflect.StructField) string {
	s := fld.Tag.Get("db")
	if s != "" {
		return s
	}

	return strings.ToLower(fld.Name)
}

func label(fld reflect.StructField) string {
	s := fld.Tag.Get("label")
	if s != "" {
		return s
	}

	return fld.Name
}
