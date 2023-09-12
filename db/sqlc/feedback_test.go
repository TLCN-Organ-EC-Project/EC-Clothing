package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/stretchr/testify/require"
)

func createRandomFeedback(t *testing.T, product Product) Feedback {
	user := createRandomUser(t)

	arg := CreateFeedbackParams{
		UserComment:      user.Username,
		ProductCommented: product.ID,
		Rating:           "Well",
		Commention:       util.RandomOwner(),
		CreatedAt:        time.Now(),
	}

	feedback, err := testQueries.CreateFeedback(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, feedback)

	require.Equal(t, arg.UserComment, feedback.UserComment)
	require.Equal(t, arg.ProductCommented, feedback.ProductCommented)
	require.Equal(t, arg.Rating, feedback.Rating)
	require.Equal(t, arg.Commention, feedback.Commention)

	return feedback
}

func TestCreateFeedback(t *testing.T) {
	product := createRandomProduct(t)
	createRandomFeedback(t, product)
}

func TestGetFeedback(t *testing.T) {
	product := createRandomProduct(t)
	feedback1 := createRandomFeedback(t, product)
	feedback2, err := testQueries.GetFeedback(context.Background(), feedback1.ID)

	// No error and no empty when create the new user
	require.NoError(t, err)
	require.NotEmpty(t, feedback2)

	// All the field of user1 and user2 must be equal
	require.Equal(t, feedback1.ID, feedback2.ID)
	require.Equal(t, feedback1.UserComment, feedback2.UserComment)
	require.Equal(t, feedback1.ProductCommented, feedback2.ProductCommented)
	require.Equal(t, feedback1.Rating, feedback2.Rating)
	require.Equal(t, feedback1.Commention, feedback2.Commention)

	// 2 times are within duration delta of each other
	require.WithinDuration(t, feedback1.CreatedAt, feedback2.CreatedAt, time.Second)
}

func TestListFeedbacksByProduct(t *testing.T) {
	product := createRandomProduct(t)
	for i := 0; i < 5; i++ {
		createRandomFeedback(t, product)
	}

	arg := ListFeedbacksParams{
		ProductCommented: product.ID,
		Limit:            5,
		Offset:           0,
	}

	feedbacks, err := testQueries.ListFeedbacks(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, feedbacks)

	require.Len(t, feedbacks, 5)

	for _, feedback := range feedbacks {
		require.NotEmpty(t, feedback)
	}
}

func TestUpdateFeedback(t *testing.T) {
	product := createRandomProduct(t)
	feedback1 := createRandomFeedback(t, product)

	arg := UpdateFeedbackParams{
		ID:         feedback1.ID,
		Rating:     "Well Update",
		Commention: util.RandomOwner(),
	}

	feedback2, err := testQueries.UpdateFeedback(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, feedback2)
	require.Equal(t, feedback1.ID, feedback2.ID)

	require.Equal(t, feedback1.UserComment, feedback2.UserComment)
	require.Equal(t, feedback1.ProductCommented, feedback2.ProductCommented)

	require.Equal(t, arg.Rating, feedback2.Rating)
	require.Equal(t, arg.Commention, feedback2.Commention)

	require.WithinDuration(t, feedback1.CreatedAt, feedback2.CreatedAt, time.Second)

}

func TestDeleteFeedback(t *testing.T) {
	product := createRandomProduct(t)
	feedback1 := createRandomFeedback(t, product)

	err := testQueries.DeleteFeedback(context.Background(), feedback1.ID)
	require.NoError(t, err)

	feedback2, err := testQueries.GetFeedback(context.Background(), feedback1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, feedback2)
}
