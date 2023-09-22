package api

import (
	"database/sql"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createProductInCategoryRequest struct {
	ProductID []int64 `json:"product_id" binding:"required,min=1"`
}

// @Summary Admin Create Products In Category
// @ID adminCreateProductInCategory
// @Produce json
// @Accept json
// @Param data body createProductInCategoryRequest true "createProductInCategoryRequest data"
// @Param id path string true "ID"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {object} db.ProductsInCategory
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/categories/{id}/products [post]
func (server *Server) adminCreateProductInCategory(ctx *gin.Context) {
	var reqGet getCategoryRequest
	var reqAdd createProductInCategoryRequest

	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&reqAdd); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var productsInCategory []db.ProductsInCategory

	for _, product := range reqAdd.ProductID {
		arg := db.CreateProductsInCategoryParams {
			ProductID: product,
			CategoryID: reqGet.ID,
		}
		productInCate, err := server.store.CreateProductsInCategory(ctx, arg)
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
		productsInCategory = append(productsInCategory, productInCate)
	}

	ctx.JSON(http.StatusOK, productsInCategory)
}

type getProductAndCategoryRequest struct {
	CategoryID int64 `uri:"id" binding:"required,min=1"`
	ProductID  int64 `uri:"product_id" binding:"required,min=1"`
}

// @Summary Admin Delete Products In Category
// @ID adminDeleteProductInCategory
// @Produce json
// @Accept json
// @Param id path string true "CategoryID"
// @Param product_id path string true "ProductID"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {string} Delete Product In Category Successfully
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/categories/{id}/products/{product_id} [delete]
func (server *Server) adminDeleteProductInCategory(ctx *gin.Context) {
	var req getProductAndCategoryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetProductsInCategoryByIDParams{
		ProductID:  req.ProductID,
		CategoryID: req.CategoryID,
	}

	productInCategory, err := server.store.GetProductsInCategoryByID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteProductsInCategory(ctx, productInCategory.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Delete Product In Category Successfully")
}

// @Summary Get Products In Category
// @ID getProductsInCategory
// @Produce json
// @Accept json
// @Param data query listProductRequest true "listProductRequest data"
// @Param id path string true "CategoryID"
// @Tags Started
// @Success 200 {object} []db.Product
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Router /api/categories/{id}/products [get]
func (server *Server) getProductsInCategory(ctx *gin.Context) {
	var req getCategoryRequest
	var reqList listProductRequest
	var result []db.Product

	if err := ctx.ShouldBindQuery(&reqList); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductsInCategoryParams {
		CategoryID: req.ID,
		Limit:  reqList.PageSize,
		Offset: (reqList.PageID - 1) * reqList.PageSize,
	}

	productsInCategory, err := server.store.ListProductsInCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, productInCategory := range productsInCategory {
		product, err := server.store.GetProduct(ctx, productInCategory.ProductID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		result = append(result, product)
	}
	ctx.JSON(http.StatusOK, result)
}