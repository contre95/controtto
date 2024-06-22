package pnl

type AccountID uint8

type Account struct {
	ID          AccountID
	Name        string
	Description string
	Type        AssetType
}

// Accounts is the repository that handles the CRUD of Accounts
type Accounts interface {
	AddAccount(a Account) error
	ListAccounts() ([]Account, error)
	GetAccount(id AccountID) (*Account, error)
}
