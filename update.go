package crud

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/pdk/crud/describe"
)

// NewUpdater creates a new updater CRUDFunc.
func NewUpdater(tableName string, exampleStruct interface{}, keyFields ...string) CRUDFunc {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build updateer: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	setColumnNames := desc.ColumnsOmitFields(keyFields...)
	setIndexes := desc.IndexesOf(setColumnNames...)

	keyColumnNames := desc.ColumnsOf(keyFields...)
	keyIndexes := desc.IndexesOf(keyColumnNames...)

	stmt := "update " + tableName + " set "

	for i, c := range setColumnNames {
		if i > 0 {
			stmt += ", "
		}

		stmt += c + " = $" + strconv.Itoa(i+1)
	}

	stmt += " where "

	for i, c := range keyColumnNames {
		if i > 0 {
			stmt += " and "
		}

		stmt += c + " = $" + strconv.Itoa(i+1+len(setColumnNames))
	}

	valueCount := len(setIndexes) + len(keyIndexes)

	return func(db dbHandle, item interface{}) error {

		itemValue := reflect.ValueOf(item)
		if itemValue.Type() != exampleItemType {
			return fmt.Errorf("crud.NewUpdateer func expected a %s, but got a %s",
				exampleItemType.String(), itemValue.Type().String())
		}

		bindValues := make([]interface{}, valueCount, valueCount)
		for p, i := range setIndexes {
			bindValues[p] = itemValue.Field(i).Interface()
		}
		for p, i := range keyIndexes {
			bindValues[p+len(setIndexes)] = itemValue.Field(i).Interface()
		}

		log.Printf("%s", stmt)
		log.Printf("%v", bindValues)

		_, err = db.Exec(stmt, bindValues...)

		return err
	}
}
