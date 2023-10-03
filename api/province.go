package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type getProvinceRequest struct {
	ID int64 `uri:"id" binding:"required,min=1,max=63"`
}

// @Summary Get Provinces By ID
// @ID getProvinceByID
// @Produce json
// @Accept json
// @Tags Started
// @Param id path string true "ID"
// @Success 200 {object} db.Province
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Router /api/provinces/{id} [get]
func (server *Server) getProvinceByID(ctx *gin.Context) {
	var req getProvinceRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	province, err := server.store.GetProvinceByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, province)
}

// @Summary Get List Provinces
// @ID listProvinces
// @Produce json
// @Accept json
// @Tags Started
// @Success 200 {array} []string
// @Failure 500 {string} error
// @Router /api/provinces [get]
func (server *Server) listProvinces(ctx *gin.Context) {

	provinces, err := server.store.ListProvinces(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, provinces)
}

type createProvinceRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createProvinces(ctx *gin.Context) {
	var req createProvinceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	province, err := server.store.CreateProvince(ctx, req.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, province)
}
