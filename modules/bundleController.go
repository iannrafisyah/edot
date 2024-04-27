package modules

import (
	auth "Edot/modules/auth/controller"
	cart "Edot/modules/cart/controller"
	product "Edot/modules/product/controller"
	stock "Edot/modules/stock/controller"
	transaction "Edot/modules/transaction/controller"
	user "Edot/modules/user/controller"
	worker "Edot/modules/worker/controller"

	"go.uber.org/fx"
)

// AppController :
var AppController = fx.Options(
	fx.Provide(user.NewController),
	fx.Provide(auth.NewController),
	fx.Provide(product.NewController),
	fx.Provide(cart.NewController),
	fx.Provide(transaction.NewController),
	fx.Provide(stock.NewController),
	fx.Provide(worker.NewController),
)
