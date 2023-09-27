package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomItemsOrder(t *testing.T, order Order) ItemsOrder {
	product := createRandomProduct(t)
	quantity := int32(util.RandomInt(1, 10))

	arg := CreateItemsOrderParams{
		BookingID: order.BookingID,
		ProductID: product.ID,
		Quantity:  quantity,
		Price:     product.Price * float64(quantity),
	}

	itemsorder, err := testQueries.CreateItemsOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, itemsorder)

	require.Equal(t, arg.BookingID, itemsorder.BookingID)
	require.Equal(t, arg.ProductID, itemsorder.ProductID)
	require.Equal(t, arg.Quantity, itemsorder.Quantity)
	require.Equal(t, arg.Price, itemsorder.Price)

	return itemsorder
}

func TestCreateItemsOrder(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)
	order := createRandomOrder(t, user, promotion)

	createRandomItemsOrder(t, order)
}

func TestGetItemsOrder(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)
	order := createRandomOrder(t, user, promotion)
	imtemsorder1 := createRandomItemsOrder(t, order)

	imtemsorder2, err := testQueries.GetItemsOrder(context.Background(), imtemsorder1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, imtemsorder2)

	require.Equal(t, imtemsorder1.ID, imtemsorder2.ID)
	require.Equal(t, imtemsorder1.BookingID, imtemsorder2.BookingID)
	require.Equal(t, imtemsorder1.ProductID, imtemsorder2.ProductID)
	require.Equal(t, imtemsorder1.Quantity, imtemsorder2.Quantity)
	require.Equal(t, imtemsorder1.Price, imtemsorder2.Price)
}

func TestListItemsOrderByBookingID(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)
	order := createRandomOrder(t, user, promotion)

	for i := 0; i < 5; i++ {
		createRandomItemsOrder(t, order)
	}

	itemsorders, err := testQueries.ListItemsOrderByBookingID(context.Background(), order.BookingID)
	require.NoError(t, err)
	require.NotEmpty(t, itemsorders)
	require.Len(t, itemsorders, 5)

	for _, itemsorder := range itemsorders {
		require.NotEmpty(t, itemsorder)
	}
}

func TestDeleteItemsOrder(t *testing.T) {
	user := createRandomUser(t)
	promotion := createRandomPromotion(t)
	order := createRandomOrder(t, user, promotion)
	imtemsorder1 := createRandomItemsOrder(t, order)

	err := testQueries.DeleteItemsOrder(context.Background(), imtemsorder1.ID)
	require.NoError(t, err)

	imtemsorder2, err := testQueries.GetItemsOrder(context.Background(), imtemsorder1.ID)
	require.Error(t, err)
	require.Empty(t, imtemsorder2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
