package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomImgProduct(t *testing.T, product Product) ImgsProduct {
	arg := CreateImgProductParams {
		ProductID: product.ID,
		Image: util.RandomOwner(),
	}

	imgproduct, err := testQueries.CreateImgProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, imgproduct)

	require.Equal(t, arg.ProductID, imgproduct.ProductID)
	require.Equal(t, arg.Image, imgproduct.Image)


	arg1 := CreateImgProductParams {
		ProductID: 15000,
		Image: util.RandomOwner(),
	}
	imgproduct2, err := testQueries.CreateImgProduct(context.Background(), arg1)
	require.Error(t, err)
	require.Empty(t, imgproduct2)
	require.Contains(t, err.Error(), "pq: insert or update on table \"imgs_product\" violates foreign key constraint \"imgs_product_product_id_fkey\"")


	return imgproduct
}

func TestCreateImgProduct(t *testing.T) {
	product := createRandomProduct(t)
	createRandomImgProduct(t, product)
}

func TestGetImgProduct(t *testing.T) {
	product := createRandomProduct(t)
	imgsproduct1 := createRandomImgProduct(t, product)
	imgsproduct2, err := testQueries.GetImgProduct(context.Background(), imgsproduct1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, imgsproduct2)

	require.Equal(t, imgsproduct1.ProductID, imgsproduct2.ProductID)
	require.Equal(t, imgsproduct1.Image, imgsproduct2.Image)

}

func TestListImgProductsByProductID(t *testing.T) {
	product := createRandomProduct(t)

	for i:=0; i<5; i++ {
		createRandomImgProduct(t, product)
	}


	arg := ListImgProductsParams {
		ProductID: product.ID,
		Limit: 5,
		Offset: 0,
	}

	imgsproducts, err := testQueries.ListImgProducts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, imgsproducts)
	require.Len(t, imgsproducts, 5)

	for _, imgsproduct := range imgsproducts {
		require.NotEmpty(t, imgsproduct)
	}

}

func TestUpdateImgProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	product2 := createRandomProduct(t)
	imgproduct := createRandomImgProduct(t, product1)
	imgproduct2 := createRandomImgProduct(t, product2)

	arg := UpdateImgProductParams {
		ID: imgproduct.ID,
		ProductID: product1.ID,
		Image: util.RandomOwner(),
	}

	arg2 := UpdateImgProductParams {
		ID: imgproduct2.ID,
		ProductID: product1.ID,
		Image: imgproduct2.Image,
	}

	update_imgproduct, err := testQueries.UpdateImgProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, update_imgproduct)
	require.Equal(t, imgproduct.ID, update_imgproduct.ID)
	require.Equal(t, product1.ID, update_imgproduct.ProductID)

	update_imgproduct2, err := testQueries.UpdateImgProduct(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, update_imgproduct2)
	require.Equal(t, imgproduct2.ID, update_imgproduct2.ID)
	require.Equal(t, product1.ID, update_imgproduct2.ProductID)
	require.Equal(t, imgproduct2.Image, update_imgproduct2.Image)
}

func TestDeleteImgProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	imgproduct1 := createRandomImgProduct(t, product1)
	err := testQueries.DeleteImgProduct(context.Background(), imgproduct1.ID)
	require.NoError(t, err)

	imgproduct2, err := testQueries.GetImgProduct(context.Background(), imgproduct1.ID)
	require.Error(t, err)
	require.Empty(t, imgproduct2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

