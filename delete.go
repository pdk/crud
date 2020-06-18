package crud

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/pdk/crud/describe"
)

// NewDeleter creates a deleter CRUDFunc
func NewDeleter(tableName string, exampleStruct interface{}, keyFields ...string) CRUDFunc {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build updateer: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	keyColumnNames := desc.ColumnsOf(keyFields...)
	keyIndexes := desc.IndexesOf(keyColumnNames...)

	stmt := "delete from " + tableName + " where "

	for i, c := range keyColumnNames {
		if i > 0 {
			stmt += " and "
		}

		stmt += c + " = $" + strconv.Itoa(i+1)
	}

	valueCount := len(keyIndexes)

	return func(db dbHandle, item interface{}) error {

		itemValue := reflect.ValueOf(item)
		if itemValue.Type() != exampleItemType {
			return fmt.Errorf("crud.NewUpdateer func expected a %s, but got a %s",
				exampleItemType.String(), itemValue.Type().String())
		}

		bindValues := make([]interface{}, valueCount, valueCount)
		for p, i := range keyIndexes {
			bindValues[p] = itemValue.Field(i).Interface()
		}

		log.Printf("%s", stmt)
		log.Printf("%v", bindValues)

		_, err = db.Exec(stmt, bindValues...)

		return err
	}
}
