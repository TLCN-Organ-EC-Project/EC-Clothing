package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	name := util.RandomOwner()
	category, err := testQueries.CreateCategory(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, name, category.Name)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	category2, err := testQueries.GetCategory(context.Background(), category1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.Name, category2.Name)
}

func TestListCategories(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomCategory(t)
	}


	categories, err := testQueries.ListCategories(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}

func TestUpdateCategory(t *testing.T) {
	category1 := createRandomCategory(t)

	arg := UpdateCategoryParams{
		ID:   category1.ID,
		Name: util.RandomOwner(),
	}

	category2, err := testQueries.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category2)
	require.Equal(t, category1.ID, category2.ID)

	require.NotEqual(t, category1.Name, category2.Name)
}

func TestDeleteCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	err := testQueries.DeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	category2, err := testQueries.GetCategory(context.Background(), category1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, category2)
}
