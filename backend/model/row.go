package model

import (
	"context"
	"fmt"

	"github.com/belgasemxd/dynamicProduct/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableRow struct {
	ID     primitive.ObjectID     `bson:"_id,omitempty"`
	Values map[string]interface{} `bson:"values"`
}

func (t *Table) GetRowByID(rowID primitive.ObjectID) (*TableRow, error) {
	rows, err := t.GetRows()
	if err != nil {
		return nil, err
	}

	for i := range rows {
		if rows[i].ID == rowID {
			return &rows[i], nil
		}
	}

	return nil, nil
}

func (t *Table) GetRowByIDAsMap(rowID primitive.ObjectID) (map[string]string, error) {
	row, err := t.GetRowByID(rowID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}

	rowMap := make(map[string]string)
	rowMap["id"] = row.ID.Hex()
	for _, col := range t.columns {
		val, ok := row.Values[col.Name]
		if !ok || val == nil {
			rowMap[col.Name] = ""
			continue
		}
		rowMap[col.Name] = fmt.Sprintf("%v", val)
	}

	return rowMap, nil
}

func (t *Table) AddRow(values map[string]interface{}) (*primitive.ObjectID, error) {
	for _, col := range t.columns {
		val, ok := values[col.Name]
		if !ok {
			val = col.DefaultValue
		}

		if !util.IsValidType(val, col.Type) {
			continue
		}

		values[col.Name] = val
	}

	res, err := t.rowsCol.InsertOne(context.Background(), TableRow{Values: values})
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (t *Table) RemoveRow(rowID primitive.ObjectID) error {
	_, err := t.rowsCol.DeleteOne(context.Background(), bson.M{"_id": rowID})
	return err
}

func (t *Table) GetRows() ([]TableRow, error) {
	cur, err := t.rowsCol.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var rows []TableRow
	if err := cur.All(context.Background(), &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (t *Table) GetRowsMap() ([]map[string]string, error) {
	rows, err := t.GetRows()
	if err != nil {
		return nil, err
	}

	var result []map[string]string
	for _, row := range rows {
		rowMap := make(map[string]string)
		rowMap["id"] = row.ID.Hex()
		for _, col := range t.columns {
			val, ok := row.Values[col.Name]
			if !ok || val == nil {
				rowMap[col.Name] = ""
				continue
			}

			rowMap[col.Name] = fmt.Sprintf("%v", val)
		}
		result = append(result, rowMap)
	}

	return result, nil
}
