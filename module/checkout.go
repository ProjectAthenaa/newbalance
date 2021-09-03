package module

import (
	"fmt"
	module "github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"net/url"
	"strings"
)

func (tk *Task) BeginCheckout() {
	checkouthtmlreq, err := tk.NewRequest("GET", "https://www.newbalance.com/cart//", nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}
	checkouthtmlreq.Headers = tk.GenerateDefaultHeaders("https://www.newbalance.com/checkout-begin/?stage=payment")

	shapeheaders, err := tk.getShapeHeaders()
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not retrieve shape headers")
		tk.Stop()
		return
	}
	for k, v := range shapeheaders {
		checkouthtmlreq.Headers[k] = []string{v}
	}

	checkouthtmlres, err := tk.Do(checkouthtmlreq)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	tk.productLineItemUUID = sonic.GrabValueFromHTMLName("productLineItemUUID", &checkouthtmlres.Body)
	tk.originalShipmentUUID = sonic.GrabValueFromHTMLName("originalShipmentUUID", &checkouthtmlres.Body)
	tk.shipmentUUID = sonic.GrabValueFromHTMLName("shipmentUUID", &checkouthtmlres.Body)
	tk.csrf_token = sonic.GrabValueFromHTMLName("csrf_token", &checkouthtmlres.Body)

}

func (tk *Task) Shipping() {
	tk.SetStatus(module.STATUS_CHECKING_OUT, "submitting shipping")

	var addrline2 string
	if tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != nil {
		addrline2 = *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2
	}

	poststring := url.QueryEscape(fmt.Sprintf(
		"productLineItemUUID=%s&"+
			"originalShipmentUUID=%s&"+
			"shipmentUUID=%s&"+
			"zipCodeErrorMsg=Please enter a valid Zip/Postal code&"+
			"dwfrm_shipping_shippingAddress_addressFields_country=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_firstName=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_lastName=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_address1=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_address2=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_city=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_states_stateCode=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_postalCode=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_phone=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_email=%s&"+
			"dwfrm_shipping_shippingAddress_addressFields_addtoemaillist=true&"+
			"csrf_token=%s&"+
			"saveShippingAddr=false",
		tk.productLineItemUUID,
		tk.originalShipmentUUID,
		tk.shipmentUUID,
		tk.Data.Profile.Shipping.ShippingAddress.Country,
		tk.Data.Profile.Shipping.FirstName,
		tk.Data.Profile.Shipping.LastName,
		tk.Data.Profile.Shipping.ShippingAddress.AddressLine,
		addrline2,
		tk.Data.Profile.Shipping.ShippingAddress.City,
		tk.Data.Profile.Shipping.ShippingAddress.StateCode,
		tk.Data.Profile.Shipping.ShippingAddress.ZIP,
		tk.formatPhone(),
		tk.Data.Profile.Email,
		tk.csrf_token,
	))

	shippingpostreq, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/CheckoutShippingServices-SubmitShipping", []byte(poststring))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}
	shippingpostreq.Headers = tk.GenerateDefaultHeaders("https://www.newbalance.com/checkout-begin/?stage=payment")

	shippingres, err := tk.Do(shippingpostreq)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	if shippingres.StatusCode >= 300 {
		tk.SetStatus(module.STATUS_CHECKOUT_ERROR, "couldnt submit shipping")
		tk.Stop()
		return
	}

	tk.SetStatus(module.STATUS_CHECKING_OUT, "submitted shipping")
}

func (tk *Task) Payment() {
	tk.SetStatus(module.STATUS_CHECKING_OUT, "submitting payment")

	var paymentform string
	if !tk.Data.Profile.Shipping.BillingIsShipping {
		var addrline2 string
		if tk.Data.Profile.Shipping.BillingAddress.AddressLine2 != nil {
			addrline2 = *tk.Data.Profile.Shipping.BillingAddress.AddressLine2
		}

		paymentform = url.QueryEscape(fmt.Sprintf(`csrf_token=%s&`+
			`localizedNewAddressTitle=New+Address&`+
			`dwfrm_billing_paymentMethod=CREDIT_CARD&`+
			`dwfrm_billing_creditCardFields_cardNumber=%s&`+
			`dwfrm_billing_creditCardFields_expirationMonth=%s&`+
			`dwfrm_billing_creditCardFields_expirationYear=%s&`+
			`dwfrm_billing_creditCardFields_securityCode=%s&`+
			`dwfrm_billing_creditCardFields_cardType=%s&`+
			`dwfrm_billing_paymentMethod=CREDIT_CARD&`+
			`dwfrm_afterpay_isAfterpayUrl=%2Fon%2Fdemandware.store%2FSites-NBUS-Site%2Fen_US%2FAfterpayRedirect-IsAfterpay&`+
			`dwfrm_afterpay_redirectAfterpayUrl=%2Fon%2Fdemandware.store%2FSites-NBUS-Site%2Fen_US%2FAfterpayRedirect-Redirect&`+
			`addressSelector=new&`+
			`dwfrm_billing_addressFields_country=%s&`+
			`dwfrm_billing_addressFields_firstName=%s&`+
			`dwfrm_billing_addressFields_lastName=%s&`+
			`dwfrm_billing_addressFields_address1=%s&`+
			`dwfrm_billing_addressFields_address2=%s&`+
			`dwfrm_billing_addressFields_city=%s&`+
			`dwfrm_billing_addressFields_states_stateCode=%s&`+
			`dwfrm_billing_addressFields_postalCode=%s&`+
			`dwfrm_billing_addressFields_email=%s&`+
			`dwfrm_billing_addressFields_phone=%s&&`+
			`dwfrm_billing_paymentMethod=CREDIT_CARD&`+
			`dwfrm_billing_creditCardFields_cardNumber=%s&`+
			`dwfrm_billing_creditCardFields_expirationMonth=%s&`+
			`dwfrm_billing_creditCardFields_expirationYear=%s&`+
			`dwfrm_billing_creditCardFields_securityCode=%s&`+
			`dwfrm_billing_creditCardFields_cardType=%s&`+
			`addressId=new&`+
			`saveBillingAddr=false`,
			tk.csrf_token,
			tk.Data.Profile.Billing.Number,
			monthregex.ReplaceAllString(tk.Data.Profile.Billing.ExpirationMonth, ""),
			"20"+tk.Data.Profile.Billing.ExpirationYear,
			tk.Data.Profile.Billing.CVV,
			tk.cardType(),
			tk.Data.Profile.Shipping.BillingAddress.Country,
			tk.Data.Profile.Shipping.FirstName,
			tk.Data.Profile.Shipping.LastName,
			tk.Data.Profile.Shipping.BillingAddress.AddressLine,
			addrline2,
			tk.Data.Profile.Shipping.BillingAddress.City,
			tk.Data.Profile.Shipping.BillingAddress.StateCode,
			tk.Data.Profile.Shipping.BillingAddress.ZIP,
			tk.Data.Profile.Email,
			tk.formatPhone(),
			tk.Data.Profile.Billing.Number,
			monthregex.ReplaceAllString(tk.Data.Profile.Billing.ExpirationMonth, ""),
			"20"+tk.Data.Profile.Billing.ExpirationYear,
			tk.Data.Profile.Billing.CVV,
			tk.cardType(),
		))
	} else {
		var addrline2 string
		if tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != nil {
			addrline2 = *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2
		}
		paymentform = url.QueryEscape(fmt.Sprintf(`csrf_token=%s&`+
			`localizedNewAddressTitle=New+Address&`+
			`dwfrm_billing_paymentMethod=CREDIT_CARD&`+
			`dwfrm_billing_creditCardFields_cardNumber=%s&`+
			`dwfrm_billing_creditCardFields_expirationMonth=%s&`+
			`dwfrm_billing_creditCardFields_expirationYear=%s&`+
			`dwfrm_billing_creditCardFields_securityCode=%s&`+
			`dwfrm_billing_creditCardFields_cardType=%s&`+
			`dwfrm_billing_paymentMethod=CREDIT_CARD&`+
			`dwfrm_afterpay_isAfterpayUrl=%2Fon%2Fdemandware.store%2FSites-NBUS-Site%2Fen_US%2FAfterpayRedirect-IsAfterpay&`+
			`dwfrm_afterpay_redirectAfterpayUrl=%2Fon%2Fdemandware.store%2FSites-NBUS-Site%2Fen_US%2FAfterpayRedirect-Redirect&`+
			`addressSelector=%s&`+
			`dwfrm_billing_addressFields_country=%s&`+
			`dwfrm_billing_addressFields_firstName=%s&`+
			`dwfrm_billing_addressFields_lastName=%s&`+
			`dwfrm_billing_addressFields_address1=%s&`+
			`dwfrm_billing_addressFields_address2=%s&`+
			`dwfrm_billing_addressFields_city=%s&`+
			`dwfrm_billing_addressFields_states_stateCode=%s&`+
			`dwfrm_billing_addressFields_postalCode=%s&`+
			`dwfrm_billing_addressFields_email=%s&`+
			`dwfrm_billing_addressFields_phone=%s&&`+
			`dwfrm_billing_paymentMethod=CREDIT_CARD&`+
			`dwfrm_billing_creditCardFields_cardNumber=%s&`+
			`dwfrm_billing_creditCardFields_expirationMonth=%s&`+
			`dwfrm_billing_creditCardFields_expirationYear=%s&`+
			`dwfrm_billing_creditCardFields_securityCode=%s&`+
			`dwfrm_billing_creditCardFields_cardType=%s&`+
			`addressId=%s&`+
			`saveBillingAddr=false`,
			tk.csrf_token,
			tk.Data.Profile.Billing.Number,
			monthregex.ReplaceAllString(tk.Data.Profile.Billing.ExpirationMonth, ""),
			"20"+tk.Data.Profile.Billing.ExpirationYear,
			tk.Data.Profile.Billing.CVV,
			tk.cardType(),
			tk.originalShipmentUUID,
			tk.Data.Profile.Shipping.ShippingAddress.Country,
			tk.Data.Profile.Shipping.FirstName,
			tk.Data.Profile.Shipping.LastName,
			tk.Data.Profile.Shipping.ShippingAddress.AddressLine,
			addrline2,
			tk.Data.Profile.Shipping.ShippingAddress.City,
			tk.Data.Profile.Shipping.ShippingAddress.StateCode,
			tk.Data.Profile.Shipping.ShippingAddress.ZIP,
			tk.Data.Profile.Email,
			tk.formatPhone(),
			tk.Data.Profile.Billing.Number,
			monthregex.ReplaceAllString(tk.Data.Profile.Billing.ExpirationMonth, ""),
			"20"+tk.Data.Profile.Billing.ExpirationYear,
			tk.Data.Profile.Billing.CVV,
			tk.cardType(),
			tk.originalShipmentUUID,
		))
	}

	paymentreq, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/CheckoutServices-SubmitPayment", []byte(paymentform))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	paymentreq.Headers = tk.GenerateDefaultHeaders("https://www.newbalance.com/checkout-begin/?stage=payment")


	shapeheaders, err := tk.getShapeHeaders()
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not retrieve shape headers")
		tk.Stop()
		return
	}
	for k, v := range shapeheaders {
		paymentreq.Headers[k] = []string{v}
	}

	paymentres, err := tk.Do(paymentreq)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	var paymentstruct *PaymentResponse
	if err = json.Unmarshal(paymentres.Body, &paymentstruct); err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	if paymentstruct.FieldErrors != nil {
		tk.SetStatus(module.STATUS_ERROR, paymentstruct.FieldErrors[0].DwfrmBillingCreditCardFieldsCardNumber)
		tk.Stop()
		return
	}

	tk.SetStatus(module.STATUS_CHECKING_OUT, "submitting final info")
}

func (tk *Task) PlaceOrder() {
	confirmationreq, err := tk.NewRequest("POST", "https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/CheckoutServices-PlaceOrder?termsconditions=undefined&DFReferenceId=", nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}
	confirmationreq.Headers = tk.GenerateDefaultHeaders("https://www.newbalance.com/checkout-begin/?stage=payment")

	res, err := tk.Do(confirmationreq)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	var confirmationresponse *ConfirmationStruct
	if err = json.Unmarshal(res.Body, &confirmationresponse); err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	if confirmationresponse.Error {
		tk.SetStatus(module.STATUS_CHECKOUT_FAILED, "couldnt review")
		tk.Restart()
		return
	}

	confirmationpagereq, err := tk.NewRequest("GET", fmt.Sprintf("https://www.newbalance.com/on/demandware.store/Sites-NBUS-Site/en_US/Order-Confirm?ID=%s&token=%s", confirmationresponse.OrderID, confirmationresponse.OrderToken), nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	confirmationpageres, err := tk.Do(confirmationpagereq)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	if !strings.Contains(string(confirmationpageres.Body), "Thank you for your order") {
		tk.SetStatus(module.STATUS_CHECKOUT_FAILED, "couldnt confirm")
		tk.Restart()
		return
	}

	tk.SetStatus(module.STATUS_CHECKED_OUT, "checked out")
}
