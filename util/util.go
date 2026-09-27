package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetEnv() (string, error) {
	if err := godotenv.Load("resources/config/.env"); err != nil {
		return "", fmt.Errorf("fehler beim Laden der .env Datei: %w", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return "", fmt.Errorf("MONGODB_URI nicht gesetzt")
	}

	return mongoURI, nil
}

func ReadYAMLFile() ([]byte, error) {
	const yamlPath = "resources/config/model.yaml"

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Lesen der YAML-Datei %s: %w", yamlPath, err)
	}

	return data, nil
}

func IsValidType(val interface{}, colType string) bool {
	switch colType {
	case "text", "string":
		_, ok := val.(string)
		return ok
	case "number", "int", "float":
		switch val.(type) {
		case int, int32, int64, float32, float64:
			return true
		}
	case "bool", "boolean":
		_, ok := val.(bool)
		return ok
	}
	return false
}

func GetCSSPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	root := filepath.Dir(wd)

	cssPath := filepath.Join(root, "resources", "css")
	return cssPath, nil
}
