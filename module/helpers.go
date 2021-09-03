package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	jsoniter "github.com/json-iterator/go"
	"github.com/ProjectAthenaa/go-credit-card"
	"regexp"
)

var (
	monthregex = regexp.MustCompile("^0")
	json  = jsoniter.ConfigFastest
)

type PaymentResponse struct {
	FieldErrors []struct {
		DwfrmBillingCreditCardFieldsCardNumber string `json:"dwfrm_billing_creditCardFields_cardNumber"`
	} `json:"fieldErrors"`
}


type ConfirmationStruct struct {
	Error       bool   `json:"error"`
	OrderID     string `json:"orderID"`
	OrderToken  string `json:"orderToken"`
}


func (tk *Task) GenerateDefaultHeaders(referrer string) http.Headers {
	return http.Headers{
		`user-agent`:         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"},
		`accept`:             {`application/json`},
		`accept-encoding`:    {`gzip, deflate, br`},
		`accept-language`:    {`en-us`},
		`content-type`:       {`application/json`},
		`sec-ch-ua`:          {`"Chromium";v="91", " Not A;Brand";v="99", "Google Chrome";v="91"`},
		`sec-ch-ua-mobile`:   {`?0`},
		`Sec-Fetch-Site`:     {`same-site`},
		`Sec-Fetch-Dest`:     {`empty`},
		`Sec-Fetch-Mode`:     {`cors`},
		`x-application-name`: {`web`},
		`referer`:            {referrer},
		`origin`:             {`https://www.target.com`},
		`Pragma`:             {`no-cache`},
		`Cache-Control`:      {`no-cache`},
		`Connection`:         {`keep-alive`},
	}
}

func (tk *Task) getShapeHeaders() (map[string]string, error){
	tk.SetStatus(module.STATUS_GENERATING_COOKIES)
	headers, err := shapeClient.GenHeaders(tk.Ctx, &shape.Site{Value: shape.SITE_NEWBALANCE})
	if err != nil{
		return nil, err
	}
	return headers.Values, nil
}

func (tk *Task) formatPhone() string{
	return fmt.Sprintf("(%s) %s-%s",tk.Data.Profile.Shipping.PhoneNumber[0:3],tk.Data.Profile.Shipping.PhoneNumber[3:6],tk.Data.Profile.Shipping.PhoneNumber[6:])
}
func (tk *Task) cardType()string{
	return creditcard.Card{Number: tk.Data.Profile.Billing.Number, Cvv: tk.Data.Profile.Billing.CVV, Month: tk.Data.Profile.Billing.ExpirationMonth, Year: "20" + tk.Data.Profile.Billing.ExpirationYear}.Company.Long
}