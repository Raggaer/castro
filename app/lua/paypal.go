package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/gopaypal"
	"github.com/yuin/gopher-lua"
	"strconv"
)

// paypal application client
var client gopaypal.Client

// CreatePaypalClient creates the paypal application client
func CreatePaypalClient(sandbox bool) {
	// Create application client for live settings
	if !sandbox {

		client = gopaypal.NewClient(
			util.Config.PayPal.PublicKey,
			util.Config.PayPal.SecretKey,
			gopaypal.LiveURL,
		)

		return
	}

	// Create application client for sandbox settings
	client = gopaypal.NewClient(
		util.Config.PayPal.PublicKey,
		util.Config.PayPal.SecretKey,
		gopaypal.SandBoxURL,
	)
}

// SetPayPalMetaTable sets the paypal metatable of the given state
func SetPayPalMetaTable(luaState *lua.LState) {
	// Create and set the paypal metatable
	paypalMetaTable := luaState.NewTypeMetatable(PayPalMetaTableName)
	luaState.SetGlobal(PayPalMetaTableName, paypalMetaTable)

	// Set all mail metatable functions
	luaState.SetFuncs(paypalMetaTable, paypalMethods)
}

// CreatePaypalPayment creates a paypal payment returning the payment URL
func CreatePaypalPayment(L *lua.LState) int {
	// Get payment price
	price := L.ToInt(4)

	// HTTP mode
	mode := "http"

	if util.Config.SSL.Enabled {
		mode = "https"
	}

	// Create paypal payment
	payment := gopaypal.Payment{
		Intent: "sale",
		Payer: gopaypal.Payer{
			PaymentMethod: "paypal",
		},
		Transactions: []gopaypal.Transaction{
			{
				Amount: gopaypal.Amount{
					Total:    strconv.Itoa(price),
					Currency: util.Config.PayPal.Currency,
					Details: gopaypal.Details{
						SubTotal: strconv.Itoa(price),
					},
				},
				Description: L.ToString(2),
				Custom:      "test",
				ItemList: gopaypal.ItemList{
					Items: []gopaypal.Item{
						{
							Name:     L.ToString(3),
							Price:    strconv.Itoa(price),
							Currency: util.Config.PayPal.Currency,
							Quantity: 1,
						},
					},
				},
			},
		},
		RedirectURL: gopaypal.RedirectURL{
			ReturnURL: fmt.Sprintf("%v://%v:%v/%v", mode, util.Config.URL, util.Config.Port, "subtopic/shop/paypal"),
			CancelURL: fmt.Sprintf("%v://%v:%v/%v", mode, util.Config.URL, util.Config.Port, "subtopic/shop/paypal"),
		},
	}

	// Create paypal payment
	paymentResponse, err := client.CreatePayment(payment)

	if err != nil {
		L.RaiseError("Cannot create paypal payment: %v", err)
		return 0
	}

	// Loop payment links
	for _, link := range paymentResponse.Links {

		// Push approval URL
		if link.Rel == "approval_url" {

			L.Push(lua.LString(link.Href))

			return 1
		}
	}

	L.RaiseError("Cannot find approval_url payment link")
	return 0
}
