package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/pdk/crud"
	q "github.com/pdk/query"

	_ "github.com/lib/pq"
)

type Foo struct {
	ID      int64
	Name    string
	Age     int
	Address string
	Salary  float32
}

var (
	crudMachine  = crud.NewMachineGetID(crud.DollarNumber, "foo", Foo{}, "ID")
	crudMachine2 = crud.NewMachine(crud.DollarNumber, "foo", Foo{}, "ID")

	// queryMachine = NewQueryMachine(crudMachine)
	queryMachine = NewQueryFactory(crudMachine)

	FooByAge = NewFooQuery(q.Where("age =", q.Bind(nil)))
)

func NewFooQuery(qry q.Builder) func(*sql.DB, ...interface{}) ([]Foo, error) {

	//queryFunc := queryMachine.NewQuery(qry)
	queryFunc := queryMachine(qry)

	return func(db *sql.DB, args ...interface{}) ([]Foo, error) {
		return typeFix(queryFunc(db, args...))
	}
}

func (f Foo) Insert(db *sql.DB) (Foo, error) {
	var err error
	f.ID, err = crudMachine.InsertGetID(db, f)
	return f, err
}

func (f Foo) Insert2(db *sql.DB) (Foo, error) {
	err := crudMachine2.Insert(db, f)
	return f, err
}

func (f Foo) Update(db *sql.DB) (Foo, error) {
	err := crudMachine.Update(db, f)
	return f, err
}

func (f Foo) Delete(db *sql.DB) (Foo, error) {
	err := crudMachine.Delete(db, f)
	return f, err
}

func QueryFoo(db *sql.DB, querySQL string, queryParams ...interface{}) ([]Foo, error) {
	return typeFix(crudMachine.Query(db, querySQL, queryParams...))
}

func QueryOneRowFoo(db *sql.DB, querySQL string, queryParams ...interface{}) (Foo, error) {
	return typeFixOne(crudMachine.QueryOneRow(db, querySQL, queryParams...))
}

func typeFix(data interface{}, err error) ([]Foo, error) {
	return data.([]Foo), err
}

func typeFixOne(data interface{}, err error) (Foo, error) {
	return data.(Foo), err
}

func main() {

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("drop table if exists foo")
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	_, err = db.Exec("create table foo (id bigserial primary key, name text, age int, address varchar(200), salary numeric(12,2))")
	if err != nil {
		log.Fatalf("cannot create table: %v", err)
	}

	log.Printf("table created!")

	f := Foo{
		Name:    "Walter",
		Age:     32,
		Address: "123 Main St",
		Salary:  1234.56,
	}

	f, err = f.Insert(db)
	if err != nil {
		log.Fatalf("failed to insert: %v", err)
	}

	log.Printf("rec inserted, new id = %d", f.ID)

	f.ID = 2
	f.Name = "Marty"
	f, err = f.Insert2(db)
	if err != nil {
		log.Printf("failed to insert again: %v", err)
	}

	log.Printf("2nd rec inserted")

	f.Name = "Myrtle"

	f, err = f.Update(db)
	if err != nil {
		log.Printf("failed to update: %v", err)
	}

	log.Printf("final: %v", f)

	log.Printf("select columns are: %v", crudMachine.Columns)

	row, err := QueryOneRowFoo(db, "select "+strings.Join(crudMachine.Columns, ", ")+
		" from foo where id = $1", 1)
	if err != nil {
		log.Printf("query one row failed: %v", err)
	}

	log.Printf("query one row: %v", row)

	rows, err := QueryFoo(db, "select "+strings.Join(crudMachine.Columns, ", ")+" from foo")
	if err != nil {
		log.Printf("query failed: %v", err)
	}

	for _, r := range rows {
		log.Printf("query result = %v", r)
	}

	rows, err = FooByAge(db, 32)
	if err != nil {
		log.Printf("query by age failed: %v", err)
	}

	for _, r := range rows {
		log.Printf("query by age result = %v", r)
	}

	f, err = f.Delete(db)
	if err != nil {
		log.Printf("delete failed: %v", err)
	}

	log.Printf("deleted")
}
