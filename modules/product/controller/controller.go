package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"Edot/models"
	"Edot/modules/product/dto"
	"Edot/modules/product/repository"
	"Edot/packages/logger"
	"Edot/packages/paginate"
	"Edot/packages/postgres"
	"Edot/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	IProductController interface {
		FindAll(ctx context.Context, reqData *dto.FindAllRequest, tx *gorm.DB) ([]*dto.FindAllResponse, *paginate.Pagination, error)
		FindByID(ctx context.Context, reqData *dto.FindByIDRequest, tx *gorm.DB) (*dto.FindByIDResponse, error)
	}

	ProductController struct {
		fx.In
		Logger            *logger.Logger
		DB                *postgres.DB
		ProductRepository repository.IProductInterface
	}
)

// NewController :
func NewController(productController ProductController) IProductController {
	return &productController
}

// FindAll :
func (r *ProductController) FindAll(ctx context.Context, reqData *dto.FindAllRequest, tx *gorm.DB) ([]*dto.FindAllResponse, *paginate.Pagination, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchAllProduct, paging, err := r.ProductRepository.FindAll(ctx, &models.Product{}, paginate.Pagination{
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

	for _, v := range fetchAllProduct {
		response = append(response, &dto.FindAllResponse{
			ID:    v.ID,
			Name:  v.Name,
			Price: v.Price,
		})
	}

	return response, paging, nil
}

// FindByID :
func (r *ProductController) FindByID(ctx context.Context, reqData *dto.FindByIDRequest, tx *gorm.DB) (*dto.FindByIDResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchProductDetail, err := r.ProductRepository.FindByID(ctx, &models.Product{
		ID: reqData.ID,
		Filter: models.Filter{
			WarehouseID: reqData.WarehouseID,
		},
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.DataNotFound, "product"),
				http.StatusNotFound,
			)
		}
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	var (
		totalStock = int(0)
		warehouses []*dto.FindByIDWarehouseStock
	)

	for _, warehouseStock := range fetchProductDetail.Stocks {
		totalStock += warehouseStock.Latest
		warehouses = append(warehouses, &dto.FindByIDWarehouseStock{
			WarehouseID:   warehouseStock.WarehouseID,
			WarehouseName: warehouseStock.Warehouse.Name,
			Latest:        warehouseStock.Latest,
			Previous:      warehouseStock.Previous,
			Qty:           warehouseStock.Qty,
		})
	}

	response := dto.FindByIDResponse{
		ID:             fetchProductDetail.ID,
		Name:           fetchProductDetail.Name,
		Price:          fetchProductDetail.Price,
		StockAvailable: totalStock,
		Warehouses:     warehouses,
	}

	return &response, nil
}
