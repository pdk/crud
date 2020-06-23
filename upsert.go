package crud

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pdk/crud/describe"
)

func NewUpserter(bindStyle BindStyle, tableName string, exampleStruct interface{}, keyFields ...string) CRUDFunc {

	desc, err := describe.Describe(exampleStruct)
	if err != nil {
		log.Fatalf("cannot build updateer: %v", err)
	}

	exampleItemType := reflect.ValueOf(exampleStruct).Type()

	insertColumnNames := desc.Columns()
	insertColumnNamesStr := strings.Join(insertColumnNames, ", ")

	keyColumnNames := desc.ColumnsOf(keyFields...)
	keyColumnNamesStr := strings.Join(keyColumnNames, ", ")

	updateColumnNames := desc.ColumnsOmitFields(keyFields...)

	onConflictClause := "do nothing"
	if len(updateColumnNames) > 0 {
		onConflictClause = "do update set"
		for i, c := range updateColumnNames {
			if i > 0 {
				onConflictClause += ","
			}
			onConflictClause += " " + c + " = excluded." + c
		}
	}

	upsertStmt := fmt.Sprintf("insert into %s (%s) values (%s) on conflict (%s) %s",
		tableName, insertColumnNamesStr, markers(bindStyle, 1, len(insertColumnNames)),
		keyColumnNamesStr, onConflictClause)

	log.Printf("upsertStmt: %s", upsertStmt)

	valueCount := len(insertColumnNames)

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

		_, err = db.Exec(upsertStmt, insertValues...)

		return err
	}
}
