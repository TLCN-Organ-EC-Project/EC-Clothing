package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T, user User, promotion Promotion) Order {

	arg := CreateOrderParams{
		BookingID:   util.RandomOrderCode(),
		UserBooking: user.Username,
		PromotionID: promotion.Title,
		Address:     util.RandomOwner(),
		Province:    util.RandomInt(1, 63),
		Amount:      util.RandomFloat(20, 500),
		Tax:         0.1,
	}

	order, err := testQueries.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.Equal(t, arg.BookingID, order.BookingID)
	require.Equal(t, arg.UserBooking, order.UserBooking)
	require.Equal(t, arg.PromotionID, order.PromotionID)
	require.Equal(t, arg.Address, order.Address)
	require.Equal(t, arg.Province, order.Province)
	require.Equal(t, arg.Amount, order.Amount)
	require.Equal(t, arg.Tax, order.Tax)

	return order
}

func TestCreateOrder(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)

	createRandomOrder(t, user, promotion)
}

func TestGetOrder(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)
	order1 := createRandomOrder(t, user, promotion)

	order2, err := testQueries.GetOrder(context.Background(), order1.BookingID)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.BookingID, order2.BookingID)
	require.Equal(t, order1.UserBooking, order2.UserBooking)
	require.Equal(t, order1.PromotionID, order2.PromotionID)
	require.Equal(t, order1.Address, order2.Address)
	require.Equal(t, order1.Province, order2.Province)
	require.Equal(t, order1.Amount, order2.Amount)
	require.Equal(t, order1.Tax, order2.Tax)
}

func TestListOrderByUser(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)

	for i := 0; i < 5; i++ {
		createRandomOrder(t, user, promotion)
	}

	arg := ListOrderByUserParams{
		UserBooking: user.Username,
		Limit:       5,
		Offset:      0,
	}

	orders, err := testQueries.ListOrderByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orders)
	require.Len(t, orders, 5)

	for _, order := range orders {
		require.NotEmpty(t, order)
	}
}

func TestListOrder(t *testing.T) {
	arg := ListOrderParams{
		Limit:  10,
		Offset: 0,
	}

	orders, err := testQueries.ListOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orders)
	require.Len(t, orders, 10)

	for _, order := range orders {
		require.NotEmpty(t, order)
	}
}

func TestDeleteOrder(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)
	order1 := createRandomOrder(t, user, promotion)

	err := testQueries.DeleteOrder(context.Background(), order1.BookingID)
	require.NoError(t, err)

	order2, err := testQueries.GetOrder(context.Background(), order1.BookingID)
	require.Error(t, err)
	require.Empty(t, order2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
