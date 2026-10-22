package config

import (
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestLoadAndDSN(t *testing.T) {
	t.Setenv("HOST", "")
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
	if config.Host != "127.0.0.1" || config.Port != "8080" || connection.Addr != "127.0.0.1:3306" || connection.DBName != "campus_message_board" || connection.Passwd != "special@:/?password" || !connection.ParseTime {
		t.Fatalf("incorrect connection configuration")
	}
}

func TestContainerAddresses(t *testing.T) {
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_NAME", "campus_message_board")
	t.Setenv("DB_USERNAME", "campus_app")
	t.Setenv("DB_PASSWORD", "container-test-password")
	settings, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	connection, err := mysql.ParseDSN(settings.DSN())
	if err != nil {
		t.Fatal(err)
	}
	if settings.Host != "0.0.0.0" || connection.Addr != "db:3306" {
		t.Fatal("容器应监听所有接口，并通过服务名连接 MySQL")
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
