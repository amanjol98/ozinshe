package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

var ErrUserNotFound = errors.New("Вы ввели неправильный email или пароль")

func (r *UserRepository) Register(
	ctx context.Context,
	name, email, password, phone_number string,
	born_at time.Time,
) error {
	sqlQuery := `
	INSERT INTO users (name, email, password, phone_number, born_at)
	VALUES($1,$2,$3,$4,$5);
	`

	_, err := r.db.Exec(
		ctx,
		sqlQuery,
		name,
		email,
		password,
		phone_number,
		born_at,
	)
	return err
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	sqlQuery := `
	SELECT
	id,
	name,
	email,
	password,
	phone_number,
	born_at
	FROM users
	WHERE email=$1;
	`

	var user models.User

	err := r.db.QueryRow(ctx, sqlQuery, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.BornAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
