package repository

import (
	"context"
	"log"
	"portfolio/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userPostgresRepo struct {
	db *pgxpool.Pool
}

func NewUserPostgresRepo(db *pgxpool.Pool) domain.UserRepository {
	return &userPostgresRepo{
		db: db,
	}
}

func (u userPostgresRepo) SignUpUser(user *domain.User) error {
	query := `
		INSERT INTO users(
			fio,
			phone_or_email,
			password
		) VALUES ($1, $2, $3)
	`
	_, err := u.db.Exec(context.Background(), query, user.FIO, user.PhoneOrEmail, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (u userPostgresRepo) SignInUser(user *domain.SignInRepo) (string, error) {
	var hashedPassword string
	query := `
		SELECT password
		FROM users
		WHERE phone_or_email = $1;
	`
	err := u.db.QueryRow(context.Background(), query, user.PhoneOrEmail).Scan(&hashedPassword)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return hashedPassword, nil
}

func (u userPostgresRepo) GetUser(phoneOrEmail string) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE phone_or_email = $1
		)

	`
	err := u.db.QueryRow(context.Background(), query, phoneOrEmail).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (u userPostgresRepo) UpdateUser(phoneORemail, password string) error {

	query := `
		UPDATE users SET
			password = $1
		WHERE 
			phone_or_email = $2
	`
	_, err := u.db.Exec(context.Background(), query, password, phoneORemail)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
