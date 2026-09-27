package model

import (
	"fmt"

	"github.com/belgasemxd/dynamicProduct/util"
	"gopkg.in/yaml.v3"
)

func (t *Table) loadColumnsFromYAML() ([]TableColumn, error) {
	data, err := util.ReadYAMLFile()
	if err != nil {
		return nil, err
	}
	var parsed map[string][]TableColumn
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	cols, ok := parsed[t.Name()]
	if !ok {
		return nil, fmt.Errorf("Tabelle %s nicht gefunden", t.Name())
	}
	return cols, nil
}

func (t *Table) diffColumns(existing, yamlCols []TableColumn) (toAdd, toRemove []TableColumn) {
	existingMap := make(map[string]TableColumn)
	for _, c := range existing {
		existingMap[c.Name] = c
	}
	yamlMap := make(map[string]TableColumn)
	for _, c := range yamlCols {
		yamlMap[c.Name] = c
	}

	for _, c := range yamlCols {
		if _, ok := existingMap[c.Name]; !ok {
			toAdd = append(toAdd, c)
		}
	}
	for _, c := range existing {
		if _, ok := yamlMap[c.Name]; !ok {
			toRemove = append(toRemove, c)
		}
	}
	return
}
