package describe_test

import (
	"testing"
	"time"

	"github.com/pdk/crud/describe"
)

// TestStruct provides a test structure
type TestStruct struct {
	Name        string
	Age         int64
	FavColor    string    `db:"favorite_color"`
	DateOfBirth string    `db:"dob" label:"Birth Date"`
	CreatedAt   time.Time `db:"created_at" label:"Created"`
	UpdatedAt   time.Time `db:"updated_at" label:"Updated"`
}

func TestDescribesStruct(t *testing.T) {

	d, err := describe.Describe(TestStruct{})
	if err != nil {
		t.Errorf("did not expect error, got %v", err)
	}

	if d.Name != "TestStruct" {
		t.Errorf("expected name to be TestStruct, got %s", d.Name)
	}

	if !sameStrings(d.Names(), []string{"Name", "Age", "FavColor", "DateOfBirth", "CreatedAt", "UpdatedAt"}) {
		t.Errorf("unexpected result for Names(): %v", d.Names())
	}

	if !sameStrings(d.Columns(), []string{"Name", "Age", "favorite_color", "dob", "created_at", "updated_at"}) {
		t.Errorf("unexpected result for Columns(): %v", d.Columns())
	}

	if !sameStrings(d.Labels(), []string{"Name", "Age", "FavColor", "Birth Date", "Created", "Updated"}) {
		t.Errorf("unexpected result for Labels(): %v", d.Labels())
	}

	if !sameStrings(d.Types(), []string{"string", "int64", "string", "string", "Time", "Time"}) {
		t.Errorf("unexpected result for Types(): %v", d.Types())
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
