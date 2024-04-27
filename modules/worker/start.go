package worker

import (
	worker "Edot/modules/worker/controller"
	"context"
)

func Start(worker worker.IWorkerController) {
	ctx := context.Background()
	worker.CancelTransaction(ctx, "*/1 * * * *")
}
