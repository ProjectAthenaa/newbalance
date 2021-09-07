package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
)

func (tk *Task) ATC() {
	req, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/Cart-AddProduct", []byte(`pid=194768477287&quantity=1&estimatedDelivery=Estimated+delivery+-+2-5+business+days+once+shipped&options=%5B%5D`))
	//req, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/Cart-AddProduct", []byte(fmt.Sprintf("pid=%s&quantity=1&estimatedDelivery=Estimated+delivery+-+2-5+business+days+once+shipped&options=%%5B%%5D", tk.VariantId)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create atc req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf(`https://www.newbalance.com/pd/~/%s.html`, tk.PID))
	//
	//cookiejar.ReleaseCookieJar(tk.FastClient.Jar)
	//tk.FastClient.Jar = nil

	//req.Headers[`x-dtpc`] = []string{`2$500472369_531h28vPUNGQCPHUHIIUBHMKQIJELBUOKKUAIFM-0e2`}
	//req.Headers["Accept"] = []string{"*/*"}

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
