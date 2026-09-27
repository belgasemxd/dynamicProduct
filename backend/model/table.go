package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Table struct {
	name       string
	columns    []TableColumn
	rowsCol    *mongo.Collection
	tablesCol  *mongo.Collection
	tableID    primitive.ObjectID
	SelectedID primitive.ObjectID `bson:"_id,omitempty"`
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) TableExists() (bool, error) {
	var result struct {
		ID string `bson:"_id"`
	}
	err := t.tablesCol.FindOne(context.Background(), bson.M{"name": t.name}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewTable(client *mongo.Client, dbName string, tableName string) *Table {
	return &Table{
		name:       tableName,
		rowsCol:    client.Database(dbName).Collection("rows"),
		tablesCol:  client.Database(dbName).Collection("tables"),
		SelectedID: primitive.NilObjectID,
	}
}

func (t *Table) GetTableFromDB() error {
	var existing struct {
		ID      primitive.ObjectID `bson:"_id"`
		Columns []TableColumn      `bson:"columns"`
	}

	result := t.tablesCol.FindOne(context.Background(), bson.M{"name": t.name})
	err := result.Decode(&existing)
	if err != nil {
		return err
	}
	t.columns = existing.Columns
	t.tableID = existing.ID
	return nil
}

func (t *Table) CreateTableFromYAML() error {
	columns, err := t.loadColumnsFromYAML()
	if err != nil {
		return err
	}

	_, _ = t.tablesCol.DeleteMany(context.Background(), bson.M{"name": t.name})
	_, _ = t.rowsCol.DeleteMany(context.Background(), bson.M{"table_name": t.name})

	res, err := t.tablesCol.InsertOne(context.Background(), bson.M{"name": t.name, "columns": columns})
	if err != nil {
		return err
	}
	t.tableID = res.InsertedID.(primitive.ObjectID)
	t.columns = columns
	return nil
}

func (t *Table) SyncTableWithYAML() error {
	columns, err := t.loadColumnsFromYAML()
	if err != nil {
		return err
	}

	existingCols, err := t.GetColumns()
	if err != nil {
		return err
	}

	toAdd, toRemove := t.diffColumns(existingCols, columns)
	if len(toAdd) > 0 {
		if err := t.addColumns(toAdd); err != nil {
			return err
		}
	}
	if len(toRemove) > 0 {
		if err := t.removeColumns(toRemove); err != nil {
			return err
		}
	}

	return nil
}
