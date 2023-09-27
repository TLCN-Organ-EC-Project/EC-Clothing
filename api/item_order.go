package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/XuanHieuHo/EC_Clothing/token"
	"github.com/gin-gonic/gin"
)

type getOrderRequest struct {
	Username  string `uri:"username" binding:"required,alphanum"`
	BookingID string `uri:"booking_id" binding:"required"`
}

type OrderTxResult struct {
	Order          db.Order        `json:"order"`
	UserOrder      userResponse    `json:"user_order"`
	ProductOrdered []db.ItemsOrder `json:"product_ordered"`
}

// @Summary User Get Detail Order By Booking ID
// @ID getDetailOrderByBookingID
// @Produce json
// @Accept json
// @Param username path string true "Username"
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags User
// @Success 200 {object} db.OrderTxResult
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/users/{username}/orders/{booking_id}/detail [get]
func (server *Server) getDetailOrderByBookingID(ctx *gin.Context) {
	var req getOrderRequest
	var result OrderTxResult

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
	result.UserOrder = newUserResponse(user)

	order, err := server.store.GetOrder(ctx, req.BookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	result.Order = order

	itemsOrder, err := server.store.ListItemsOrderByBookingID(ctx, order.BookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	result.ProductOrdered = itemsOrder

	ctx.JSON(http.StatusOK, result)
}

type adminGetOrderRequest struct {
	BookingID string `uri:"booking_id" binding:"required"`
}

// @Summary Admin Get Detail Order By Booking ID
// @ID adminGetDetailOrderByBookingID
// @Produce json
// @Accept json
// @Param booking_id path string true "BookingID"
// @Security bearerAuth
// @Tags Admin
// @Success 200 {object} db.OrderTxResult
// @Failure 400 {string} error
// @Failure 401 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/orders/{booking_id}/detail [get]
func (server *Server) adminGetDetailOrderByBookingID(ctx *gin.Context) {
	var req adminGetOrderRequest
	var result OrderTxResult

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, req.BookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	result.Order = order

	user, err := server.store.GetUser(ctx, order.UserBooking)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	result.UserOrder = newUserResponse(user)

	itemsOrder, err := server.store.ListItemsOrderByBookingID(ctx, order.BookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	result.ProductOrdered = itemsOrder

	ctx.JSON(http.StatusOK, result)
}
