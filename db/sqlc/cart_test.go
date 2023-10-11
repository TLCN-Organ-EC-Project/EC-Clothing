package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomCart(t *testing.T, user User, product Product) Cart {
	quantity := util.RandomInt(1, 10)
	arg := CreateCartParams {
		Username: user.Username,
		ProductID: product.ID,
		Quantity: int32(quantity),
		Size: "S",
		Price: product.Price * float64(quantity),
	}

	cart, err := testQueries.CreateCart(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cart)

	require.Equal(t, arg.Username, cart.Username)
	require.Equal(t, arg.ProductID, cart.ProductID)
	require.Equal(t, arg.Quantity, cart.Quantity)
	require.Equal(t, arg.Size, cart.Size)
	require.Equal(t, arg.Price, cart.Price)

	return cart
}

func TestCreateCart(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t)

	createRandomCart(t, user, product)
}

func TestGetCart(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t)
	cart1 := createRandomCart(t, user, product)

	cart2, err := testQueries.GetCart(context.Background(), cart1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cart2)

	require.Equal(t, cart1.Username, cart2.Username)
	require.Equal(t, cart1.ProductID, cart2.ProductID)
	require.Equal(t, cart1.Quantity, cart2.Quantity)
	require.Equal(t, cart1.Size, cart2.Size)
	require.Equal(t, cart1.Price, cart2.Price)
}

func TestListCartOfUser(t *testing.T) {
	user := createRandomUser(t)

	for i:=0; i < 5; i++ {
		product := createRandomProduct(t)
		createRandomCart(t, user, product)
	}

	arg := ListCartOfUserParams {
		Username: user.Username,
		Limit: 5,
		Offset: 0,
	}
	carts, err := testQueries.ListCartOfUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, carts)
	require.Len(t, carts, 5)

	for _, cart := range carts {
		require.NotEmpty(t, cart)
		require.Equal(t, user.Username, cart.Username)
	}
}

func TestUpdateCart(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t)
	cart1 := createRandomCart(t, user, product)

	quantity := util.RandomInt(1, 20)

	arg := UpdateCartParams {
		ID: cart1.ID,
		Quantity: int32(quantity),
		Size: "M",
		Price: product.Price * float64(quantity) ,
	}

	cart2, err := testQueries.UpdateCart(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cart2)

	require.Equal(t, cart1.ID, cart2.ID)
	require.Equal(t, cart1.ProductID, cart2.ProductID)
	require.Equal(t, cart2.Size, "M")
}

func TestDeleteCart(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t)
	cart := createRandomCart(t, user, product)

	err := testQueries.DeleteCart(context.Background(), cart.ID)
	require.NoError(t, err)

	cart2, err := testQueries.GetCart(context.Background(), cart.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, cart2)
}

