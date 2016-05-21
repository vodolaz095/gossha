package config

import (
	"fmt"
	"os"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	hmdr2, err := GetHomeDir()
	if err != nil {
		t.Errorf("Error getting home directory %s", err)
	}
	fmt.Printf("Homedir path found at %s\n", hmdr2)
}

func TestGetDatabasePath(t *testing.T) {
	dbpath, err := GetDatabasePath()
	if err != nil {
		t.Errorf("Error getting database path")
	}
	fmt.Printf("Database path found at %s\n", dbpath)
}

func TestGetPrivateKeyPath(t *testing.T) {
	dbpath, err := GetPrivateKeyPath()
	if err != nil {
		t.Errorf("Error getting database path")
	}
	fmt.Printf("Private key path found at %s\n", dbpath)
}

func TestGetPublicKeyPath(t *testing.T) {
	dbpath, err := GetPublicKeyPath()
	if err != nil {
		t.Errorf("Error getting database path")
	}
	fmt.Printf("Public key path found at %s\n", dbpath)
}

func TestInitConfig(t *testing.T) {
	os.Setenv("HOME", "/tmp/gossha_test/")
	cfg, err := InitConfig()
	if err != nil {
		t.Errorf("Error initializing config - %s", err)
	}

	fmt.Println(cfg)

	if cfg.Driver != "sqlite3" {
		t.Error("Wrong driver being used")
	}
	if !cfg.Debug {
		t.Error("Non debug mode!")
	}
	if cfg.ExecuteOnMessage == "" {
		t.Error("ExecuteOnMessage not empty")
	}
	if cfg.ExecuteOnPrivateMessage == "" {
		t.Error("ExecuteOnPrivateMessage not empty")
	}

}
