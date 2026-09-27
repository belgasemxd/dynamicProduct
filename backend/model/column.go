package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type TableColumn struct {
	Name         string      `yaml:"name" bson:"name"`
	Label        string      `yaml:"label" bson:"label"`
	Type         string      `yaml:"type" bson:"type"`
	Options      []string    `yaml:"options,omitempty" bson:"options,omitempty"`
	DefaultValue interface{} `yaml:"defaultValue" bson:"default_value"`
	Description  string      `yaml:"description" bson:"description"`
}

func (t *Table) RemoveColumn(col TableColumn) error {
	_, err := t.tablesCol.UpdateOne(context.Background(), bson.M{"_id": t.tableID}, bson.M{"$pull": bson.M{"columns": bson.M{"name": col.Name}}})
	if err != nil {
		return err
	}

	if err = t.RemoveColumnFromRows(col); err != nil {
		return err
	}

	_, err = t.GetColumns()
	return err
}

func (t *Table) RemoveColumnFromRows(col TableColumn) error {
	_, err := t.rowsCol.UpdateMany(context.Background(), bson.M{}, bson.M{"$unset": bson.M{"values." + col.Name: ""}})
	return err
}

func (t *Table) addColumns(cols []TableColumn) error {
	for _, c := range cols {
		if err := t.AddColumn(c); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) removeColumns(cols []TableColumn) error {
	for _, c := range cols {
		if err := t.RemoveColumn(c); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) AddColumn(col TableColumn) error {
	_, err := t.tablesCol.UpdateOne(context.Background(), bson.M{"_id": t.tableID}, bson.M{"$push": bson.M{"columns": col}})
	if err != nil {
		return err
	}
	t.columns = append(t.columns, col)

	_, err = t.rowsCol.UpdateMany(context.Background(), bson.M{}, bson.M{"$set": bson.M{"values." + col.Name: col.DefaultValue}})
	return err
}

func (t *Table) GetColumns() ([]TableColumn, error) {
	var table struct {
		Columns []TableColumn `bson:"columns"`
	}
	err := t.tablesCol.FindOne(context.Background(), bson.M{"_id": t.tableID}).Decode(&table)
	if err != nil {
		return nil, err
	}
	t.columns = table.Columns
	return t.columns, nil
}
