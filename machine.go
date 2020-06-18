package crud

import (
	"database/sql"
)

// dbHandle is an interface where we can accept either a *sql.DB or a *sql.Trx
type dbHandle interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// CRUDFunc is a function that executes some function in a database for a
// particular value.
type CRUDFunc func(dbHandle, interface{}) error

// CRUDFuncGetID is a function that executes some function (insert) in a
// database, and returns an integer result (the newly generated ID value).
type CRUDFuncGetID func(dbHandle, interface{}) (int64, error)

// QueryFunc is a function that executes a query, with some bind parameters, and returns a value.
type QueryFunc func(dbHandle, string, ...interface{}) (interface{}, error)

// Machine is a holder of various CRUD functions. For Insert, InsertGetID,
// Update and Delete the first value returned will be the input value, and
// should be type-asserted back to original type. The InsertGetID func should
// update the input value with a new ID value from the database. Query will
// return a slice of the original type and should be converted to that.
// QueryOneRow will return a single instance of the struct type.
type Machine struct {
	Insert      CRUDFunc
	InsertGetID CRUDFuncGetID
	Update      CRUDFunc
	Delete      CRUDFunc
	Query       QueryFunc
	QueryOneRow QueryFunc
}

func NewMachine(tableName string, exampleStruct interface{}, keyFields ...string) Machine {

	return Machine{
		Insert:      NewInserter(tableName, exampleStruct),
		Update:      NewUpdater(tableName, exampleStruct, keyFields...),
		Delete:      NewDeleter(tableName, exampleStruct, keyFields...),
		Query:       NewQueryFunc(tableName, exampleStruct),
		QueryOneRow: NewQueryOneRowFunc(tableName, exampleStruct),
	}
}

func NewMachineGetID(tableName string, exampleStruct interface{}, idField string) Machine {

	return Machine{
		InsertGetID: NewAutoIncrIDInserter(tableName, exampleStruct, idField),
		Update:      NewUpdater(tableName, exampleStruct, idField),
		Delete:      NewDeleter(tableName, exampleStruct, idField),
		Query:       NewQueryFunc(tableName, exampleStruct),
		QueryOneRow: NewQueryOneRowFunc(tableName, exampleStruct),
	}
}
