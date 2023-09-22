// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	ChangeUserPassword(ctx context.Context, arg ChangeUserPasswordParams) (User, error)
	CreateCategory(ctx context.Context, name string) (Category, error)
	CreateDescriptionProduct(ctx context.Context, arg CreateDescriptionProductParams) (DescriptionsProduct, error)
	CreateFeedback(ctx context.Context, arg CreateFeedbackParams) (Feedback, error)
	CreateImgProduct(ctx context.Context, arg CreateImgProductParams) (ImgsProduct, error)
	CreateItemsOrder(ctx context.Context, arg CreateItemsOrderParams) (ItemsOrder, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateProductsInCategory(ctx context.Context, arg CreateProductsInCategoryParams) (ProductsInCategory, error)
	CreatePromotion(ctx context.Context, arg CreatePromotionParams) (Promotion, error)
	CreateProvince(ctx context.Context, name string) (Province, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCategory(ctx context.Context, id int64) error
	DeleteDescriptionProduct(ctx context.Context, productID int64) error
	DeleteFeedback(ctx context.Context, id int64) error
	DeleteImgProduct(ctx context.Context, id int64) error
	DeleteItemsOrder(ctx context.Context, id int64) error
	DeleteOrder(ctx context.Context, bookingID string) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteProductsInCategory(ctx context.Context, id int64) error
	DeletePromotion(ctx context.Context, id int64) error
	DeleteStore(ctx context.Context, arg DeleteStoreParams) error
	DeleteUser(ctx context.Context, username string) error
	GetCategory(ctx context.Context, id int64) (Category, error)
	GetDescriptionProductByID(ctx context.Context, productID int64) (DescriptionsProduct, error)
	GetFeedback(ctx context.Context, id int64) (Feedback, error)
	GetImgProduct(ctx context.Context, id int64) (ImgsProduct, error)
	GetItemsOrder(ctx context.Context, id int64) (ItemsOrder, error)
	GetOrder(ctx context.Context, bookingID string) (Order, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetProductsInCategoryByID(ctx context.Context, arg GetProductsInCategoryByIDParams) (ProductsInCategory, error)
	GetPromotion(ctx context.Context, title string) (Promotion, error)
	GetProvince(ctx context.Context, name string) (Province, error)
	GetRole(ctx context.Context, name string) (Role, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetStore(ctx context.Context, arg GetStoreParams) (Store, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByResetPassToken(ctx context.Context, resetPasswordToken string) (User, error)
	ListCategories(ctx context.Context) ([]Category, error)
	ListDescriptionProduct(ctx context.Context, arg ListDescriptionProductParams) ([]DescriptionsProduct, error)
	ListFeedbacks(ctx context.Context, arg ListFeedbacksParams) ([]Feedback, error)
	ListImgProducts(ctx context.Context, productID int64) ([]ImgsProduct, error)
	ListItemsOrderByBookingID(ctx context.Context, arg ListItemsOrderByBookingIDParams) ([]ItemsOrder, error)
	ListOrder(ctx context.Context, arg ListOrderParams) ([]Order, error)
	ListOrderByUser(ctx context.Context, arg ListOrderByUserParams) ([]Order, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListProductsInCategory(ctx context.Context, arg ListProductsInCategoryParams) ([]ProductsInCategory, error)
	ListPromotions(ctx context.Context, arg ListPromotionsParams) ([]Promotion, error)
	ListProvinces(ctx context.Context) ([]string, error)
	ListStore(ctx context.Context, productID int64) ([]Store, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateDescriptionProduct(ctx context.Context, arg UpdateDescriptionProductParams) (DescriptionsProduct, error)
	UpdateFeedback(ctx context.Context, arg UpdateFeedbackParams) (Feedback, error)
	UpdateImgProduct(ctx context.Context, arg UpdateImgProductParams) (ImgsProduct, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductsInCategory(ctx context.Context, arg UpdateProductsInCategoryParams) (ProductsInCategory, error)
	UpdatePromotion(ctx context.Context, arg UpdatePromotionParams) (Promotion, error)
	UpdateResetPasswordToken(ctx context.Context, arg UpdateResetPasswordTokenParams) (User, error)
	UpdateStore(ctx context.Context, arg UpdateStoreParams) (Store, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
