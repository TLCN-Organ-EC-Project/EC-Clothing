package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomOwner()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	province := createRandomProvince(t)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandonEmail(),
		Phone:          util.RandomPhoneNo(),
		Address:        util.RandomOwner(),
		Province:       province.ID,
		Role:           2,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	// No error and no empty when create the new user
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// All the field of arg and user must be equal
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.Address, user.Address)
	require.Equal(t, arg.Province, user.Province)
	require.Equal(t, arg.Role, user.Role)

	//
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	// No error and no empty when create the new user
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	// All the field of user1 and user2 must be equal
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user2.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user2.Address, user2.Address)
	require.Equal(t, user1.Province, user2.Province)
	require.Equal(t, user1.Role, user2.Role)

	// 2 times are within duration delta of each other
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)

	// No error and no empty when create the new user
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	// All the field of user1 and user2 must be equal
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user2.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user2.Address, user2.Address)
	require.Equal(t, user1.Province, user2.Province)
	require.Equal(t, user1.Role, user2.Role)

	// 2 times are within duration delta of each other
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestListUsers(t *testing.T) {
	arg := ListUsersParams{
		Limit: 5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	province := createRandomProvince(t)
	arg := UpdateUserParams {
		Username: user1.Username,
		FullName: util.RandomOwner(),
		Email: util.RandonEmail(),
		Phone: util.RandomPhoneNo(),
		Address: util.RandomOwner(),
		Province: province.ID,
		
	}
	
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestChangeUserPassword(t *testing.T) {
	newPassword, err := util.HashPassword(util.RandomOwner())
	require.NoError(t, err)
	user1 := createRandomUser(t)

	arg := ChangeUserPasswordParams{
		Username: user1.Username,
		HashedPassword: newPassword,
		PasswordChangedAt: time.Now(),
	}

	user2, err := testQueries.ChangeUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.NotZero(t, user2.PasswordChangedAt)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}