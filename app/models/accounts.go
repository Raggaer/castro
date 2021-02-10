package models

import (
	"database/sql"

	"github.com/raggaer/castro/app/database"
)

// Account struct used for tfs accounts
type Account struct {
	ID              int64
	Name            string
	Premium_ends_at int64
	Password        string
	Email           string
	Creation        int64
	Secret          sql.NullString
}

// CastroAccount struct used for castro custom accounts
type CastroAccount struct {
	ID        int64
	AccountID int64
	Points    int
	Admin     bool
}

// GetAccountByName gets an account and its castro account by the account name
func GetAccountByName(name string) (Account, CastroAccount, error) {
	// Placeholders for query values
	account := Account{}
	castroAccount := CastroAccount{}

	// Get account from database
	if err := database.DB.Get(&account, "SELECT id, name, password, premium_ends_at, email, creation, secret FROM accounts WHERE name = ?", name); err != nil {
		return account, castroAccount, err
	}

	// Get castro account from database
	if err := database.DB.Get(&castroAccount, "SELECT id, points, admin FROM castro_accounts WHERE account_id = ?", account.ID); err != nil {
		return account, castroAccount, err
	}

	return account, castroAccount, nil
}
