package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tempFileName := os.TempDir() + "config.yaml"
	configContent := `
  http:
    gin_mode: debug
    port: 8080
  database:
    driver: mysql
    username: root
    address: localhost
    database_name: test
  log:
    filename: /var/log/yourname.log 
  `

	err := os.WriteFile(tempFileName, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Cannot create temporary file: %s", err)
	}

	t.Run("NewConfig", func(t *testing.T) {
		oldValue := os.Getenv(configPathEnvVarName)
		os.Setenv(configPathEnvVarName, tempFileName)
		defer os.Setenv(configPathEnvVarName, oldValue)

		config, err := newConfig()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		fmt.Println(config.Database)
	})
}
