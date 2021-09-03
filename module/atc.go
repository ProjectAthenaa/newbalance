package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
)

func (tk *Task) ATC() {
	req, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/Cart-AddProduct", []byte(fmt.Sprintf("pid=%s&quantity=1&estimatedDelivery=&options=%%5B%%5D", tk.VariantId)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create atc req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf(`https://www.newbalance.com/pd/~/%s.html`, tk.PID))

	shapeheaders, err := tk.getShapeHeaders()
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not retrieve shape headers")
		tk.Stop()
		return
	}
	for k, v := range shapeheaders {
		req.Headers[k] = []string{v}
	}

	tk.SetStatus(module.STATUS_ADDING_TO_CART, "adding to cart")
	tk.Do(req)
	tk.SetStatus(module.STATUS_WAITING_FOR_CHECKOUT, "moving to checkout")
}
