package api

import (
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/gin-gonic/gin"
)

type addImageProductRequest struct {
	Images []string `json:"images" binding:"required"`
}

// @Summary Admin Add Image of Product
// @ID adminAddImageProduct
// @Produce json
// @Accept json
// @Param data body addImageProductRequest true "addImageProductRequest data"
// @Param id path string true "ID"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {object} []db.ImgsProduct
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/products/{id} [post]
func (server *Server) adminAddImageProduct(ctx *gin.Context) {
	var reqProduct getProductRequest
	if err := ctx.ShouldBindUri(&reqProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addImageProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddImageProductTxParams{
		Images: req.Images,
		ID:     reqProduct.ID,
	}

	rsp, err := server.store.CreateImgProductTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
