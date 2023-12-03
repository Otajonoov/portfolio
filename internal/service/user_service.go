package service

import (
	"portfolio/internal/domain"
	"portfolio/internal/pkg/utils/password"
	"portfolio/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userService struct {
	repository domain.UserRepository
	factory    *domain.UserFactory
}

func NewUserService(db *pgxpool.Pool) domain.UserService {
	return &userService{
		repository: repository.NewUserPostgresRepo(db),
	}
}

func (u userService) SignUpUserByPhone(user *domain.SignUp) error {

	resp := u.factory.FactoryUserByPhone(user)
	err := u.repository.SignUpUser(resp)
	if err != nil {
		return err
	}

	return nil
}

func (u userService) SignUpUserByEmail(user *domain.SignUp) error {

	resp := u.factory.FactoryUserByEmail(user)
	err := u.repository.SignUpUser(resp)
	if err != nil {
		return err
	}

	return nil
}

func (u userService) SignInUserByPhone(user *domain.SignIn) bool {

	resp := u.factory.SignInPhoneFactory(user)
	hashedPassword, err := u.repository.SignInUser(resp)
	if err != nil {
		return false
	}

	ok := password.CheckPassword(resp.Password, hashedPassword)
	return ok == nil
}

func (u userService) SignInUserByEmail(user *domain.SignIn) bool {

	resp := u.factory.SignInEmailFactory(user)
	hashedPassword, err := u.repository.SignInUser(resp)
	if err != nil {
		return false
	}

	ok := password.CheckPassword(resp.Password, hashedPassword)
	return ok == nil
}

func (u userService) GetUser(phoneORemail string) bool {

	ok := u.repository.GetUser(phoneORemail)
	return ok
}

func (u userService) UpdatePasswordByPhone(user * domain.PasswordReset) error {

	err := u.repository.UpdateUser(user.PhoneNumber, user.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u userService) UpdatePasswordByEmail(user * domain.PasswordReset) error {

	err := u.repository.UpdateUser(user.Email, user.NewPassword)
	if err != nil {
		return err
	}

	return nil
}
