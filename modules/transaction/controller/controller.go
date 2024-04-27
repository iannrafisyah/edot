package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"Edot/config"
	"Edot/models"
	cartController "Edot/modules/cart/controller"
	cartDto "Edot/modules/cart/dto"
	productController "Edot/modules/product/controller"
	productDto "Edot/modules/product/dto"
	stockController "Edot/modules/stock/controller"
	stockDto "Edot/modules/stock/dto"
	"Edot/modules/transaction/dto"
	"Edot/modules/transaction/repository"
	"Edot/packages/logger"
	"Edot/packages/paginate"
	"Edot/utilities"

	"github.com/brianvoe/gofakeit/v7"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	ITransactionController interface {
		FindByID(ctx context.Context, reqData *dto.FindByIDRequest, tx *gorm.DB) (*dto.FindByIDResponse, error)
		FindAll(ctx context.Context, reqData *dto.FindAllRequest, tx *gorm.DB) ([]*dto.FindAllResponse, *paginate.Pagination, error)
		Create(ctx context.Context, reqData *dto.CreateRequest, tx *gorm.DB) (*dto.FindByIDResponse, error)
		Payment(ctx context.Context, reqData *dto.PaymentRequest, tx *gorm.DB) (*dto.FindByIDResponse, error)
	}

	TransactionController struct {
		fx.In
		TransactionRepository repository.ITransactionInterface
		ProductController     productController.IProductController
		CartController        cartController.ICartController
		StockController       stockController.IStockController
		Logger                *logger.Logger
	}
)

// NewController :
func NewController(transactionController TransactionController) ITransactionController {
	return &transactionController
}

// Create :
func (r *TransactionController) Create(ctx context.Context, reqData *dto.CreateRequest, tx *gorm.DB) (*dto.FindByIDResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	var (
		totalAmount         float64
		tax                 = config.Get().Tax
		transactionProducts []*models.TransactionProduct
	)

	fetchAllCart, err := r.CartController.FindAll(ctx, &cartDto.FindAllRequest{
		ContextUserID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	for _, item := range fetchAllCart {
		fetchProduct, err := r.ProductController.FindByID(ctx, &productDto.FindByIDRequest{
			ContextUserID: reqData.ContextUserID,
			ID:            item.ProductID,
			WarehouseID:   item.WarehouseID,
		}, tx)
		if err != nil {
			r.Logger.Error(err)
			return nil, err
		}

		if fetchProduct.StockAvailable < item.Qty {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.InsufficientStock, fetchProduct.Name),
				http.StatusBadRequest,
			)
		}

		totalAmount += fetchProduct.Price

		transactionProducts = append(transactionProducts, &models.TransactionProduct{
			Name:          item.Name,
			Price:         item.Price,
			Qty:           item.Qty,
			ProductID:     item.ProductID,
			WarehouseID:   item.WarehouseID,
			ToWarehouseID: item.ToWarehouseID,
			Snapshot: models.TransactionProductSnapshot{
				ID:    fetchProduct.ID,
				Name:  fetchProduct.Name,
				Price: fetchProduct.Price,
			},
		})
	}

	taxTotal := totalAmount * tax / 100
	grandTotal := totalAmount + taxTotal

	transactionID, err := r.TransactionRepository.Create(ctx, &models.Transaction{
		Invoice:    gofakeit.UUID(),
		Tax:        taxTotal,
		UserID:     reqData.ContextUserID,
		Amount:     totalAmount,
		Status:     models.TransactionStatusUnpaid,
		Type:       reqData.TransactionType,
		GrandTotal: grandTotal,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	if len(transactionProducts) <= 0 {
		return nil, utilities.ErrorRequest(
			errors.New(utilities.BadRequest),
			http.StatusBadRequest,
		)
	}

	for _, transactionProduct := range transactionProducts {
		transactionProduct.TransactionID = *transactionID
		if _, err := r.TransactionRepository.CreateTransactionProduct(ctx, transactionProduct, tx); err != nil {
			r.Logger.Error(err)
			return nil, utilities.ErrorRequest(
				errors.New(utilities.InternalServiceError),
				http.StatusInternalServerError,
			)
		}
	}

	if err := r.CartController.Delete(ctx, &cartDto.DeleteRequest{
		ContextUserID: reqData.ContextUserID,
	}, tx); err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	fetchTransaction, err := r.FindByID(ctx, &dto.FindByIDRequest{
		ID:            *transactionID,
		ContextUserID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return fetchTransaction, nil
}

// FindByID :
func (r *TransactionController) FindByID(ctx context.Context, reqData *dto.FindByIDRequest, tx *gorm.DB) (*dto.FindByIDResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchTransaction, err := r.TransactionRepository.FindByID(ctx, &models.Transaction{
		ID:     reqData.ID,
		UserID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.DataNotFound, "Transaction"),
				http.StatusNotFound,
			)
		}
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	var transactionProduct []*dto.FindByIDProduct

	for _, product := range fetchTransaction.TransactionProducts {
		transactionProduct = append(transactionProduct, &dto.FindByIDProduct{
			Name:  product.Name,
			Price: product.Price,
			Qty:   product.Qty,
			Total: product.Price * float64(product.Qty),
		})
	}

	return &dto.FindByIDResponse{
		ID:         fetchTransaction.ID,
		Invoice:    fetchTransaction.Invoice,
		Total:      fetchTransaction.Amount,
		Tax:        fetchTransaction.Tax,
		GrandTotal: fetchTransaction.GrandTotal,
		Status:     fetchTransaction.Status.String(),
		Type:       fetchTransaction.Type.String(),
		Products:   transactionProduct,
		CreatedAt:  fetchTransaction.CreatedAt,
	}, nil
}

// FindAll :
func (r *TransactionController) FindAll(ctx context.Context, reqData *dto.FindAllRequest, tx *gorm.DB) ([]*dto.FindAllResponse, *paginate.Pagination, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchAllTransaction, paging, err := r.TransactionRepository.FindAll(ctx, &models.Transaction{
		UserID: reqData.ContextUserID,
		Type:   reqData.Type,
	}, paginate.Pagination{
		Limit: reqData.Limit,
		Page:  reqData.Page,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		return nil, nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	response := []*dto.FindAllResponse{}

	for _, transaction := range fetchAllTransaction {
		response = append(response, &dto.FindAllResponse{
			ID:         transaction.ID,
			Invoice:    transaction.Invoice,
			Total:      transaction.Amount,
			Tax:        transaction.Tax,
			GrandTotal: transaction.GrandTotal,
			Status:     transaction.Status.String(),
			Type:       transaction.Type.String(),
			CreatedAt:  transaction.CreatedAt,
		})
	}

	return response, paging, nil
}

// Payment :
func (r *TransactionController) Payment(ctx context.Context, reqData *dto.PaymentRequest, tx *gorm.DB) (*dto.FindByIDResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchTransaction, err := r.TransactionRepository.FindByID(ctx, &models.Transaction{
		ID:     reqData.ID,
		UserID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.DataNotFound, "Transaction"),
				http.StatusNotFound,
			)
		}
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	if fetchTransaction.Status == models.TransactionStatusPaid {
		return nil, utilities.ErrorRequest(
			errors.New(utilities.TransactionAlreadyPaid),
			http.StatusBadRequest,
		)
	}

	for _, product := range fetchTransaction.TransactionProducts {
		fetchProduct, err := r.ProductController.FindByID(ctx, &productDto.FindByIDRequest{
			ContextUserID: reqData.ContextUserID,
			ID:            product.ProductID,
			WarehouseID:   product.WarehouseID,
		}, tx)
		if err != nil {
			r.Logger.Error(err)
			return nil, err
		}

		if fetchProduct.StockAvailable < product.Qty {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.InsufficientStock, fetchProduct.Name),
				http.StatusBadRequest,
			)
		}

		if err := r.StockController.Create(ctx, &stockDto.AddCartRequest{
			ContextUserID: reqData.ContextUserID,
			Qty:           product.Qty,
			WarehouseID:   product.WarehouseID,
			Operator:      models.StockOperatorDecrement,
			TransactionID: product.TransactionID,
			ProductID:     product.ProductID,
		}, tx); err != nil {
			r.Logger.Error(err)
			return nil, err
		}

		if fetchTransaction.Type == models.TransactionTypeTransferStock &&
			product.ToWarehouseID != nil {
			if err := r.StockController.Create(ctx, &stockDto.AddCartRequest{
				ContextUserID: reqData.ContextUserID,
				Qty:           product.Qty,
				WarehouseID:   *product.ToWarehouseID,
				Operator:      models.StockOperatorIncrement,
				TransactionID: product.TransactionID,
				ProductID:     product.ProductID,
			}, tx); err != nil {
				r.Logger.Error(err)
				return nil, err
			}
		}
	}

	if err := r.TransactionRepository.UpdateStatus(ctx, &models.Transaction{
		ID:     reqData.ID,
		UserID: reqData.ContextUserID,
		Status: models.TransactionStatusPaid,
	}, tx); err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	fetchNewTransaction, err := r.FindByID(ctx, &dto.FindByIDRequest{
		ID:            fetchTransaction.ID,
		ContextUserID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return fetchNewTransaction, nil
}
