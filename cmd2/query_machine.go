package main

import (
	"database/sql"
	"log"

	"github.com/pdk/crud"
	q "github.com/pdk/query"
)

type QueryMachine struct {
	CrudMachine  crud.Machine
	QueryBuilder q.Builder
}

func NewQueryMachine(mach crud.Machine) QueryMachine {
	return QueryMachine{
		CrudMachine:  mach,
		QueryBuilder: q.Select(mach.Columns...).From(mach.TableName),
	}
}

func (qm QueryMachine) NewQuery(querySpecs q.Builder) func(*sql.DB, ...interface{}) (interface{}, error) {

	qry := qm.QueryBuilder.Merge(querySpecs)

	querySQL := qry.SQL()
	defaultArgs := qry.BindValues()

	return func(db *sql.DB, args ...interface{}) (interface{}, error) {
		return qm.CrudMachine.Query(db, querySQL, CombineArgs(defaultArgs, args)...)
	}
}

func CombineArgs(defaultArgs, newArgs []interface{}) []interface{} {

	if len(newArgs) > len(defaultArgs) {
		log.Fatalf("CombineArgs received more newArgs than defaultArgs. defaultArgs = %v, newArgs = %v",
			defaultArgs, newArgs)
	}

	args := make([]interface{}, 0, len(defaultArgs))
	args = append(args, defaultArgs[:len(defaultArgs)-len(newArgs)]...)
	args = append(args, newArgs...)

	return args
}

func NewQueryFactory(crud crud.Machine) func(q.Builder) func(*sql.DB, ...interface{}) (interface{}, error) {

	baseQuery := q.Select(crud.Columns...).From(crud.TableName)

	return func(querySpecs q.Builder) func(*sql.DB, ...interface{}) (interface{}, error) {

		theQuery := baseQuery.Merge(querySpecs)

		querySQL := theQuery.SQL()
		defaultArgs := theQuery.BindValues()

		return func(db *sql.DB, args ...interface{}) (interface{}, error) {
			return crud.Query(db, querySQL, CombineArgs(defaultArgs, args)...)
		}
	}
}
