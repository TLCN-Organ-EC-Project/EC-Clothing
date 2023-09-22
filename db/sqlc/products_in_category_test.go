package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomProductsInCategory(t *testing.T, category Category) ProductsInCategory {
	product := createRandomProduct(t)
	arg := CreateProductsInCategoryParams{
		ProductID:  product.ID,
		CategoryID: category.ID,
	}

	productincategory, err := testQueries.CreateProductsInCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productincategory)

	require.Equal(t, arg.ProductID, productincategory.ProductID)
	require.Equal(t, arg.CategoryID, productincategory.CategoryID)

	return productincategory
}

func TestRandomProductsInCategory(t *testing.T) {
	category := createRandomCategory(t)
	createRandomProductsInCategory(t, category)
}

func TestGetProductsInCategory(t *testing.T) {
	category := createRandomCategory(t)
	productincategory1 := createRandomProductsInCategory(t, category)
	productincategory2, err := testQueries.GetProductsInCategoryByID(context.Background(), GetProductsInCategoryByIDParams{
		CategoryID: productincategory1.CategoryID,
		ProductID: productincategory1.ProductID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, productincategory2)

	require.Equal(t, productincategory1.CategoryID, productincategory2.CategoryID)
	require.Equal(t, productincategory1.ProductID, productincategory2.ProductID)
}

func TestListProductsInCategory(t *testing.T) {
	category := createRandomCategory(t)
	for i := 0; i < 5; i++ {
		createRandomProductsInCategory(t, category)
	}

	arg := ListProductsInCategoryParams{
		Limit:      5,
		Offset:     0,
		CategoryID: category.ID,
	}

	productsincategory, err := testQueries.ListProductsInCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productsincategory)
	require.Len(t, productsincategory, 5)

	for _, productincategory := range productsincategory {
		require.NotEmpty(t, productincategory)
	}
}

func TestUpdateProductsInCategory(t *testing.T) {
	category := createRandomCategory(t)
	productincategory1 := createRandomProductsInCategory(t, category)

	product2 := createRandomProduct(t)
	category2 := createRandomCategory(t)

	arg := UpdateProductsInCategoryParams{
		ID:         productincategory1.ID,
		ProductID:  product2.ID,
		CategoryID: category2.ID,
	}

	productincategory2, err := testQueries.UpdateProductsInCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productincategory2)
	require.Equal(t, productincategory1.ID, productincategory2.ID)

	require.Equal(t, arg.ProductID, productincategory2.ProductID)
	require.Equal(t, arg.CategoryID, productincategory2.CategoryID)
}

func TestDeleteProductsInCategory(t *testing.T) {
	category := createRandomCategory(t)
	productincategory1 := createRandomProductsInCategory(t, category)

	err := testQueries.DeleteProductsInCategory(context.Background(), productincategory1.ID)
	require.NoError(t, err)

	productincategory2, err := testQueries.GetProductsInCategoryByID(context.Background(), GetProductsInCategoryByIDParams{
		CategoryID: productincategory1.CategoryID,
		ProductID: productincategory1.ProductID,
	})
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, productincategory2)
}
