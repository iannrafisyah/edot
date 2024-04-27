package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"Edot/models"
	"Edot/modules/cart/dto"
	"Edot/modules/cart/repository"
	productController "Edot/modules/product/controller"
	productDto "Edot/modules/product/dto"
	"Edot/packages/logger"
	"Edot/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	ICartController interface {
		FindAll(ctx context.Context, reqData *dto.FindAllRequest, tx *gorm.DB) ([]*dto.FindAllResponse, error)
		AddCart(ctx context.Context, reqData *dto.AddCartRequest, tx *gorm.DB) error
		Delete(ctx context.Context, reqData *dto.DeleteRequest, tx *gorm.DB) error
	}

	CartController struct {
		fx.In
		Logger            *logger.Logger
		CartRepository    repository.ICartInterface
		ProductController productController.IProductController
	}
)

// NewController :
func NewController(cartController CartController) ICartController {
	return &cartController
}

// AddCart :
func (r *CartController) AddCart(ctx context.Context, reqData *dto.AddCartRequest, tx *gorm.DB) error {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	if err := r.CartRepository.Delete(ctx, &models.Cart{
		UserID: reqData.ContextUserID,
	}, tx); err != nil {
		return err
	}

	for _, item := range reqData.Items {
		fetchProduct, err := r.ProductController.FindByID(ctx, &productDto.FindByIDRequest{
			ContextUserID: reqData.ContextUserID,
			ID:            item.ProductID,
			WarehouseID:   item.WarehouseID,
		}, tx)
		if err != nil {
			r.Logger.Error(err)
			return err
		}

		if fetchProduct.Warehouses == nil {
			return utilities.ErrorRequest(
				fmt.Errorf(utilities.DataNotFound, "Warehouse"),
				http.StatusBadRequest,
			)
		}

		totalStockPendingPayment, err := r.CartRepository.TotalQtyWithTrxPending(ctx, &models.Cart{
			ProductID:   item.ProductID,
			WarehouseID: item.WarehouseID,
		}, tx)
		if err != nil {
			r.Logger.Error(err)
			return utilities.ErrorRequest(
				errors.New(utilities.InternalServiceError),
				http.StatusInternalServerError,
			)
		}

		if totalStockPendingPayment == nil {
			totalStockPendingPayment = utilities.Int64Pointer(0)
		}

		stockAvailable := fetchProduct.StockAvailable - int(*totalStockPendingPayment)
		if stockAvailable < item.Qty {
			return utilities.ErrorRequest(
				fmt.Errorf(utilities.InsufficientStock, fetchProduct.Name),
				http.StatusBadRequest,
			)
		}

		if _, err := r.CartRepository.Create(ctx, &models.Cart{
			Name:          fetchProduct.Name,
			Price:         fetchProduct.Price,
			Qty:           item.Qty,
			UserID:        reqData.ContextUserID,
			ProductID:     fetchProduct.ID,
			WarehouseID:   item.WarehouseID,
			ToWarehouseID: item.ToWarehouseID,
		}, tx); err != nil {
			r.Logger.Error(err)
			return utilities.ErrorRequest(
				errors.New(utilities.InternalServiceError),
				http.StatusInternalServerError,
			)
		}
	}

	return nil
}

// FindAll :
func (r *CartController) FindAll(ctx context.Context, reqData *dto.FindAllRequest, tx *gorm.DB) ([]*dto.FindAllResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchAllCart, err := r.CartRepository.FindByUserID(ctx, &models.Cart{
		UserID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	if len(fetchAllCart) <= 0 {
		return nil, utilities.ErrorRequest(
			fmt.Errorf(utilities.DataNotFound, "Cart"),
			http.StatusOK,
		)
	}

	response := []*dto.FindAllResponse{}

	for _, cart := range fetchAllCart {
		response = append(response, &dto.FindAllResponse{
			ProductID:     cart.ProductID,
			Name:          cart.Name,
			Price:         cart.Price,
			Qty:           cart.Qty,
			WarehouseID:   cart.WarehouseID,
			ToWarehouseID: cart.ToWarehouseID,
		})
	}

	return response, nil
}

// Delete :
func (r *CartController) Delete(ctx context.Context, reqData *dto.DeleteRequest, tx *gorm.DB) error {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	if err := r.CartRepository.Delete(ctx, &models.Cart{
		UserID: reqData.ContextUserID,
	}, tx); err != nil {
		return err
	}

	return nil
}
