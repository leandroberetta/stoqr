package database

import (
	"os"
	"testing"
)

func TestCheckPostgresConfigMissingHost(t *testing.T) {
	want := "postgres host is empty"
	if err := checkPostgresConfig(); err != nil && err.Error() != want {
		t.Errorf("wrong error: want %v, got %v", want, err)
	}
}

func TestCheckPostgresConfigMissingPort(t *testing.T) {
	os.Setenv("STOQR_API_DB_HOST", "localhost")

	want := "postgres port is empty"
	if err := checkPostgresConfig(); err != nil && err.Error() != want {
		t.Errorf("wrong error: want %v, got %v", want, err)
	}
}

func TestCheckPostgresConfigMissingUser(t *testing.T) {
	os.Setenv("STOQR_API_DB_HOST", "localhost")
	os.Setenv("STOQR_API_DB_PORT", "5432")

	want := "postgres user is empty"
	if err := checkPostgresConfig(); err != nil && err.Error() != want {
		t.Errorf("wrong error: want %v, got %v", want, err)
	}
}

func TestCheckPostgresConfigMissingPassword(t *testing.T) {
	os.Setenv("STOQR_API_DB_HOST", "localhost")
	os.Setenv("STOQR_API_DB_PORT", "5432")
	os.Setenv("STOQR_API_DB_USER", "postgres")

	want := "postgres password is empty"
	if err := checkPostgresConfig(); err != nil && err.Error() != want {
		t.Errorf("wrong error: want %v, got %v", want, err)
	}
}

func TestCheckPostgresConfigMissingName(t *testing.T) {
	os.Setenv("STOQR_API_DB_HOST", "localhost")
	os.Setenv("STOQR_API_DB_PORT", "5432")
	os.Setenv("STOQR_API_DB_USER", "postgres")
	os.Setenv("STOQR_API_DB_PASSWORD", "postgres")

	want := "postgres db name is empty"
	if err := checkPostgresConfig(); err != nil && err.Error() != want {
		t.Errorf("wrong error: want %v, got %v", want, err)
	}
}

func TestCheckPostgresConfigOK(t *testing.T) {
	os.Setenv("STOQR_API_DB_HOST", "localhost")
	os.Setenv("STOQR_API_DB_PORT", "5432")
	os.Setenv("STOQR_API_DB_USER", "postgres")
	os.Setenv("STOQR_API_DB_PASSWORD", "postgres")
	os.Setenv("STOQR_API_DB_NAME", "postgres")

	if err := checkPostgresConfig(); err != nil {
		t.Errorf("wrong error: want %v, got %v", nil, err)
	}
}

func TestGeneratePostgresConnectionString(t *testing.T) {
	os.Setenv("STOQR_API_DB_HOST", "localhost")
	os.Setenv("STOQR_API_DB_PORT", "5432")
	os.Setenv("STOQR_API_DB_USER", "postgres")
	os.Setenv("STOQR_API_DB_PASSWORD", "postgres")
	os.Setenv("STOQR_API_DB_NAME", "postgres")

	want := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	if got := generatePostgresConnectionString(); got != want {
		t.Errorf("wrong connection string: want %v, got %v", want, got)
	}
}
