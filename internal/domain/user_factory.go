package domain

import "portfolio/internal/pkg/utils/password"

type UserFactory struct{}

func (u *UserFactory) FactoryUserByPhone(source *SignUp) *User {

	hashPas, _ := password.HashPassword(source.NewPassword)
	return &User{
		FIO:          source.FIO,
		PhoneOrEmail: source.PhoneNumber,
		Password:     hashPas,
	}
}

func (u *UserFactory) FactoryUserByEmail(source *SignUp) *User {

	hashPas, _ := password.HashPassword(source.NewPassword)
	return &User{
		FIO:          source.FIO,
		PhoneOrEmail: source.Email,
		Password:     hashPas,
	}
}

func (u *UserFactory) SignInPhoneFactory(source *SignIn) *SignInRepo {

	return &SignInRepo{
		PhoneOrEmail: source.PhoneNumber,
		Password: source.Password,
	}
}

func (u *UserFactory) SignInEmailFactory(source *SignIn) *SignInRepo {

	return &SignInRepo{
		PhoneOrEmail: source.Email,
		Password: source.Password,
	}
}