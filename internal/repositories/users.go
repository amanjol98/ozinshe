package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

var ErrUserNotFound = errors.New("пользователь не найден")

var ErrEmailAlreadyExists = errors.New("email уже существует")

func (r *UserRepository) Register(
	ctx context.Context,
	email, password string,
) error {
	sqlQuery := `
	INSERT INTO users (email, password)
	VALUES($1,$2);
	`

	_, err := r.db.Exec(
		ctx,
		sqlQuery,
		email,
		password,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}

		return err
	}

	return nil
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
	born_at,
	role
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
		&user.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (models.User, error) {
	sqlQuery := `
		SELECT
			id,
			name,
			email,
			phone_number,
			born_at,
			role
			FROM users
			WHERE id=$1;
	`

	var user models.User

	err := r.db.QueryRow(ctx, sqlQuery, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PhoneNumber,
		&user.BornAt,
		&user.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	userID int,
	name string,
	phoneNumber string,
	bornAt *time.Time,
) error {
	sqlQuery := `
		UPDATE users
		SET name=$1,
			phone_number=$2,
			born_at=$3
		WHERE id=$4;
		`

	_, err := r.db.Exec(ctx, sqlQuery, name, phoneNumber, bornAt, userID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, password string) error {
	sqlQuery := `
	UPDATE users
	SET password=$1
	WHERE id=$2;
	`

	_, err := r.db.Exec(ctx, sqlQuery, password, userID)
	return err
}

func (r *UserRepository) GetByIDWithPassword(ctx context.Context, userID int) (models.User, error) {
	sqlQuery := `
	SELECT
		id,
		name,
		email,
		password,
		phone_number,
		born_at,
		role
		FROM users
		WHERE id=$1;
	`

	var user models.User

	err := r.db.QueryRow(ctx, sqlQuery, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.BornAt,
		&user.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
