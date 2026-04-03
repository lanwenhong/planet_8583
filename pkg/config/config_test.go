package config

import (
	"os"
	"path"
	"testing"
)

func TestConfig(t *testing.T) {
	dir, _ := os.Getwd()
	ParseConfig(path.Join(dir, "planted.ini"))
}
