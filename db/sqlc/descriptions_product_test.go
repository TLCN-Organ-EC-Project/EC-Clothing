package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomDescriptionsProduct(t *testing.T, product Product) DescriptionsProduct {

	arg := CreateDescriptionProductParams{
		ProductID:   product.ID,
		Gender:      util.RandomOwner(),
		Material:    util.RandomOwner(),
		Size:        "Unisex",
		SizeOfModel: "M",
	}

	descriptionsproduct, err := testQueries.CreateDescriptionProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, descriptionsproduct)

	require.Equal(t, arg.ProductID, descriptionsproduct.ProductID)
	require.Equal(t, arg.Gender, descriptionsproduct.Gender)
	require.Equal(t, arg.Material, descriptionsproduct.Material)
	require.Equal(t, arg.Size, descriptionsproduct.Size)
	require.Equal(t, arg.SizeOfModel, descriptionsproduct.SizeOfModel)
	return descriptionsproduct
}

func TestCreateDescriptionsProduct(t *testing.T) {
	product := createRandomProduct(t)
	createRandomDescriptionsProduct(t, product)
}

func TestGetDescriptionsProduct(t *testing.T) {
	product := createRandomProduct(t)
	descriptionsproduct1 := createRandomDescriptionsProduct(t, product)
	descriptionsproduct2, err := testQueries.GetDescriptionProductByID(context.Background(), descriptionsproduct1.ProductID)
	require.NoError(t, err)
	require.NotEmpty(t, descriptionsproduct2)
	require.Equal(t, descriptionsproduct1.ProductID, descriptionsproduct2.ProductID)
	require.Equal(t, product.ID, descriptionsproduct2.ProductID)

	require.Equal(t, descriptionsproduct1.Gender, descriptionsproduct2.Gender)
	require.Equal(t, descriptionsproduct1.Material, descriptionsproduct2.Material)
	require.Equal(t, descriptionsproduct1.Size, descriptionsproduct2.Size)
	require.Equal(t, descriptionsproduct1.SizeOfModel, descriptionsproduct2.SizeOfModel)

}

func TestListDescriptionsProduct(t *testing.T) {
	for i := 0; i < 5; i++ {
		product := createRandomProduct(t)
		createRandomDescriptionsProduct(t, product)
	}

	arg := ListDescriptionProductParams{
		Limit:  5,
		Offset: 0,
	}

	descriptionsproducts, err := testQueries.ListDescriptionProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, descriptionsproducts)
	require.Len(t, descriptionsproducts, 5)

	for _, descriptionsproduct := range descriptionsproducts {
		require.NotEmpty(t, descriptionsproduct)
	}
}

func TestUpdateDescriptionProduct(t *testing.T) {
	product := createRandomProduct(t)
	descriptionsproduct1 := createRandomDescriptionsProduct(t, product)

	arg := UpdateDescriptionProductParams{
		ProductID:   product.ID,
		Gender:      util.RandomOwner(),
		Material:    util.RandomOwner(),
		Size:        "L",
		SizeOfModel: "M",
	}

	descriptionsproduct2, err := testQueries.UpdateDescriptionProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, descriptionsproduct2)
	require.Equal(t, descriptionsproduct1.ID, descriptionsproduct2.ID)
	require.Equal(t, descriptionsproduct1.ProductID, descriptionsproduct2.ProductID)

	require.Equal(t, arg.Gender, descriptionsproduct2.Gender)
	require.Equal(t, arg.Material, descriptionsproduct2.Material)
	require.Equal(t, arg.Size, descriptionsproduct2.Size)
	require.Equal(t, arg.SizeOfModel, descriptionsproduct2.SizeOfModel)
}

func TestDeleteDescriptionProduct(t *testing.T) {
	product := createRandomProduct(t)
	descriptionsproduct1 := createRandomDescriptionsProduct(t, product)

	err := testQueries.DeleteDescriptionProduct(context.Background(), descriptionsproduct1.ProductID)
	require.NoError(t, err)

	descriptionsproduct2, err := testQueries.GetDescriptionProductByID(context.Background(), descriptionsproduct1.ProductID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, descriptionsproduct2)
}
