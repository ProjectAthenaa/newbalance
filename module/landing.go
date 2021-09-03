package module

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
)

func (tk *Task) GetProductPage(){
	req, err := tk.NewRequest("GET", tk.productUrl, nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create atc req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(`https://www.newbalance.com/`)
	_, err = tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not retrieve product page")
		tk.Stop()
		return
	}
}