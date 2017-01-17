package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"gopkg.in/gomail.v2"
)

// SetMailMetaTable sets the mail metatable of the given state
func SetMailMetaTable(luaState *lua.LState) {
	// Create and set the mail metatable
	mailMetaTable := luaState.NewTypeMetatable(MailMetaTableName)
	luaState.SetGlobal(MailMetaTableName, mailMetaTable)

	// Set all mail metatable functions
	luaState.SetFuncs(mailMetaTable, mailMethods)
}

// SendMail sends a mail to the given direction
func SendMail(L *lua.LState) int {
	// Get information table
	tbl := L.Get(2)

	// Check for valid type
	if tbl.Type() != lua.LTTable {

		L.ArgError(1, "Invalid email send type. Expected table")
		return 0
	}

	// Convert table to map
	info := TableToMap(tbl.(*lua.LTable))

	// Get to header
	to, ok := info["to"].(string)

	if !ok {
		L.ArgError(1, "Missing 'to' table field")
		return 0
	}

	// Get subject
	subject, ok := info["subject"].(string)

	if !ok {
		L.ArgError(1, "Missing 'subject' table field")
		return 0
	}

	// Get email body
	body, ok := info["body"].(string)

	if !ok {
		L.ArgError(1, "Missing 'body' table field")
		return 0
	}

	// Create new gomail object
	m := gomail.NewMessage()

	// Set from header
	m.SetHeader("From", util.Config.Mail.Username)

	// Set to header
	m.SetHeader("To", to)

	// Set subject
	m.SetHeader("Subject", subject)

	// Set body
	m.SetBody("text/html", body)

	// Create dialer
	d := gomail.NewDialer(
		util.Config.Mail.Server,
		util.Config.Mail.Port,
		util.Config.Mail.Username,
		util.Config.Mail.Password,
	)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		L.RaiseError("Cannot send email: %v", err)
		return 0
	}

	return 0
}
