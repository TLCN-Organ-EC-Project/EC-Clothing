package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
