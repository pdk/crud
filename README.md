# crud

Database CRUD operations in go using reflection

Attempt to support all CRUD (insert, update, delete, query) operations with a
minimal amount of boilerplate code, either written or generated.

At system initialization, functions are created and stored in a `Machine`. The
functions may be used directly, or methods can be written to encapsulate the
needed type assertions.

Example:

    import (
        "database/sql"
        "github.com/pdk/crud"
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
    )

    func (f Foo) Insert(db *sql.DB) (Foo, error) {
        var err error
        f.ID, err = crudMachine.InsertGetID(db, f)
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
        results, err := crudMachine.Query(db, querySQL, queryParams...)
        return results.([]Foo), err
    }

    func QueryOneRowFoo(db *sql.DB, querySQL string, queryParams ...interface{}) (Foo, error) {
        result, err := crudMachine.QueryOneRow(db, querySQL, queryParams...)
        return result.(Foo), err
    }

Struct tags, ala `db:"bar"` are supported to specifying column names, if they
differ from field names.

Limitations: Only supports postgres style ("$1") bind markers.
