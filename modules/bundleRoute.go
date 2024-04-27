package modules

import (
	auth "Edot/modules/auth/route"
	cart "Edot/modules/cart/route"
	product "Edot/modules/product/route"
	stock "Edot/modules/stock/route"
	transaction "Edot/modules/transaction/route"
	user "Edot/modules/user/route"
	worker "Edot/modules/worker"

	"go.uber.org/fx"
)

// AppRoute :
var AppRoute = fx.Options(
	fx.Invoke(user.NewRoute),
	fx.Invoke(auth.NewRoute),
	fx.Invoke(product.NewRoute),
	fx.Invoke(cart.NewRoute),
	fx.Invoke(transaction.NewRoute),
	fx.Invoke(stock.NewRoute),
	fx.Invoke(worker.Start),
)
