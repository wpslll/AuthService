package repository

import (
	"AuthService/internal/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	conn *pgxpool.Pool
}

func NewDB(ctx context.Context, connstring string) (*Database, error) {
	db, err := pgxpool.New(ctx, connstring)
	if err != nil { return nil, err }
	return &Database{ conn: db }, nil
}

func (d *Database) Create(ctx context.Context, user domain.User) error {
	query1 := `INSERT INTO Users (username, password)
			VALUES ($1, $2)`
	var pgErr *pgconn.PgError
	if _, err := d.conn.Exec(ctx, query1, user.Username, user.Password); err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
    		return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (d *Database) Auth(ctx context.Context, user domain.User) error {
	query := `SELECT username
	FROM Users
	WHERE username = $1 AND password = $2`
	var foundUser domain.User
	if err := d.conn.QueryRow(ctx, query, user.Username, user.Password).Scan(&foundUser.Username); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		} else {
			return err
		}
	}
	return nil
}
