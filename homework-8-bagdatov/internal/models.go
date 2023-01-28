package internal

import (
	"context"
	"errors"
	"sync"
	"time"
)

type User struct {
	ID       int
	Username string
	Password string
}

type CryptoWallet struct {
	Name    string
	OwnerID int
	Amount  int64
	Stop    chan struct{}
	sync.RWMutex
}

func (c *CryptoWallet) Mine() {

	ctxTime, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-ctxTime.Done():
		c.Lock()
		c.Amount++
		c.Unlock()
	case <-c.Stop:
		return
	}
}

type UserBase struct {
	store   map[int]*User
	wallets map[string]*CryptoWallet
	mu      sync.Mutex
}

var errNotAuthorized = errors.New("invalid user credentials")
var errExists = errors.New("object with these credentials already exists")
var errNotFound = errors.New("object with these credentials not exists")

func NewUserBase() *UserBase {
	return &UserBase{
		store:   map[int]*User{0: {0, "admin", "admin"}},
		wallets: make(map[string]*CryptoWallet),
	}
}

func (ub *UserBase) FindWalletsByID(id int) (walletNames []string) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	for _, w := range ub.wallets {
		if w.OwnerID == id {
			walletNames = append(walletNames, w.Name)
		}
	}
	return
}

func (ub *UserBase) FindUser(username, password string) (*User, error) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	for _, user := range ub.store {
		if user.Username == username && user.Password == password {
			return user, nil
		}
	}
	return nil, errNotAuthorized
}

func (ub *UserBase) FindUsernameByID(id int) (string, error) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	u, ok := ub.store[id]
	if !ok {
		return "", errNotFound
	}
	return u.Username, nil
}

func (ub *UserBase) CreateUser(id int, username, password string) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	_, ok := ub.store[id]
	if ok {
		return errExists
	}

	for _, u := range ub.store {
		if u.Username == username {
			return errExists
		}
	}

	ub.store[id] = &User{
		ID:       id,
		Username: username,
		Password: password,
	}

	return nil
}

func (ub *UserBase) CreateWallet(name string, ownerID int) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	_, ok := ub.wallets[name]
	if ok {
		return errExists
	}

	ub.wallets[name] = &CryptoWallet{
		Name:    name,
		OwnerID: ownerID,
		Stop:    make(chan struct{}),
	}

	return nil
}

func (ub *UserBase) WalletSum(walletName string, requesterID int) (int64, error) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	w, ok := ub.wallets[walletName]
	if !ok {
		return 0, errNotFound
	}

	var num int64
	w.Lock()
	if w.OwnerID != requesterID {
		return 0, errNotAuthorized
	}
	num = w.Amount
	w.Unlock()

	return num, nil
}

func (ub *UserBase) StartMine(walletName string, requesterID int) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	w, ok := ub.wallets[walletName]
	if !ok {
		return errNotFound
	}

	w.Lock()
	if w.OwnerID != requesterID {
		return errNotAuthorized
	}
	w.Unlock()

	go func() {
		for {
			select {
			case <-w.Stop:
				return
			default:
				w.Mine()
			}
		}
	}()

	return nil
}

func (ub *UserBase) StopMine(walletName string, requesterID int) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	w, ok := ub.wallets[walletName]
	if !ok {
		return errNotFound
	}

	w.Lock()
	if w.OwnerID != requesterID {
		return errNotAuthorized
	}
	w.Unlock()

	w.Stop <- struct{}{}

	return nil
}
