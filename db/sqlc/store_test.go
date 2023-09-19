package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

const (
	S        = "S"
	M        = "M"
	L        = "L"
	XL       = "XL"
	XXL      = "XXL"
	OVERSIZE = "OVERSIZE"
)

func createRandomStore(t *testing.T, product Product, size string) Store {
	arg := CreateStoreParams{
		ProductID: product.ID,
		Quantity:  int32(util.RandomInt(20, 50)),
		Size:      size,
	}

	store, err := testQueries.CreateStore(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, store)

	require.Equal(t, arg.ProductID, store.ProductID)
	require.Equal(t, arg.Size, store.Size)
	require.Equal(t, arg.Quantity, store.Quantity)
	return store
}

func TestCreateStore(t *testing.T) {
	product := createRandomProduct(t)
	size := util.RandomSize()
	createRandomStore(t, product, size)
}

func TestGetStore(t *testing.T) {
	product := createRandomProduct(t)
	size := util.RandomSize()
	store1 := createRandomStore(t, product, size)
	arg := GetStoreParams{
		ProductID: store1.ProductID,
		Size:      store1.Size,
	}

	store2, err := testQueries.GetStore(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, store2)

	require.Equal(t, store1.ID, store2.ID)
	require.Equal(t, store1.ProductID, store2.ProductID)
	require.Equal(t, store1.Size, store2.Size)
	require.Equal(t, store1.Quantity, store2.Quantity)
}

func TestListStore(t *testing.T) {
	product := createRandomProduct(t)
	sizes := []string{S, M, L, XL, XXL, OVERSIZE}

	for _, size := range sizes {
		createRandomStore(t, product, size)
	}

	arg := ListStoreParams{
		ProductID: product.ID,
		Limit:     6,
		Offset:    0,
	}

	stores, err := testQueries.ListStore(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, stores)
	require.Len(t, stores, 6)

	for _, store := range stores {
		require.NotEmpty(t, store)
	}
}

func TestUpdateStore(t *testing.T) {
	product := createRandomProduct(t)
	size := util.RandomSize()
	store := createRandomStore(t, product, size)

	arg := UpdateStoreParams{
		ProductID: store.ProductID,
		Size:      store.Size,
		Quantity:  int32(util.RandomInt(10, 56)),
	}

	update_store, err := testQueries.UpdateStore(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, update_store)

	require.Equal(t, store.ID, update_store.ID)
	require.Equal(t, store.ProductID, update_store.ProductID)
	require.Equal(t, store.Size, update_store.Size)
	require.NotEqual(t, store.Quantity, update_store.Quantity)

}
func TestDeleteStore(t *testing.T) {
	product := createRandomProduct(t)
	size := util.RandomSize()
	store := createRandomStore(t, product, size)

	arg := DeleteStoreParams{
		ProductID: store.ProductID,
		Size:      store.Size,
	}

	err := testQueries.DeleteStore(context.Background(), arg)
	require.NoError(t, err)

	store2, err := testQueries.GetStore(context.Background(), GetStoreParams{
		ProductID: store.ProductID,
		Size:      store.Size,
	})
	require.Error(t, err)
	require.Empty(t, store2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
