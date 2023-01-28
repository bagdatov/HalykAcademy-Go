package domain

import "context"

// UseCase is bussiness logic over our web service.
type UseCase interface {
	// FindUser is used to for authorization
	// and called along at middleware, before any request proceeded.
	FindUser(ctx context.Context, username, password string) (*User, error)

	// Returns CryptoWallet details.
	// It is validating the requester, so no one can access wallet except owner.
	FindWallet(ctx context.Context, requester *User, walletName string) (*CryptoWallet, error)

	// Any authorized can access details about user
	// including his ID, Username and list of wallets, but not the password and wallet sums.
	GetUserByID(ctx context.Context, id int) (*User, error)

	// Register new user, returns error if username/id already exist.
	CreateUser(ctx context.Context, u *User) error

	// Register new wallet for user, returns error if walletname is busy.
	CreateWallet(ctx context.Context, cw *CryptoWallet) error

	// At this moment, we have ability to delete user/wallet,
	// but this functionality is not added to web service yet.
	// It is designed so only owner can delete them.
	DeleteUser(ctx context.Context, requester *User, targetID int) error
	DeleteWallet(ctx context.Context, requester *User, targetWallet string) error

	// Starts goroutine for mining, and registers wallet as active wallet.
	// Also validating if it is an owner.
	// Duplicated request for start/stop returns error.
	// Stop is sending signal (through wallet channel) to goroutine to end mining.
	// And removes wallet from active wallets.
	StartMine(ctx context.Context, requester *User, walletName string) error
	StopMine(ctx context.Context, requester *User, walletName string) error
}

// Repository is representing database.
// At this moment it is PostgreSQL, but can be updated anytime,
// so bussiness logic will not be affected.
type Repository interface {
	FindUser(ctx context.Context, username, password string) (*User, error)
	FindWallet(ctx context.Context, walletName string) (*CryptoWallet, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetWalletsByUserID(ctx context.Context, id int) ([]string, error)
	CreateUser(ctx context.Context, u *User) error
	CreateWallet(ctx context.Context, cw *CryptoWallet) error
	DeleteUser(ctx context.Context, u *User) error
	DeleteWallet(ctx context.Context, cw *CryptoWallet) error
	UpdateWalletSum(ctx context.Context, walletName string, toAdd int) error
	CloseConnection()
}
