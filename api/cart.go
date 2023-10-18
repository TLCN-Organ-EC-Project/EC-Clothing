package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/XuanHieuHo/EC_Clothing/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createCartRequest struct {
	ProductID int64  `json:"product_id" binding:"required,min=1"`
	Size      string `json:"size" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required,min=1"`
}

// @Summary User Add Cart
// @ID createCart
// @Produce json
// @Accept json
// @Param data body createCartRequest true "createCartRequest data"
// @Param username path string true "Username"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} db.Cart
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/carts [post]
func (server *Server) createCart(ctx *gin.Context) {
	var reqUser getUserRequest
	var req createCartRequest

	if err := ctx.ShouldBindUri(&reqUser); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, reqUser.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("user doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	store, err := server.store.GetStore(ctx, db.GetStoreParams{
		ProductID: product.ID,
		Size: req.Size,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Quantity > int64(store.Quantity) {
		err := errors.New("quantity is not enough")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	existCart, err := server.store.GetCartDetails(ctx, db.GetCartDetailsParams{
		Username: reqUser.Username,
		ProductID: req.ProductID,
		Size: req.Size,
	})
	// nếu lỗi không có hàng nào thì tạo mới, nếu đã tổn tại thì update lại số lượng và tiền
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreateCartParams{
				Username:  reqUser.Username,
				ProductID: req.ProductID,
				Quantity:  int32(req.Quantity),
				Size:      req.Size,
				Price:     product.Price * float64(req.Quantity),
			}
		
			cart, err := server.store.CreateCart(ctx, arg)
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
			ctx.JSON(http.StatusOK, cart)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	quantity := existCart.Quantity + int32(req.Quantity)
	if quantity > store.Quantity {
		err := errors.New("quantity is not enough")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	updateCart, err := server.store.UpdateCart(ctx, db.UpdateCartParams{
		ID: existCart.ID,
		Quantity: quantity,
		Size: req.Size,
		Price: existCart.Price + product.Price * float64(req.Quantity),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, updateCart)
}

type listCartRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=10,max=20"`
}

type listCartResponse struct {
	Carts []struct {
		Cart db.Cart `json:"cart"`
		Product db.Product `json:"product"`
	} `json:"carts"`
}

// @Summary User Get List Cart
// @ID listCartOfUser
// @Produce json
// @Accept json
// @Param username path string true "Username"
// @Param data query listCartRequest true "listCartRequest data"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} listCartResponse
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/carts [get]
func (server *Server) listCartOfUser(ctx *gin.Context) {
	var req getUserRequest
	var reqList listCartRequest
	var result listCartResponse

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&reqList); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("user doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	listCarts, err := server.store.ListCartOfUser(ctx, db.ListCartOfUserParams{
		Username: user.Username,
		Limit:    reqList.PageSize,
		Offset:   (reqList.PageID - 1) * reqList.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	for _, cart := range listCarts {
		product, err := server.store.GetProduct(ctx, cart.ProductID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		result.Carts = append(result.Carts, struct {
			Cart db.Cart  `json:"cart"`
			Product db.Product `json:"product"`
		}{cart, product})
	}

	ctx.JSON(http.StatusOK, result)
}

type updateCartRequest struct {
	Size     string `json:"size" binding:"required"`
	Quantity int64  `json:"quantity" binding:"required,min=1"`
}

type getCartRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
	CartID   int64  `uri:"cart_id" binding:"required"`
}

// @Summary User Update Cart
// @ID updatCart
// @Produce json
// @Accept json
// @Param data body updateCartRequest true "updateCartRequest data"
// @Param username path string true "Username"
// @Param cart_id path string true "CartID"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} db.Cart
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/carts/{cart_id} [put]
func (server *Server) updateCart(ctx *gin.Context) {
	var reqGet getCartRequest
	var req updateCartRequest

	if err := ctx.ShouldBindUri(&reqGet); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, reqGet.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("user doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	cart, err := server.store.GetCart(ctx, reqGet.CartID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if cart.Username != authPayload.Username {
		err := errors.New("cart doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, cart.ProductID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateCartParams{
		ID:       cart.ID,
		Quantity: int32(req.Quantity),
		Size:     req.Size,
		Price:    product.Price * float64(req.Quantity),
	}

	updatedCart, err := server.store.UpdateCart(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, updatedCart)
}

// @Summary User Delete Cart
// @ID deleteCart
// @Produce json
// @Accept json
// @Tags User
// @Param username path string true "Username"
// @Param cart_id path string true "CartID"
// @Security bearerAuth
// @Success 200 {string} successfully
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/carts/{cart_id} [delete]
func (server *Server) deleteCart(ctx *gin.Context) {
	var req getCartRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("user doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	cart, err := server.store.GetCart(ctx, req.CartID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if cart.Username != authPayload.Username {
		err := errors.New("cart doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteCart(ctx, req.CartID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Delete Cart Successfully")
}

// @Summary User Delete Carts Of User
// @ID deleteCartOfUser
// @Produce json
// @Accept json
// @Tags User
// @Param username path string true "Username"
// @Security bearerAuth
// @Success 200 {string} successfully
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/carts [delete]
func (server *Server) deleteCartOfUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("user doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteCartOfUser(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Delete Cart Of User Successfully")
}