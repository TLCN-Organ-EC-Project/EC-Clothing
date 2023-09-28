package api

import (
	"database/sql"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Admin Create New Category
// @ID createCategory
// @Produce json
// @Accept json
// @Tags Admin
// @Security bearerAuth
// @Param data body createCategoryRequest true "createCategoryRequest data"
// @Security bearerAuth
// @Success 200 {object} db.Category
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/categories [post]
func (server *Server) adminCreateCategory(ctx *gin.Context) {
	var req createCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.CreateCategory(ctx, req.Name)
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

	ctx.JSON(http.StatusOK, category)
}

type getCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// @Summary Get Category
// @ID getCategory
// @Produce json
// @Accept json
// @Tags Started
// @Param id path string true "ID"
// @Success 200 {object} db.Category
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/categories/{id} [get]
func (server *Server) getCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, category)
}

// @Summary Get List Category
// @ID listCategory
// @Produce json
// @Accept json
// @Tags Started
// @Success 200 {array} []db.Category
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Router /api/categories [get]
func (server *Server) listCategory(ctx *gin.Context) {
	categories, err := server.store.ListCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var listCategory []db.Category
	listCategory = append(listCategory, categories...)

	ctx.JSON(http.StatusOK, listCategory)
}

type updateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Admin Update Category
// @ID updateCategory
// @Produce json
// @Accept json
// @Tags Admin
// @Security bearerAuth
// @Param data body updateCategoryRequest true "updateCategoryRequest data"
// @Param id path string true "ID"
// @Security bearerAuth
// @Success 200 {object} db.Category
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/categories/{id} [put]
func (server *Server) adminUpdateCategory(ctx *gin.Context) {
	var reqGet getCategoryRequest
	var reqUpdate updateCategoryRequest

	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&reqUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCategoryParams{
		ID:   reqGet.ID,
		Name: reqUpdate.Name,
	}

	categoryUpdate, err := server.store.UpdateCategory(ctx, arg)
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

	ctx.JSON(http.StatusOK, categoryUpdate)
}

// @Summary Admin Delete Category
// @ID adminDeleteCategory
// @Produce json
// @Accept json
// @Tags Admin
// @Param id path string true "ID"
// @Security bearerAuth
// @Success 200 {string} successfully
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/categories/{id} [delete]
func (server *Server) adminDeleteCategory(ctx *gin.Context) {
	var reqGet getCategoryRequest

	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategory(ctx, reqGet.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(ctx, category.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Delete Category Successfully")
}