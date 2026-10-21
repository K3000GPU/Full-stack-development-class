package config

import (
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestLoadAndDSN(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "campus_message_board")
	t.Setenv("DB_USERNAME", "test_user")
	t.Setenv("DB_PASSWORD", "special@:/?password")
	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	connection, err := mysql.ParseDSN(config.DSN())
	if err != nil {
		t.Fatal(err)
	}
	if config.Port != "8080" || connection.Addr != "127.0.0.1:3306" || connection.DBName != "campus_message_board" || connection.Passwd != "special@:/?password" || !connection.ParseTime {
		t.Fatalf("incorrect connection configuration")
	}
}

func TestRejectsMissingCredentialsAndInvalidPorts(t *testing.T) {
	for _, name := range []string{"DB_NAME", "DB_USERNAME", "DB_PASSWORD", "PORT", "DB_PORT"} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("DB_NAME", "campus_message_board")
			t.Setenv("DB_USERNAME", "test_user")
			t.Setenv("DB_PASSWORD", "test_password")
			t.Setenv("PORT", "8080")
			t.Setenv("DB_PORT", "3306")
			value := ""
			if name == "PORT" || name == "DB_PORT" {
				value = "70000"
			}
			t.Setenv(name, value)
			if _, err := Load(); err == nil {
				t.Fatal("expected configuration error")
			}
		})
	}
}
