package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/XuanHieuHo/EC_Clothing/db/sqlc"
	"github.com/gin-gonic/gin"
)

type totalIncomeMonthlyRequest struct {
	Month int `form:"month" binding:"required,min=1,max=12"`
	Year  int `form:"year" binding:"required,min=2023"`
}

// @Summary Admin Get Income Monthly
// @ID getTotalIncomeMonthly
// @Produce json
// @Accept json
// @Tags Revenue/Admin
// @Param data query totalIncomeMonthlyRequest true "totalIncomeMonthlyRequest data"
// @Security bearerAuth
// @Success 200 {float} incomeMonthly
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/income/monthly [get]
func (server *Server) getTotalIncomeMonthly(ctx *gin.Context) {
	var req totalIncomeMonthlyRequest
	var totalIncome float64
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startMonth := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)
	endMonth := time.Date(req.Year, time.Month(req.Month)+1, 1, 0, 0, 0, 0, time.Local).Add(-time.Hour * 24)

	_, err := server.store.GetOrderByDate(ctx, db.GetOrderByDateParams{
		BookingDate:   startMonth,
		BookingDate_2: endMonth,
		Status:        "validated",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			totalIncome = 0
			ctx.JSON(http.StatusOK, totalIncome)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.TotalIncomeParams{
		BookingDate:   startMonth,
		BookingDate_2: endMonth,
		Status:        "validated",
	}

	totalIncome, err = server.store.TotalIncome(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, totalIncome)
}

type totalIncomeYearlyRequest struct {
	Year int `form:"year" binding:"required,min=2023"`
}

// @Summary Admin Get Income Yearly
// @ID getTotalIncomeYearly
// @Produce json
// @Accept json
// @Tags Revenue/Admin
// @Param data query totalIncomeYearlyRequest true "totalIncomeYearlyRequest data"
// @Security bearerAuth
// @Success 200 {float} incomeYearly
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/income/yearly [get]
func (server *Server) getTotalIncomeYearly(ctx *gin.Context) {
	var req totalIncomeYearlyRequest
	var totalIncome float64

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startYear := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.Local)
	endYear := time.Date(req.Year, 12, 31, 0, 0, 0, 0, time.Local)

	_, err := server.store.GetOrderByDate(ctx, db.GetOrderByDateParams{
		BookingDate:   startYear,
		BookingDate_2: endYear,
		Status:        "validated",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			totalIncome = 0
			ctx.JSON(http.StatusOK, totalIncome)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.TotalIncomeParams{
		BookingDate:   startYear,
		BookingDate_2: endYear,
		Status:        "validated",
	}

	totalIncome, err = server.store.TotalIncome(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, totalIncome)
}

// @Summary Admin Get Statistics Product
// @ID getStatisticsProduct
// @Produce json
// @Accept json
// @Tags Revenue/Admin
// @Security bearerAuth
// @Success 200 {object} db.StatisticsProductRow
// @Failure 400 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/admin/statistics_product [get]
func (server *Server) getStatisticsProduct(ctx *gin.Context) {

	statistics, err := server.store.StatisticsProduct(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, statistics)
}