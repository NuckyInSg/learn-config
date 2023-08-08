package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Can not get pwd")
	}
	tempFileName := pwd + "/config.yaml"
	fmt.Println(tempFileName)
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

	err = os.WriteFile(tempFileName, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Cannot create temporary file: %s", err)
	}
	defer os.Remove(tempFileName)

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
