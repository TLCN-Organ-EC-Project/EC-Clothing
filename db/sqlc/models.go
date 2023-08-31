// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DescriptionsProduct struct {
	ID          int64  `json:"id"`
	ProductID   int64  `json:"product_id"`
	Gender      string `json:"gender"`
	Material    string `json:"material"`
	Size        string `json:"size"`
	SizeOfModel string `json:"size_of_model"`
}

type Feedback struct {
	ID               int64     `json:"id"`
	UserComment      string    `json:"user_comment"`
	ProductCommented int64     `json:"product_commented"`
	Rating           string    `json:"rating"`
	Commention       string    `json:"commention"`
	CreatedAt        time.Time `json:"created_at"`
}

type ImgsProduct struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Image     string `json:"image"`
}

type ItemsOrder struct {
	ID        int64   `json:"id"`
	BookingID string  `json:"booking_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	BookingID   string    `json:"booking_id"`
	UserBooking string    `json:"user_booking"`
	PromotionID string    `json:"promotion_id"`
	Status      string    `json:"status"`
	BookingDate time.Time `json:"booking_date"`
	Address     string    `json:"address"`
	Province    int64     `json:"province"`
	// must be positive
	Tax           float64 `json:"tax"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

type Product struct {
	ID          int64   `json:"id"`
	ProductName string  `json:"product_name"`
	Thumb       string  `json:"thumb"`
	Price       float64 `json:"price"`
}

type ProductsInCategory struct {
	ID         int64 `json:"id"`
	CategoryID int64 `json:"category_id"`
	ProductID  int64 `json:"product_id"`
}

type Promotion struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DiscountPercent float64   `json:"discount_percent"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

type Province struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Store struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Size      string `json:"size"`
	Quantity  int32  `json:"quantity"`
}

type User struct {
	Username                 string    `json:"username"`
	HashedPassword           string    `json:"hashed_password"`
	FullName                 string    `json:"full_name"`
	Email                    string    `json:"email"`
	Phone                    string    `json:"phone"`
	Address                  string    `json:"address"`
	Province                 int64     `json:"province"`
	Role                     int64     `json:"role"`
	PasswordChangedAt        time.Time `json:"password_changed_at"`
	CreatedAt                time.Time `json:"created_at"`
	ResetPasswordToken       string    `json:"reset_password_token"`
	RspasswordTokenExpiredAt time.Time `json:"rspassword_token_expired_at"`
}
