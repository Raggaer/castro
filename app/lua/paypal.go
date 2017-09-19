package lua

import (
	"strconv"

	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/gopaypal"
	"github.com/yuin/gopher-lua"
)

// paypal application client
var client gopaypal.Client

// CreatePaypalClient creates the paypal application client
func CreatePaypalClient(sandbox bool) {
	// Create application client for live settings
	if !sandbox {

		client = gopaypal.NewClient(
			util.Config.Configuration.PayPal.PublicKey,
			util.Config.Configuration.PayPal.SecretKey,
			gopaypal.LiveURL,
		)

		return
	}

	// Create application client for sandbox settings
	client = gopaypal.NewClient(
		util.Config.Configuration.PayPal.PublicKey,
		util.Config.Configuration.PayPal.SecretKey,
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
	price := L.ToInt(3)

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
					Currency: util.Config.Configuration.PayPal.Currency,
					Details: gopaypal.Details{
						SubTotal: strconv.Itoa(price),
					},
				},
				Description: L.ToString(2),
				Custom:      L.ToString(4),
				ItemList: gopaypal.ItemList{
					Items: []gopaypal.Item{
						{
							Name:     L.ToString(2),
							Price:    strconv.Itoa(price),
							Currency: util.Config.Configuration.PayPal.Currency,
							Quantity: 1,
						},
					},
				},
			},
		},
		RedirectURL: gopaypal.RedirectURL{
			ReturnURL: L.ToString(6),
			CancelURL: L.ToString(5),
		},
	}

	// Create paypal payment
	info, err := client.CreatePayment(payment)

	if err != nil {
		L.RaiseError("Cannot create paypal payment: %v", err)
		return 0
	}

	// Data table
	tbl := L.NewTable()

	// Set payment fields
	tbl.RawSetString("State", lua.LString(info.State))
	tbl.RawSetString("Custom", lua.LString(info.Transactions[0].Custom))
	tbl.RawSetString("Price", lua.LNumber(price))
	tbl.RawSetString("Name", lua.LString(info.Transactions[0].Description))
	tbl.RawSetString("PaymentID", lua.LString(info.ID))
	tbl.RawSetString("PayerStatus", lua.LString(info.Payer.Status))

	// Loop payment links
	for _, link := range info.Links {

		// If link is approval add to main table
		if link.Rel == "approval_url" {

			// Set approval link
			tbl.RawSetString("Link", lua.LString(link.Href))

			break
		}
	}

	// Push data table
	L.Push(tbl)

	return 1
}

// GetPaypalPayment gets a paypal approved payment
func GetPaypalPayment(L *lua.LState) int {
	// Get payment identifier
	id := L.Get(2)

	// Check valid identifier
	if id.Type() != lua.LTString {
		L.ArgError(1, "Invalid payment identifier type. Expected string")
		return 0
	}

	// Get payment information
	info, err := client.PaymentInformation(id.String())

	if err != nil {

		// Log if development mode
		if util.Config.Configuration.IsDev() {
			util.Logger.Logger.Errorf("Cannot get paypal payment information: %v", err)
		}

		L.Push(lua.LNil)
		return 1
	}

	// Validate approved payment
	if len(info.Transactions) > 1 {
		L.Push(lua.LNil)
		return 1
	}

	// Result table
	tbl := L.NewTable()

	// Get payment price
	price, err := strconv.ParseFloat(info.Transactions[0].Amount.Total, 10)

	if err != nil {
		L.RaiseError("Cannot get payment price: %v", err)
		return 0
	}

	// Set payment fields
	tbl.RawSetString("State", lua.LString(info.State))
	tbl.RawSetString("Custom", lua.LString(info.Transactions[0].Custom))
	tbl.RawSetString("Price", lua.LNumber(price))
	tbl.RawSetString("Name", lua.LString(info.Transactions[0].Description))
	tbl.RawSetString("PaymentID", lua.LString(id.String()))
	tbl.RawSetString("PayerID", lua.LString(info.Payer.Info.ID))
	tbl.RawSetString("PayerStatus", lua.LString(info.Payer.Status))

	// Push result table
	L.Push(tbl)

	return 1
}

// ExecutePaypalPayment executes the given paypal payment
func ExecutePaypalPayment(L *lua.LState) int {
	// Get payment identifier
	id := L.Get(2)

	// Check valid identifier
	if id.Type() != lua.LTString {
		L.ArgError(1, "Invalid payment identifier type. Expected string")
		return 0
	}

	// Get payer identifier
	payerID := L.Get(3)

	// Check valid identifier
	if payerID.Type() != lua.LTString {
		L.ArgError(2, "Invalid payer identifier type. Expected string")
	}

	// Execute paypal payment
	_, err := client.ExecutePayment(id.String(), payerID.String())

	if err != nil {
		// Log if development mode
		if util.Config.Configuration.IsDev() {
			util.Logger.Logger.Errorf("Cannot execute paypal payment: %v", err)
		}

		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(true))
	return 1
}
