package handlers

import (
	"net/http"

	"github.com/belgasemxd/dynamicProduct/backend/model"
	"github.com/belgasemxd/dynamicProduct/frontend"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableHandler struct {
	Table   *model.Table
	RowsMap []map[string]string
	Columns []model.TableColumn
}

func (h *TableHandler) GetHome(ctx *gin.Context) {
	var err error = nil

	h.Table.SelectedID = primitive.NilObjectID

	h.RowsMap, err = h.Table.GetRowsMap()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	h.Columns, err = h.Table.GetColumns()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	h.RestAll(ctx)
}

func (h *TableHandler) RestAll(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", frontend.Page(h.Columns, h.RowsMap, h.Table.Name(), h.Table.SelectedID.Hex()))
}

func (h *TableHandler) RestContent(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", templ.Join(frontend.Table(h.Columns, h.RowsMap, h.Table.SelectedID.Hex()), frontend.Inputs(h.Columns)))
}

func (h *TableHandler) SelectHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	h.Table.SelectedID = objectID

	rowMap, err := h.Table.GetRowByIDAsMap(objectID)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "", templ.Join(frontend.Table(h.Columns, h.RowsMap, h.Table.SelectedID.Hex()), frontend.InputsWithValues(h.Columns, rowMap)))
}

func (h *TableHandler) SearchHandler(ctx *gin.Context) {
	var err error = nil

	text := ctx.Param("text")

	h.RowsMap, err = h.Table.SearchMap(text)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	h.Table.SelectedID = primitive.NilObjectID
	h.RestContent(ctx)
}

func (h *TableHandler) RestInputsHanlder(ctx *gin.Context) {
	h.Table.SelectedID = primitive.NilObjectID
	h.RestContent(ctx)
}

func (h *TableHandler) RestTablesHandler(ctx *gin.Context) {
	h.Table.SelectedID = primitive.NilObjectID
	var err error = nil
	h.RowsMap, err = h.Table.GetRowsMap()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	h.RestContent(ctx)
}

func (h *TableHandler) SaveHandler(ctx *gin.Context) {
	data := make(map[string]interface{})

	if err := ctx.Request.ParseForm(); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	for key, values := range ctx.Request.PostForm {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	if h.Table.SelectedID == primitive.NilObjectID {
		_, err := h.Table.AddRow(data)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}
	} else {
		err := h.Table.UpdateRow(h.Table.SelectedID, data)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}
	}
	var err error = nil

	h.RowsMap, err = h.Table.GetRowsMap()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	h.Table.SelectedID = primitive.NilObjectID
	h.RestContent(ctx)

}

func (h *TableHandler) DeleteHandler(ctx *gin.Context) {
	if h.Table.SelectedID != primitive.NilObjectID {
		err := h.Table.RemoveRow(h.Table.SelectedID)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}
	}
	var err error = nil

	h.RowsMap, err = h.Table.GetRowsMap()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	h.Table.SelectedID = primitive.NilObjectID
	h.RestContent(ctx)

}

func (h *TableHandler) MigrateHandler(ctx *gin.Context) {
	if err := h.Table.SyncTableWithYAML(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	var err error = nil

	h.RowsMap, err = h.Table.GetRowsMap()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	h.Table.SelectedID = primitive.NilObjectID

	ctx.Redirect(http.StatusSeeOther, "/")

}
