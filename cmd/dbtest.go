package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pdk/crud"

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
	crudMachine = crud.NewMachineGetID("foo", Foo{}, "ID")
	inserter2   = crud.NewInserter("foo", Foo{})
)

func (f Foo) Insert(db *sql.DB) (Foo, error) {
	var err error
	f.ID, err = crudMachine.InsertGetID(db, f)
	return f, err
}

func (f Foo) Insert2(db *sql.DB) (Foo, error) {
	err := inserter2(db, f)
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

	f, err = f.Delete(db)
	if err != nil {
		log.Printf("delete failed: %v", err)
	}

	log.Printf("deleted")
}
