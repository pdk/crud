package crud

import (
	"log"
	"reflect"

	"github.com/pdk/crud/describe"
)

// NewQueryFunc creates a function to execute a query and return all results as
// a slice of the original struct type.
func NewQueryFunc(tableName string, exampleStruct interface{}) QueryFunc {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build query func: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	colCount := len(desc.Columns())

	return func(db dbHandle, querySQL string, queryParameters ...interface{}) (interface{}, error) {

		values := reflect.MakeSlice(reflect.SliceOf(exampleItemType), 0, 0)

		rows, err := db.Query(querySQL, queryParameters...)
		if err != nil {
			return values.Interface(), err
		}
		defer rows.Close()

		for rows.Next() {

			nextStructPtr := reflect.New(exampleItemType)
			nextStruct := nextStructPtr.Elem()

			scanItems := []interface{}{}
			for i := 0; i < colCount; i++ {
				nextFld := nextStruct.Field(i).Addr().Interface()
				scanItems = append(scanItems, nextFld)
			}

			err := rows.Scan(scanItems...)
			if err != nil {
				return values.Interface(), err
			}

			values = reflect.Append(values, nextStruct)
		}

		return values.Interface(), rows.Err()
	}
}

// NewQueryOneRowFunc creates a function to execute a QueryOneRow and return an
// instance of the original struct type.
func NewQueryOneRowFunc(tableName string, exampleStruct interface{}) QueryFunc {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build query func: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	colCount := len(desc.Columns())

	return func(db dbHandle, querySQL string, queryParameters ...interface{}) (interface{}, error) {

		nextStructPtr := reflect.New(exampleItemType)
		nextStruct := nextStructPtr.Elem()

		scanItems := []interface{}{}
		for i := 0; i < colCount; i++ {
			nextFld := nextStruct.Field(i).Addr().Interface()
			scanItems = append(scanItems, nextFld)
		}

		err := db.QueryRow(querySQL, queryParameters...).
			Scan(scanItems...)

		return nextStruct.Interface(), err
	}
}
