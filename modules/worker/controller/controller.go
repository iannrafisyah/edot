package worker

import (
	"Edot/models"
	transactionRepository "Edot/modules/transaction/repository"
	"Edot/packages/logger"
	"Edot/packages/postgres"
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/fx"
)

type (
	IWorkerController interface {
		CancelTransaction(ctx context.Context, duration string)
	}

	WorkerController struct {
		fx.In
		Logger                *logger.Logger
		GoCron                *gocron.Scheduler `optional:"true"`
		DB                    *postgres.DB
		TransactionRepository transactionRepository.ITransactionInterface
	}
)

func NewController(workerController WorkerController) IWorkerController {
	workerController.GoCron = gocron.NewScheduler(time.UTC)
	workerController.GoCron.StartAsync()
	return &workerController
}

func (r *WorkerController) CancelTransaction(ctx context.Context, duration string) {
	r.GoCron.Cron(duration).Do(func() {
		tx := r.DB.Gorm.Begin()
		fetchAllTransaction, err := r.TransactionRepository.FindAllOrderExpired(ctx, &models.Transaction{}, tx)
		if err != nil {
			r.Logger.Error(err)
			tx.Rollback()
			return
		}

		for _, transaction := range fetchAllTransaction {
			newTx := r.DB.Gorm.Begin()
			if err := r.TransactionRepository.UpdateStatus(ctx, &models.Transaction{
				UserID: transaction.UserID,
				ID:     transaction.ID,
				Status: models.TransactionStatusCancel,
			}, tx); err != nil {
				r.Logger.Error(err)
				newTx.Rollback()
				return
			}
			newTx.Commit()
		}
		tx.Commit()
		r.Logger.Infof(`Success process %v data`, len(fetchAllTransaction))
	})
}
