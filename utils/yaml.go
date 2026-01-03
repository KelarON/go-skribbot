package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// WriteStructToYAMLFile writes a struct to a YAML file.
// The file is created if it does not exist, and overwritten if it does.
func WriteStructToYAMLFile(filename string, data interface{}) error {

	yamlAdditionalData := []byte("# drawing_type: 0 - lines (faster, but more errors), 1 - points (slower, but better quality)")

	// Marshal the struct to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %v", err)
	}

	// Write the YAML data to a file
	err = os.WriteFile(filename, append(yamlData, yamlAdditionalData...), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
