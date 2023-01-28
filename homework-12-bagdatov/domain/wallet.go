package domain

import "sync"

// CryptoWallet includes wallet details and channel
// to stop mining process.
type CryptoWallet struct {
	ID        int           `json:"walletID"`
	Name      string        `json:"walletName"`
	OwnerName string        `json:"ownerName"`
	OwnerID   int           `json:"ownerID"`
	Amount    int64         `json:"amount"`
	Stop      chan struct{} `json:"-"`
}

// ActiveWallets for mining process.
// Protected with mutex for concurrent use.
type ActiveWallets struct {
	Wallets map[string]*CryptoWallet
	sync.RWMutex
}
