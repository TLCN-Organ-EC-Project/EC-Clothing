package db

import (
	"context"
	"testing"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomRole(t *testing.T) Role {
	name := util.RandomPhoneNo()
	role, err := testQueries.CreateRole(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, name, role.Name)

	return role
}

func TestCreateRole(t *testing.T) {
	createRandomRole(t)
}
