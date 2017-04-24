---
Name: paypal
---

# PayPal metatable

Provides access to the Paypal REST API

- [paypal:createPayment(description, price, custom, cancel_url, return_url)](#createpayment)
- [paypal:paymentInformation(payment_id)](#paymentinformation)
- [paypal:executePayment(payment_id, payer_id)](#executepayment)

# createPayment

Creates a new PayPal payment. You must have PayPal configured on your `config.toml` file. These are the list of mandatory parameters:

- description: payment description text.
- price: payment price using the `config.toml` currency.
- custom: custom data to pass to PayPal.
- cancel_url: where to redirect the user if canceling the payment.
- return_url: where to redirect the user after payment is approved.

```lua
local payment = paypal:createPayment(
    "test", 12, "custom", "www.website.com/cancel", "www.website.com/return"
)
--[[
payment.State = "created"
payment.Custom = "custom"
payment.Price = 12
payment.PaymentID = "PAY-XXXXX-XXXXX"
payment.PayerStatus = "verified"
payment.Link = "www.paypal.com/pay"
]]--
```

The returning table contains these fields:

- State: payment state.
- Custom: your custom passed data.
- Price: payment price.
- PaymentID: payment identifier.
- PayerStatus: account status of the payer. Verified or not.
- Link: payment processing link. Where users can approve the payment.

# paymentInformation

Returns information about the given PayPal payment.

```lua
local payment = paypal:paymentInformation("PAY-XXXXX-XXXX")
--[[
payment.State = "created"
payment.Custom = "custom"
payment.Price = 12
payment.PaymentID = "PAY-XXXXX-XXXXX"
payment.PayerID = "USER-XXXXX"
payment.PayerStatus = "verified"
payment.Link = "www.paypal.com/pay"
]]--
```

The returning table contains these fields:

- State: payment state.
- Custom: your custom passed data.
- Price: payment price.
- PaymentID: payment identifier.
- PayerID: PayPal payer identifier
- PayerStatus: account status of the payer. Verified or not.
- Link: payment processing link. Where users can approve the payment.

# executePayment

Executes the given PayPal payment. A `payer_id` is also needed (payment needs to be approved).

```lua
local success = paypal:executePayment("PAY-XXXXX-XXXX")
-- success = true
```

Executing a payment will return a boolean value indicating success or not.