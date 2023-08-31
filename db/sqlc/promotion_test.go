package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomPromotion(t *testing.T) Promotion {

	startDate := time.Now()
	endDate := startDate.Add(time.Duration(5) * 24 * time.Hour)
	arg := CreatePromotionParams {
		Title: util.RandomOwner(),
		Description: util.RandomOwner(),
		DiscountPercent: float64(util.RandomInt(5, 20)),
		StartDate: startDate,
		EndDate: endDate,
	}

	promotion, err := testQueries.CreatePromotion(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, promotion)

	require.Equal(t, arg.Title, promotion.Title)
	require.Equal(t, arg.Description, promotion.Description)
	require.Equal(t, arg.DiscountPercent, promotion.DiscountPercent)

	require.WithinDuration(t, arg.StartDate, promotion.StartDate, time.Second)
	require.WithinDuration(t, arg.EndDate, promotion.EndDate, time.Second)
	require.WithinDuration(t, promotion.StartDate, promotion.EndDate.Add(time.Duration(-5) * 24 * time.Hour), time.Second)
	require.NotZero(t, promotion.EndDate)

	return promotion
}

func TestCreatePromotion(t *testing.T) {
	createRandomPromotion(t)
}

func TestGetPromotion(t *testing.T) {
	promotion1 := createRandomPromotion(t)
	promotion2, err := testQueries.GetPromotion(context.Background(), promotion1.Title)
	
	require.NoError(t, err)
	require.NotEmpty(t, promotion2)

	require.Equal(t, promotion1.Title, promotion2.Title)
	require.Equal(t, promotion1.Description, promotion2.Description)
	require.Equal(t, promotion1.DiscountPercent, promotion2.DiscountPercent)
	
	require.WithinDuration(t, promotion1.StartDate, promotion2.StartDate, time.Second)
	require.WithinDuration(t, promotion1.EndDate, promotion2.EndDate, time.Second)
}

func TestListPromotions(t *testing.T) {
	for i:=0; i<5; i++ {
		createRandomPromotion(t)
	}

	arg := ListPromotionsParams{
		Limit: 5,
		Offset: 5,
	}

	promotions, err := testQueries.ListPromotions(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, promotions)
	require.Len(t, promotions, 5)

	for _,promotion := range promotions {
		require.NotEmpty(t, promotion)
	}
}

func TestUpdatePromotion(t *testing.T) {
	promotion1 := createRandomPromotion(t)
	arg := UpdatePromotionParams {
		ID: promotion1.ID,
		Description: util.RandomOwner(),
		DiscountPercent: float64(util.RandomInt(1, 20)),
		EndDate: promotion1.StartDate.Add(time.Duration(2) * 24 * time.Hour),
	}

	promotion2, err := testQueries.UpdatePromotion(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, promotion2)
	require.Equal(t, promotion1.ID, promotion2.ID)

	require.Equal(t, promotion1.ID, promotion2.ID)
	require.Equal(t, arg.Description, promotion2.Description)
	require.Equal(t, arg.DiscountPercent, promotion2.DiscountPercent)
	require.WithinDuration(t, arg.EndDate, promotion2.EndDate, time.Second)

	require.WithinDuration(t, promotion1.StartDate, promotion2.StartDate, time.Second)
	require.WithinDuration(t, promotion2.StartDate, promotion2.EndDate.Add(time.Duration(-2) * 24 * time.Hour), time.Second)
}

func TestDeletePromotion(t *testing.T) {
	promotion1 := createRandomPromotion(t)
	err := testQueries.DeletePromotion(context.Background(), promotion1.ID)
	require.NoError(t, err)

	promotion2, err := testQueries.GetPromotion(context.Background(), promotion1.Title)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, promotion2)
}