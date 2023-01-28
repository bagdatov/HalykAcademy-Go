package usecase

import (
	"context"
	"time"
	"web/domain"
)

type myUseCase struct {
	db domain.Repository
	domain.ActiveWallets
}

// NewUseCase returns new struct that represents domain.UseCase interface.
// It also includes the initilized map with mutex for safe concurrent mining process.
func NewUseCase(db domain.Repository) domain.UseCase {
	return &myUseCase{
		db: db,
		ActiveWallets: domain.ActiveWallets{
			Wallets: make(map[string]*domain.CryptoWallet),
		},
	}
}

// FindUser is designed so that in case of non-result underlying database shall return user as nil,
// so this prevents from sql.ErrorNoRows error.
func (uc *myUseCase) FindUser(ctx context.Context, username, password string) (*domain.User, error) {

	u, err := uc.db.FindUser(ctx, username, password)

	if u == nil {
		return nil, domain.ErrNotFound

	} else if err != nil {
		return nil, err
	}

	return u, nil
}

// FindWallet is designed so that in case of non-result underlying database shall return user as nil,
// so this prevents from sql.ErrorNoRows error.
func (uc *myUseCase) FindWallet(ctx context.Context, requester *domain.User, walletName string) (*domain.CryptoWallet, error) {

	w, err := uc.db.FindWallet(ctx, walletName)
	if w == nil {
		return nil, domain.ErrNotFound

	} else if err != nil {
		return nil, err
	}

	if w.OwnerName != requester.Username {
		return nil, domain.ErrNotAuthorized
	}

	return w, nil
}

func (uc *myUseCase) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := uc.db.GetUserByID(ctx, id)

	if user == nil {
		return nil, domain.ErrNotFound

	} else if err != nil {
		return nil, err
	}

	walletNames, err := uc.db.GetWalletsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Wallets = walletNames

	return user, nil
}

func (uc *myUseCase) CreateUser(ctx context.Context, u *domain.User) error {
	return uc.db.CreateUser(ctx, u)
}

func (uc *myUseCase) CreateWallet(ctx context.Context, cw *domain.CryptoWallet) error {
	return uc.db.CreateWallet(ctx, cw)
}

func (uc *myUseCase) DeleteUser(ctx context.Context, requester *domain.User, targetID int) error {

	u, err := uc.db.GetUserByID(ctx, targetID)
	if u == nil {
		return domain.ErrNotFound

	} else if err != nil {
		return err
	}

	if requester.ID != targetID {
		return domain.ErrNotAuthorized
	}

	return uc.db.DeleteUser(ctx, u)
}

func (uc *myUseCase) DeleteWallet(ctx context.Context, requester *domain.User, targetWallet string) error {

	w, err := uc.FindWallet(ctx, requester, targetWallet)
	if err != nil {
		return err
	}

	return uc.db.DeleteWallet(ctx, w)
}

func (uc *myUseCase) StartMine(ctx context.Context, requester *domain.User, walletName string) error {
	w, err := uc.FindWallet(ctx, requester, walletName)
	if err != nil {
		return err
	}
	w.Stop = make(chan struct{})

	uc.Lock()
	defer uc.Unlock()

	_, ok := uc.Wallets[walletName]
	if ok {
		return domain.ErrInitilized
	}

	uc.Wallets[walletName] = w
	go func() {
		for {
			ctxTimeOut, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			select {
			case <-w.Stop:
				cancel()
				return
			case <-ctxTimeOut.Done():
				uc.db.UpdateWalletSum(context.Background(), w.Name, 1)
			}
		}
	}()
	return nil
}

func (uc *myUseCase) StopMine(ctx context.Context, requester *domain.User, walletName string) error {
	_, err := uc.FindWallet(ctx, requester, walletName)
	if err != nil {
		return err
	}

	uc.Lock()
	w, ok := uc.Wallets[walletName]
	if !ok {
		uc.Unlock()
		return domain.ErrNotInitilized
	}
	delete(uc.Wallets, walletName)
	uc.Unlock()
	w.Stop <- struct{}{}

	return nil
}
