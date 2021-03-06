package database

import (
	"testing"
)

func TestConnectSQLite(t *testing.T) {
	want := "sqlite"
	if db := Connect(); db.Name() != want {
		t.Errorf("wrong database: want %v, got %v", want, db.Name())
	}
}
