package pg

import (
	"context"
	"errors"
	"fmt"
	"time"
	"web/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type sqlRepository struct {
	*pgxpool.Pool
}

const (
	username = "postgres"
	password = "postgres"
	hostname = "postgres"
	port     = 5432
	dbName   = "postgres"
)

func NewSQLRepository() (domain.Repository, error) {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, dbName)

	config, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MaxConnLifetime = 5 * time.Minute

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err

	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &sqlRepository{db}, nil
}

func (db *sqlRepository) FindUser(ctx context.Context, username, password string) (*domain.User, error) {
	user := &domain.User{}

	err := db.QueryRow(ctx,
		`SELECT ID,Name,Password 
		FROM users 
		WHERE Name=$1 AND Password=$2`,
		username, password,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	return user, err
}

func (db *sqlRepository) FindWallet(ctx context.Context, walletName string) (*domain.CryptoWallet, error) {
	w := &domain.CryptoWallet{}

	err := db.QueryRow(ctx,
		`SELECT wallets.ID, wallets.Name, wallets.OwnerID, wallets.Amount, users.Name
		FROM wallets
		LEFT JOIN users ON users.ID=wallets.OwnerID
		WHERE wallets.Name=$1`, walletName,
	).Scan(&w.ID, &w.Name, &w.OwnerID, &w.Amount, &w.OwnerName)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return w, err
}

func (db *sqlRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}

	err := db.QueryRow(ctx,
		`SELECT ID,Name 
		FROM users 
		WHERE ID=$1`,
		id,
	).Scan(&user.ID, &user.Username)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (db *sqlRepository) GetWalletsByUserID(ctx context.Context, id int) ([]string, error) {
	wallets := []string{}

	rows, err := db.Query(ctx,
		`SELECT Name 
		FROM wallets 
		WHERE OwnerID=$1`, id,
	)

	if err == pgx.ErrNoRows {
		return nil, nil

	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var w string

		if err := rows.Scan(&w); err != nil {
			return nil, err
		}

		wallets = append(wallets, w)
	}

	return wallets, rows.Err()
}

func (db *sqlRepository) CreateUser(ctx context.Context, u *domain.User) error {
	_, err := db.Exec(ctx,
		`INSERT INTO users(ID,Name,Password) 
		VALUES($1,$2,$3)`,
		u.ID, u.Username, u.Password,
	)

	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrExists
		}
	}
	return err
}

func (db *sqlRepository) CreateWallet(ctx context.Context, cw *domain.CryptoWallet) error {
	_, err := db.Exec(ctx,
		`INSERT INTO wallets(Name,OwnerID,Amount)
		VALUES($1,$2,$3)`,
		cw.Name, cw.OwnerID, cw.Amount,
	)
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrExists
		}
	}
	return err
}

func (db *sqlRepository) DeleteUser(ctx context.Context, u *domain.User) error {
	_, err := db.Exec(ctx, "DELETE FROM users WHERE ID=$1", u.ID)
	return err
}

func (db *sqlRepository) DeleteWallet(ctx context.Context, cw *domain.CryptoWallet) error {
	_, err := db.Exec(ctx, "DELETE FROM wallets WHERE Name=$1", cw.Name)
	return err
}

func (db *sqlRepository) CloseConnection() {
	db.Close()
}

func (db *sqlRepository) UpdateWalletSum(ctx context.Context, walletName string, toAdd int) error {
	_, err := db.Exec(ctx, "UPDATE wallets SET Amount=Amount+$1 WHERE Name=$2", toAdd, walletName)
	return err
}
