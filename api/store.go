package api

import (
	"fmt"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type addProductToStoreRequest struct {
	Size     []string `json:"size" binding:"required"`
	Quantity []int32  `json:"quantity" binding:"required,min=1"`
}

// @Summary Admin Add Product To Store
// @ID adminAddProductToStore
// @Produce json
// @Accept json
// @Param data body addProductToStoreRequest true "addProductToStoreRequest data"
// @Param id path string true "ID"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {object} []db.Store
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/products/{id}/store [post]
func (server *Server) adminAddProductToStore(ctx *gin.Context) {
	var reqGet getProductRequest
	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addProductToStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if len(req.Quantity) != len(req.Size) {
		err := fmt.Errorf("retype full 2 properties")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var storeProduct []db.Store
	for i := range req.Size {
		size := req.Size[i]
		quantity := req.Quantity[i]

		arg := db.CreateStoreParams{
			ProductID: reqGet.ID,
			Size:      size,
			Quantity:  quantity,
		}

		store, err := server.store.CreateStore(ctx, arg)
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

		storeProduct = append(storeProduct, store)
	}

	ctx.JSON(http.StatusOK, storeProduct)
}

// @Summary Admin Update Product To Store
// @ID adminUpdateProductToStore
// @Produce json
// @Accept json
// @Param data body addProductToStoreRequest true "addProductToStoreRequest data"
// @Param id path string true "ID"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {object} []db.Store
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/products/{id}/store [put]
func (server *Server) adminUpdateProductToStore(ctx *gin.Context) {
	var reqGet getProductRequest
	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req addProductToStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var storeProduct []db.Store
	for i := range req.Size {
		size := req.Size[i]
		quantity := req.Quantity[i]

		arg := db.UpdateStoreParams{
			ProductID: reqGet.ID,
			Size:      size,
			Quantity:  quantity,
		}

		store, err := server.store.UpdateStore(ctx, arg)
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

		storeProduct = append(storeProduct, store)
	}

	ctx.JSON(http.StatusOK, storeProduct)
}
