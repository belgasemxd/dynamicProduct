package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *Table) UpdateRow(rowID primitive.ObjectID, updates map[string]interface{}) error {
	columns, err := t.GetColumns()
	if err != nil {
		return err
	}

	set := make(map[string]interface{})
	for _, col := range columns {
		if val, ok := updates[col.Name]; ok {
			set["values."+col.Name] = val
		}
	}

	if len(set) == 0 {
		return nil
	}

	filter := map[string]interface{}{"_id": rowID}
	update := map[string]interface{}{"$set": set}

	_, err = t.rowsCol.UpdateOne(context.Background(), filter, update)
	return err
}
