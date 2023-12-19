package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/XuanHieuHo/EC_Clothing/token"
	"github.com/gin-gonic/gin"
)

type createOrderRequest struct {
	PromotionID   string   `json:"promotion_id"`
	Address       string   `json:"address" binding:"required"`
	Province      string   `json:"province" binding:"required"`
	PaymentMethod string   `json:"payment_method" binding:"required"`
	ProductID     []int64  `json:"product_id" binding:"required,min=1"`
	Size          []string `json:"size" binding:"required"`
	Quantity      []int64  `json:"quantity" binding:"required,min=1"`
}

// @Summary User Create Order
// @ID createOrder
// @Produce json
// @Accept json
// @Param data body createOrderRequest true "createOrderRequest data"
// @Param username path string true "UserOrder"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} db.OrderTxResult
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/orders [post]
func (server *Server) createOrder(ctx *gin.Context) {
	var reqUser getUserRequest
	var req createOrderRequest

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

	arg := db.OrderTxParams{
		Username:      user.Username,
		PromotionID:   req.PromotionID,
		Address:       req.Address,
		Province:      req.Province,
		PaymentMethod: req.PaymentMethod,
		ProductID:     req.ProductID,
		Size:          req.Size,
		Quantity:      req.Quantity,
	}

	result, err := server.store.OrderTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Summary User Get Order
// @ID getOrder
// @Produce json
// @Accept json
// @Param username path string true "Username"
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} db.Order
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/orders/{booking_id} [get]
func (server *Server) getOrder(ctx *gin.Context) {
	var req getOrderRequest

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

	order, err := server.store.GetOrder(ctx, req.BookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

type listOrderRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// @Summary User Get List Order
// @ID listOrderByUser
// @Produce json
// @Accept json
// @Param username path string true "Username"
// @Param data query listOrderRequest true "listOrderRequest data"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} []db.Order
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/orders [get]
func (server *Server) listOrderByUser(ctx *gin.Context) {
	var req getUserRequest
	var reqList listOrderRequest

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

	listOrders, err := server.store.ListOrderByUser(ctx, db.ListOrderByUserParams{
		UserBooking: user.Username,
		Limit:       reqList.PageSize,
		Offset:      (reqList.PageID - 1) * reqList.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, listOrders)
}

// @Summary Admin Get List Order
// @ID adminListOrder
// @Produce json
// @Accept json
// @Param data query listOrderRequest true "listOrderRequest data"
// @Security bearerAuth
// @Tags Admin
// @Success 200 {object} []db.Order
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/orders [get]
func (server *Server) adminListOrder(ctx *gin.Context) {
	var reqList listOrderRequest

	if err := ctx.ShouldBindQuery(&reqList); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	listOrders, err := server.store.ListOrder(ctx, db.ListOrderParams{
		Limit:  reqList.PageSize,
		Offset: (reqList.PageID - 1) * reqList.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, listOrders)
}

// @Summary Admin Get Order By Booking ID
// @ID adminGetOrderByBookingID
// @Produce json
// @Accept json
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags Admin
// @Success 200 {object} db.Order
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/orders/{booking_id} [get]
func (server *Server) adminGetOrderByBookingID(ctx *gin.Context) {
	var req adminGetOrderRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.BookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// @Summary Admin Get List Order By User
// @ID adminListOrderByUser
// @Produce json
// @Accept json
// @Param username path string true "Username"
// @Param data query listOrderRequest true "listOrderRequest data"
// @Security bearerAuth
// @Tags Admin
// @Success 200 {object} []db.Order
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/orders/users/{username} [get]
func (server *Server) adminListOrderByUser(ctx *gin.Context) {
	var req getUserRequest
	var reqList listOrderRequest

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

	listOrders, err := server.store.ListOrderByUser(ctx, db.ListOrderByUserParams{
		UserBooking: user.Username,
		Limit:       reqList.PageSize,
		Offset:      (reqList.PageID - 1) * reqList.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, listOrders)
}

type updateOrderRequest struct {
	Address   string   `json:"address" binding:"required"`
	Province  string   `json:"province" binding:"required"`
	ProductID []int64  `json:"product_id" binding:"required,min=1"`
	Size      []string `json:"size" binding:"required"`
	Quantity  []int64  `json:"quantity" binding:"required,min=1"`
}

// @Summary User Update Order
// @ID updateOrder
// @Produce json
// @Accept json
// @Param data body updateOrderRequest true "updateOrderRequest data"
// @Param username path string true "UserOrder"
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} db.UpdateOrderTxResult
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/orders/{booking_id} [put]
func (server *Server) updateOrder(ctx *gin.Context) {
	var reqGet getOrderRequest
	var req updateOrderRequest

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

	order, err := server.store.GetOrder(ctx, reqGet.BookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if order.UserBooking != authPayload.Username {
		err := errors.New("order doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateOrderTxParams{
		Username:  user.Username,
		Address:   req.Address,
		Province:  req.Province,
		ProductID: req.ProductID,
		Size:      req.Size,
		Quantity:  req.Quantity,
		BookingID: order.BookingID,
	}

	result, err := server.store.UpdateOrderTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Summary User Cancel Order
// @ID cancelOrder
// @Produce json
// @Accept json
// @Param username path string true "UserBooking"
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags User
// @Success 200 {string} successfully
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/orders/{booking_id}/cancel [put]
func (server *Server) cancelOrder(ctx *gin.Context) {
	var req getOrderRequest
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

	order, err := server.store.GetOrder(ctx, req.BookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if order.UserBooking != authPayload.Username {
		err := errors.New("order doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.CancelOrderParams {
		BookingID:       req.BookingID,
		UserBooking:     req.Username,
	}

	rsp, err := server.store.CancelOrderTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, rsp)
}

// @Summary User Confirm Order
// @ID confirmOrder
// @Produce json
// @Accept json
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags Admin
// @Success 200 {string} db.Order
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/orders/{booking_id}/confirm [put]
func (server *Server) confirmOrder(ctx *gin.Context) {
	var req adminGetOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.BookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateStatusOrderParams {
		BookingID: order.BookingID,
		Status: "confirm",
	}

	rsp, err := server.store.UpdateStatusOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, rsp)
}