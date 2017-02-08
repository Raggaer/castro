package models

import "github.com/raggaer/castro/app/database"

// Account struct used for tfs accounts
type Account struct {
	ID       int64
	Name     string
	Premdays int
	Password string
	Email    string
	Lastday  int64
	Creation int64
}

// CastroAccount struct used for castro custom accounts
type CastroAccount struct {
	ID        int64
	AccountID int64
	Name      string
	Points    int
	Admin     bool
}

// GetAccountByName gets an account and its castro account by the account name
func GetAccountByName(name string) (Account, CastroAccount, error) {
	// Placeholders for query values
	account := Account{}
	castroAccount := CastroAccount{}

	// Get account from database
	if err := database.DB.Get(&account, "SELECT id, name, password, premdays, email, lastday, creation FROM accounts WHERE name = ?", name); err != nil {
		return account, castroAccount, err
	}

	// Get castro account from database
	if err := database.DB.Get(&castroAccount, "SELECT id, name, points, admin FROM castro_accounts WHERE account_id = ?", account.ID); err != nil {
		return account, castroAccount, err
	}

	return account, castroAccount, nil
}
