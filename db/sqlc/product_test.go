package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T) Product {
	arg := CreateProductParams{
		ProductName: util.RandomOwner(),
		Thumb:       util.RandomOwner(),
		Price:       float64(util.RandomInt(20, 50)),
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.ProductName, product.ProductName)
	require.Equal(t, arg.Thumb, product.Thumb)
	require.Equal(t, arg.Price, product.Price)

	return product
}

func TestCreateProduct(t *testing.T) {
	createRandomProduct(t)
}

func TestGetProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	product2, err := testQueries.GetProduct(context.Background(), product1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ProductName, product2.ProductName)
	require.Equal(t, product1.Thumb, product2.Thumb)
	require.Equal(t, product1.Price, product2.Price)
}

func TestListProducts(t *testing.T) {

	for i := 0; i < 5; i++ {
		createRandomProduct(t)
	}

	arg := ListProductsParams{
		Limit:  5,
		Offset: 5,
	}

	products, err := testQueries.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, products)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}

}

func TestUpdateProduct(t *testing.T) {
	product1 := createRandomProduct(t)

	arg := UpdateProductParams{
		ID:          product1.ID,
		ProductName: util.RandomOwner(),
		Thumb:       util.RandomOwner(),
		Price:       float64(util.RandomInt(2, 50)),
	}

	product2, err := testQueries.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product2)
	require.Equal(t, product1.ID, product2.ID)

	require.Equal(t, arg.ProductName, product2.ProductName)
	require.Equal(t, arg.Thumb, product2.Thumb)
	require.Equal(t, arg.Price, product2.Price)
}

func TestDeleteProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	err := testQueries.DeleteProduct(context.Background(), product1.ID)
	require.NoError(t, err)

	product2, err := testQueries.GetProduct(context.Background(), product1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, product2)
}
