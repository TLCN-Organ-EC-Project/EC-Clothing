package api

import (
	"log"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

	var imgProducts []db.ImgsProduct
	for _, image := range req.Images {
		img, err := cld.Upload.Upload(ctx, image, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg := db.CreateImgProductParams {
			ProductID: reqProduct.ID,
			Image: img.SecureURL,
		}

		imgProduct, err := server.store.CreateImgProduct(ctx, arg)
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

		imgProducts = append(imgProducts, imgProduct)
	}
		
	ctx.JSON(http.StatusOK, imgProducts)
}
