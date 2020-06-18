package crud

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pdk/crud/describe"
)

// NewAutoIncrIDInserter generates a new Inserter. To be called on startup. Will
// terminate program if there is any error. The item value passed to the
// generated function MUST be a pointer to the type of struct passed as example.
// f := NewInserterWithAutoIncrID("foo", "ID", DaStruct{})
// err := f(db, &daStructInstance)
func NewAutoIncrIDInserter(tableName string, exampleStruct interface{}, idField string) CRUDFuncGetID {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build inserter: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	insertColumnNames := desc.ColumnsOmitFields(idField)
	insertColumnNamesStr := strings.Join(insertColumnNames, ", ")
	insertColumnIndexes := desc.IndexesOf(insertColumnNames...)
	valueCount := len(insertColumnNames)

	idColumnName := desc.ColumnNameOfField(idField)

	insertStmt := fmt.Sprintf("insert into %s (%s) values (%s) returning %s",
		tableName, insertColumnNamesStr, markers(len(insertColumnIndexes)), idColumnName)

	return func(db dbHandle, item interface{}) (int64, error) {

		itemValue := reflect.ValueOf(item)
		if itemValue.Type() != exampleItemType {
			return 0, fmt.Errorf("crud.NewAutoIncrIDInserter func expected a %s, but got a %s",
				exampleItemType.String(), itemValue.Type().String())
		}

		insertValues := make([]interface{}, valueCount, valueCount)
		for i, index := range insertColumnIndexes {
			insertValues[i] = itemValue.Field(index).Interface()
		}

		var newID int64
		err = db.QueryRow(insertStmt, insertValues...).Scan(&newID)

		return newID, err
	}
}

// NewInserter generates a new Inserter. To be called on startup. Will
// terminate program if there is any error. The item value passed to the
// generated function MUST be a pointer to the type of struct passed as example.
// f := NewInserter("foo", "ID", DaStruct{})
// err := f(db, &daStructInstance)
func NewInserter(tableName string, exampleStruct interface{}) CRUDFunc {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build inserter: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	insertColumnNames := desc.Columns()
	insertColumnNamesStr := strings.Join(insertColumnNames, ", ")
	valueCount := len(insertColumnNames)

	insertStmt := fmt.Sprintf("insert into %s (%s) values (%s)",
		tableName, insertColumnNamesStr, markers(len(insertColumnNames)))

	return func(db dbHandle, item interface{}) error {

		itemValue := reflect.ValueOf(item)
		if itemValue.Type() != exampleItemType {
			return fmt.Errorf("crud.NewInserter func expected a %s, but got a %s",
				exampleItemType.String(), itemValue.Type().String())
		}

		insertValues := make([]interface{}, valueCount, valueCount)
		for i := range insertColumnNames {
			insertValues[i] = itemValue.Field(i).Interface()
		}

		_, err = db.Exec(insertStmt, insertValues...)

		return err
	}
}
