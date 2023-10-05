package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/XuanHieuHo/EC_Clothing/util"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/lib/pq"
)

type Stores interface {
	Querier
	CreateProductTx(ctx context.Context, arg CreateProductTxParams) (CreateProductTxResults, error)
	CreateImgProductTx(ctx context.Context, arg AddImageProductTxParams) ([]ImgsProduct, error)
	OrderTx(ctx context.Context, arg OrderTxParams) (OrderTxResult, error)
	UpdateOrderTx(ctx context.Context, arg UpdateOrderTxParams) (UpdateOrderTxResult, error)
	CancelOrderTx(ctx context.Context, arg CancelOrderParams) (string, error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Stores {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type CreateProductTxParams struct {
	ProductName string  `json:"product_name"`
	Thumb       string  `json:"thumb"`
	Price       float64 `json:"price"`
	Gender      string  `json:"gender"`
	Material    string  `json:"material"`
	Size        string  `json:"size"`
	SizeOfModel string  `json:"size_of_model"`
}

type CreateProductTxResults struct {
	ProductName string  `json:"product_name"`
	Thumb       string  `json:"thumb"`
	Price       float64 `json:"price"`
	Gender      string  `json:"gender"`
	Material    string  `json:"material"`
	Size        string  `json:"size"`
	SizeOfModel string  `json:"size_of_model"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user User) userResponse {
	return userResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

func (store *SQLStore) CreateProductTx(ctx context.Context, arg CreateProductTxParams) (CreateProductTxResults, error) {
	var result CreateProductTxResults

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	err = store.execTx(ctx, func(q *Queries) error {
		var err error

		cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
		if err != nil {
			return err
		}

		params := uploader.UploadParams{
			Folder:         "ec-clothing",
			Format:         "jpg",
			Transformation: "f_auto,fl_lossy,q_auto:eco,dpr_auto,w_auto",
		}

		thumb, err := cld.Upload.Upload(ctx, arg.Thumb, params)
		if err != nil {
			return err
		}

		argProduct := CreateProductParams{
			ProductName: arg.ProductName,
			Thumb:       thumb.SecureURL,
			Price:       arg.Price,
		}

		product, err := store.CreateProduct(ctx, argProduct)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "foreign_key_violation", "unique_violation":
					return err
				}
			}
			return err
		}

		argDescriptionsProduct := CreateDescriptionProductParams{
			ProductID:   product.ID,
			Gender:      arg.Gender,
			Material:    arg.Material,
			Size:        arg.Size,
			SizeOfModel: arg.SizeOfModel,
		}

		descriptionsProduct, err := store.CreateDescriptionProduct(ctx, argDescriptionsProduct)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "foreign_key_violation", "unique_violation":
					return err
				}
			}
			return err
		}

		rsp := CreateProductTxResults{
			ProductName: product.ProductName,
			Thumb:       product.Thumb,
			Price:       product.Price,
			Gender:      descriptionsProduct.Gender,
			Material:    descriptionsProduct.Material,
			Size:        descriptionsProduct.Size,
			SizeOfModel: descriptionsProduct.SizeOfModel,
		}

		result = rsp
		return nil
	})
	return result, err
}

type AddImageProductTxParams struct {
	Images []string `json:"images"`
	ID     int64    `json:"id"`
}

func (store *SQLStore) CreateImgProductTx(ctx context.Context, arg AddImageProductTxParams) ([]ImgsProduct, error) {
	var imgProducts []ImgsProduct

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	err = store.execTx(ctx, func(q *Queries) error {
		var err error

		cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
		if err != nil {
			return err
		}

		params := uploader.UploadParams{
			Folder:         "ec-clothing",
			Format:         "jpg",
			Transformation: "f_auto,fl_lossy,q_auto:eco,dpr_auto,w_auto",
		}

		for _, image := range arg.Images {
			img, err := cld.Upload.Upload(ctx, image, params)
			if err != nil {
				return err
			}

			arg := CreateImgProductParams{
				ProductID: arg.ID,
				Image:     img.SecureURL,
			}

			imgProduct, err := store.CreateImgProduct(ctx, arg)
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok {
					switch pqErr.Code.Name() {
					case "foreign_key_violation", "unique_violation":
						return err
					}
				}
				return err
			}
			imgProducts = append(imgProducts, imgProduct)
		}
		return nil
	})

	return imgProducts, err
}

type OrderTxParams struct {
	Username      string   `json:"username"`
	PromotionID   string   `json:"promotion_id"`
	Address       string   `json:"address"`
	Province      string   `json:"province"`
	PaymentMethod string   `json:"payment_method"`
	ProductID     []int64  `json:"product_id"`
	Size          []string `json:"size"`
	Quantity      []int64  `json:"quantity"`
}

type OrderTxResult struct {
	Order          Order        `json:"order"`
	UserOrder      userResponse `json:"user_order"`
	ProductOrdered []ItemsOrder `json:"product_ordered"`
}

func (store *SQLStore) OrderTx(ctx context.Context, arg OrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		if (len(arg.ProductID) != len(arg.Size)) || (len(arg.ProductID) != len(arg.Quantity)) || (len(arg.Size) != len(arg.Quantity)) {
			err := errors.New("the length of one or more of the properties is not equal")
			return err
		}

		user, err := q.GetUser(ctx, arg.Username)
		if err != nil {
			return err
		}

		result.UserOrder = newUserResponse(user)

		bookingID := util.RandomOrderCode()

		province, err := q.GetProvince(ctx, arg.Province)
		if err != nil {
			return err
		}

		var discount float64
		if arg.PromotionID != "none" {
			promotion, err := q.GetPromotion(ctx, arg.PromotionID)
			if err != nil {
				if err == sql.ErrNoRows {
					err = fmt.Errorf("promotion code doesn't exist")
					return err
				} else {
					return err
				}
			}

			if time.Now().After(promotion.EndDate) {
				err = fmt.Errorf("promotion code has expired")
				return err
			}
			discount = promotion.DiscountPercent / 100
		} else {
			discount = 0
		}

		argOrder := CreateOrderParams{
			BookingID:   bookingID,
			UserBooking: user.Username,
			PromotionID: arg.PromotionID,
			Address:     arg.Address,
			Province:    province.ID,
			Tax:         0.1,
			Amount:      0,
		}
		order, err := q.CreateOrder(ctx, argOrder)
		if err != nil {
			return err
		}

		var amount float64
		amount = 0
		for i, productID := range arg.ProductID {
			size := arg.Size[i]
			quantity := arg.Quantity[i]

			product, err := q.GetProduct(ctx, productID)
			if err != nil {
				return err
			}

			store, err := q.GetStore(ctx, GetStoreParams{
				ProductID: product.ID,
				Size:      size,
			})
			if err != nil {
				return err
			}

			if quantity > int64(store.Quantity) {
				err := errors.New("quantity is not enough")
				return err
			}

			argItemOrder := CreateItemsOrderParams{
				BookingID: order.BookingID,
				ProductID: productID,
				Quantity:  int32(quantity),
				Size:      size,
				Price:     product.Price * float64(quantity),
			}

			amount = amount + argItemOrder.Price

			_, err = q.CreateItemsOrder(ctx, argItemOrder)
			if err != nil {
				return err
			}

			_, err = q.UpdateStore(ctx, UpdateStoreParams{
				ProductID: store.ProductID,
				Size:      store.Size,
				Quantity:  store.Quantity - argItemOrder.Quantity,
			})
			if err != nil {
				return err
			}

		}

		amount = (amount + order.Tax*amount) * (1 - discount)

		result.Order, err = q.UpdateAmountOfOrder(ctx, UpdateAmountOfOrderParams{
			BookingID: order.BookingID,
			Amount:    amount,
		})
		if err != nil {
			return err
		}

		result.ProductOrdered, err = q.ListItemsOrderByBookingID(ctx, order.BookingID)
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

type UpdateOrderTxParams struct {
	Username  string   `json:"username"`
	BookingID string   `json:"booking_id"`
	Address   string   `json:"address"`
	Province  string   `json:"province"`
	ProductID []int64  `json:"product_id"`
	Size      []string `json:"size"`
	Quantity  []int64  `json:"quantity"`
}

type UpdateOrderTxResult struct {
	Order          Order        `json:"order"`
	UserOrder      userResponse `json:"user_order"`
	ProductOrdered []ItemsOrder `json:"product_ordered"`
}

func (store *SQLStore) UpdateOrderTx(ctx context.Context, arg UpdateOrderTxParams) (UpdateOrderTxResult, error) {
	var result UpdateOrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		if (len(arg.ProductID) != len(arg.Size)) || (len(arg.ProductID) != len(arg.Quantity)) || (len(arg.Size) != len(arg.Quantity)) {
			err := errors.New("the length of one or more of the properties is not equal")
			return err
		}

		user, err := q.GetUser(ctx, arg.Username)
		if err != nil {
			return err
		}

		result.UserOrder = newUserResponse(user)

		province, err := q.GetProvince(ctx, arg.Province)
		if err != nil {
			return err
		}

		argUpdateOrder := UpdateOrderParams{
			BookingID: arg.BookingID,
			Address:   arg.Address,
			Province:  province.ID,
		}
		order, err := q.UpdateOrder(ctx, argUpdateOrder)
		if err != nil {
			return err
		}

		var discount float64
		if order.PromotionID != "none" {
			promotion, err := q.GetPromotion(ctx, order.PromotionID)
			if err != nil {
				if err == sql.ErrNoRows {
					err = fmt.Errorf("promotion code doesn't exist")
					return err
				} else {
					return err
				}
			}

			if time.Now().After(promotion.EndDate) {
				err = fmt.Errorf("promotion code has expired")
				return err
			}
			discount = promotion.DiscountPercent / 100
		} else {
			discount = 0
		}

		items, err := q.ListItemsOrderByBookingID(ctx, arg.BookingID)
		if err != nil {
			return err
		}
		for _, item := range items {
			store, err := q.GetStore(ctx, GetStoreParams{
				ProductID: item.ProductID,
				Size:      item.Size,
			})
			if err != nil {
				return err
			}

			_, err = q.UpdateStore(ctx, UpdateStoreParams{
				ProductID: item.ProductID,
				Size:      item.Size,
				Quantity:  item.Quantity + store.Quantity,
			})
			if err != nil {
				return err
			}
		}
		err = q.DeleteItemsOrderByBookingID(ctx, arg.BookingID)
		if err != nil {
			return err
		}

		var amount float64
		amount = 0
		for i, productID := range arg.ProductID {
			size := arg.Size[i]
			quantity := arg.Quantity[i]

			product, err := q.GetProduct(ctx, productID)
			if err != nil {
				return err
			}

			store, err := q.GetStore(ctx, GetStoreParams{
				ProductID: product.ID,
				Size:      size,
			})
			if err != nil {
				return err
			}

			if quantity > int64(store.Quantity) {
				err := errors.New("quantity is not enough")
				return err
			}

			argItemOrder := CreateItemsOrderParams{
				BookingID: arg.BookingID,
				ProductID: productID,
				Quantity:  int32(quantity),
				Size:      size,
				Price:     product.Price * float64(quantity),
			}

			amount = amount + argItemOrder.Price

			_, err = q.CreateItemsOrder(ctx, argItemOrder)
			if err != nil {
				return err
			}

			_, err = q.UpdateStore(ctx, UpdateStoreParams{
				ProductID: store.ProductID,
				Size:      store.Size,
				Quantity:  store.Quantity - argItemOrder.Quantity,
			})
			if err != nil {
				return err
			}
		}

		amount = (amount + order.Tax*amount) * (1 - discount)

		result.Order, err = q.UpdateAmountOfOrder(ctx, UpdateAmountOfOrderParams{
			BookingID: order.BookingID,
			Amount:    amount,
		})
		if err != nil {
			return err
		}

		result.ProductOrdered, err = q.ListItemsOrderByBookingID(ctx, order.BookingID)
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

type CancelOrderParams struct {
	BookingID   string `json:"booking_id"`
	UserBooking string `json:"user_booking"`
}

func (store *SQLStore) CancelOrderTx(ctx context.Context, arg CancelOrderParams) (string, error) {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		order, err := q.GetOrder(ctx, arg.BookingID)
		if err != nil {
			return err
		}

		items, err := q.ListItemsOrderByBookingID(ctx, order.BookingID)
		if err != nil {
			return err
		}
		for _, item := range items {
			store, err := q.GetStore(ctx, GetStoreParams{
				ProductID: item.ProductID,
				Size:      item.Size,
			})
			if err != nil {
				return err
			}

			_, err = q.UpdateStore(ctx, UpdateStoreParams{
				ProductID: item.ProductID,
				Size:      item.Size,
				Quantity:  item.Quantity + store.Quantity,
			})
			if err != nil {
				return err
			}
		}
		err = q.DeleteItemsOrderByBookingID(ctx, order.BookingID)
		if err != nil {
			return err
		}

		err = q.DeleteOrder(ctx, order.BookingID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return "Cancelling booking successfully", err
}
