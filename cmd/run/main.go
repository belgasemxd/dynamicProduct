package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/belgasemxd/dynamicProduct/backend/handlers"
	"github.com/belgasemxd/dynamicProduct/backend/model"
	"github.com/belgasemxd/dynamicProduct/resources"

	"github.com/belgasemxd/dynamicProduct/frontend"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/belgasemxd/dynamicProduct/util"
)

func main() {

	mongoURI, err := util.GetEnv()
	if err != nil {
		log.Fatal(err)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("2: ", err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("3: ", err)
	}

	tableName := "Equipment"
	dbName := "EquipmentDB"

	table := model.NewTable(client, dbName, tableName)
	err = table.GetTableFromDB()

	fmt.Println("Connected")

	router := gin.Default()

	router.StaticFS("/static", http.FS(resources.CssFS))

	router.HTMLRender = &frontend.TemplRender{}

	tableHandler := handlers.TableHandler{Table: table}

	tableHandler.Columns, err = tableHandler.Table.GetColumns()

	if err != nil {
		log.Fatal("3: ", err)
	}

	tableHandler.RowsMap, err = tableHandler.Table.GetRowsMap()
	if err != nil {
		log.Fatal("3: ", err)
	}

	router.GET("/", tableHandler.GetHome)
	router.GET("/reset-search", tableHandler.RestTablesHandler)
	router.GET("/reset-form", tableHandler.RestInputsHanlder)
	router.GET("/select/:id", tableHandler.SelectHandler)
	router.POST("/save", tableHandler.SaveHandler)
	router.DELETE("/delete", tableHandler.DeleteHandler)
	router.GET("/migrate", tableHandler.MigrateHandler)
	router.GET("/search/:text", tableHandler.SearchHandler)

	router.Run("localhost:8080")

}
