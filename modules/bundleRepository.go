package modules

import (
	cart "Edot/modules/cart/repository"
	product "Edot/modules/product/repository"
	stock "Edot/modules/stock/repository"
	transaction "Edot/modules/transaction/repository"
	user "Edot/modules/user/repository"

	"go.uber.org/fx"
)

// AppRepository :
var AppRepository = fx.Options(
	fx.Provide(user.NewRepository),
	fx.Provide(product.NewRepository),
	fx.Provide(cart.NewRepository),
	fx.Provide(transaction.NewRepository),
	fx.Provide(stock.NewRepository),
)
