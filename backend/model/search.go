package model

import (
	"fmt"
	"strings"
)

func (t *Table) SearchMap(query string) ([]map[string]string, error) {
	rows, err := t.Search(query)
	if err != nil {
		return nil, err
	}

	var results []map[string]string
	for _, row := range rows {
		m := make(map[string]string)
		m["id"] = row.ID.Hex() // falls du wie bei GetRowsMap auch die ID speichern willst
		for _, col := range t.columns {
			val, ok := row.Values[col.Name]
			if !ok || val == nil {
				m[col.Name] = ""
				continue
			}
			m[col.Name] = fmt.Sprintf("%v", val)
		}
		results = append(results, m)
	}

	return results, nil
}

func (t *Table) Search(query string) ([]TableRow, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	if strings.Contains(query, ":") {
		return t.searchStructured(query)
	}

	return t.searchFreeText(query)
}

func (t *Table) searchFreeText(q string) ([]TableRow, error) {
	qLower := strings.ToLower(q)

	rows, err := t.GetRows()
	if err != nil {
		return nil, err
	}

	var results []TableRow
	for _, row := range rows {
		for _, val := range row.Values {
			valStr := fmt.Sprintf("%v", val)
			if strings.Contains(strings.ToLower(valStr), qLower) {
				results = append(results, row)
				break
			}
		}
	}
	return results, nil
}

func (t *Table) searchStructured(q string) ([]TableRow, error) {
	rows, err := t.GetRows()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(q, ",")
	conditions := make(map[string]string)
	for _, p := range parts {
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(kv[0]))
		val := strings.TrimSpace(strings.ToLower(kv[1]))
		conditions[key] = val
	}

	var results []TableRow
RowLoop:
	for _, row := range rows {
		for key, cond := range conditions {
			v, ok := row.Values[key]
			if !ok {
				continue RowLoop
			}
			valStr := strings.ToLower(fmt.Sprintf("%v", v))
			if !strings.Contains(valStr, cond) {
				continue RowLoop
			}
		}
		results = append(results, row)
	}
	return results, nil
}
