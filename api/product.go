package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type getProductRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type createProductRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	Thumb       string  `json:"thumb" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Gender      string  `json:"gender" binding:"required"`
	Material    string  `json:"material" binding:"required"`
	Size        string  `json:"size" binding:"required"`
	SizeOfModel string  `json:"size_of_model" binding:"required"`
}

type createProductResponse struct {
	ProductName string  `json:"product_name"`
	Thumb       string  `json:"thumb"`
	Price       float64 `json:"price"`
	Gender      string  `json:"gender"`
	Material    string  `json:"material"`
	Size        string  `json:"size"`
	SizeOfModel string  `json:"size_of_model"`
}

func newProductResponse(product db.Product, descriptionsProduct db.DescriptionsProduct) createProductResponse {
	return createProductResponse{
		ProductName: product.ProductName,
		Thumb:       product.Thumb,
		Price:       product.Price,
		Gender:      descriptionsProduct.Gender,
		Material:    descriptionsProduct.Material,
		Size:        descriptionsProduct.Size,
		SizeOfModel: descriptionsProduct.SizeOfModel,
	}
}

// @Summary Admin Create Product
// @ID adminCreateProduct
// @Produce json
// @Accept json
// @Param data body createProductRequest true "createProductRequest data"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {object} createProductResponse
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/products [post]
func (server *Server) adminCreateProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := uploader.UploadParams{
		Folder:         "ec-clothing",
		Format:         "jpg",
		Transformation: "f_auto,fl_lossy,q_auto:eco,dpr_auto,w_auto",
	}

	thumb, err := cld.Upload.Upload(ctx, req.Thumb, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argProduct := db.CreateProductParams{
		ProductName: req.ProductName,
		Thumb:       thumb.SecureURL,
		Price:       req.Price,
	}

	product, err := server.store.CreateProduct(ctx, argProduct)
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

	argDescriptionsProduct := db.CreateDescriptionProductParams{
		ProductID:   product.ID,
		Gender:      req.Gender,
		Material:    req.Material,
		Size:        req.Size,
		SizeOfModel: req.SizeOfModel,
	}

	descriptionsProduct, err := server.store.CreateDescriptionProduct(ctx, argDescriptionsProduct)
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

	rsp := newProductResponse(product, descriptionsProduct)
	ctx.JSON(http.StatusOK, rsp)
}

type getProductByNameResponse struct {
	Product []struct {
		db.Product             `json:"product"`
		db.DescriptionsProduct `json:"descriptions_product"`
		Stores                 []db.Store       `json:"stores"`
		Images                 []db.ImgsProduct `json:"images"`
		listFeedbackResponse   `json:"list_of_feedbacks"`
	} `json:"product"`
}

// @Summary Get Product By ID
// @ID getProductByID
// @Produce json
// @Accept json
// @Param data query listFeedbackRequest true "listFeedbackRequest data"
// @Param id path string true "ID"
// @Tags Started
// @Success 200 {object} getProductByNameResponse
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/products/{id} [get]
func (server *Server) getProductByID(ctx *gin.Context) {
	var result getProductByNameResponse
	var resultFeedbacks listFeedbackResponse
	var reqList listFeedbackRequest
	if err := ctx.ShouldBindQuery(&reqList); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqGet getProductRequest
	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, reqGet.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	descriptionsProduct, err := server.store.GetDescriptionProductByID(ctx, product.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listStore, err := server.store.ListStore(ctx, product.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listImage, err := server.store.ListImgProducts(ctx, product.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	feedbacks, err := server.store.ListFeedbacks(ctx, db.ListFeedbacksParams{
		ProductCommented: product.ID,
		Limit:            reqList.PageSize,
		Offset:           (reqList.PageID - 1) * reqList.PageSize,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, feedback := range feedbacks {
		user, err := server.store.GetUser(ctx, feedback.UserComment)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		userResult := newUserResponse(user)
		resultFeedbacks.Feedbacks = append(resultFeedbacks.Feedbacks, struct {
			db.Feedback `json:"feedback"`
			User        userResponse `json:"commentor"`
		}{feedback, userResult})
	}

	result.Product = append(result.Product, struct {
		db.Product             `json:"product"`
		db.DescriptionsProduct `json:"descriptions_product"`
		Stores                 []db.Store       `json:"stores"`
		Images                 []db.ImgsProduct `json:"images"`
		listFeedbackResponse   `json:"list_of_feedbacks"`
	}{product, descriptionsProduct, listStore, listImage, resultFeedbacks})

	ctx.JSON(http.StatusOK, result)
}

type listProductRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=10,max=20"`
}

// @Summary List Product
// @ID listProduct
// @Produce json
// @Accept json
// @Param data query listProductRequest true "listProductRequest data"
// @Tags Started
// @Success 200 {object} []db.Product
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/products [get]
func (server *Server) listProduct(ctx *gin.Context) {
	var reqList listProductRequest
	if err := ctx.ShouldBindQuery(&reqList); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductsParams{
		Limit:  reqList.PageSize,
		Offset: (reqList.PageID - 1) * reqList.PageSize,
	}

	products, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// @Summary Admin Update Product
// @ID adminUpdateProduct
// @Produce json
// @Accept json
// @Param data body createProductRequest true "createProductRequest data"
// @Param id path string true "ID"
// @Tags Admin
// @Security bearerAuth
// @Success 200 {object} createProductResponse
// @Failure 400 {string} error
// @Failure 403 {string} error
// @Failure 500 {string} error
// @Router /api/admin/products/{id} [put]
func (server *Server) adminUpdateProduct(ctx *gin.Context) {
	var reqGet getProductRequest
	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	params := uploader.UploadParams{
		Folder:         "ec-clothing",
		Format:         "jpg",
		Transformation: "f_auto,fl_lossy,q_auto:eco,dpr_auto,w_auto",
	}

	thumb, err := cld.Upload.Upload(ctx, req.Thumb, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argUpdateProduct := db.UpdateProductParams{
		ID:          reqGet.ID,
		ProductName: req.ProductName,
		Thumb:       thumb.SecureURL,
		Price:       req.Price,
	}

	updateProduct, err := server.store.UpdateProduct(ctx, argUpdateProduct)
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

	argUpdateDescriptionsProduct := db.UpdateDescriptionProductParams{
		ProductID:   updateProduct.ID,
		Gender:      req.Gender,
		Material:    req.Material,
		Size:        req.Size,
		SizeOfModel: req.SizeOfModel,
	}

	updateDescriptionsProduct, err := server.store.UpdateDescriptionProduct(ctx, argUpdateDescriptionsProduct)
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

	rsp := newProductResponse(updateProduct, updateDescriptionsProduct)
	ctx.JSON(http.StatusOK, rsp)
}