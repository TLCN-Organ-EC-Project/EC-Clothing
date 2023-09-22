package db

import (
	"context"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomProvince(t *testing.T) Province {
	name := util.RandomOwner()
	province, err := testQueries.CreateProvince(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, province)

	require.Equal(t, name, province.Name)

	return province
}

func TestCreateProvince(t *testing.T) {
	createRandomProvince(t)
}
