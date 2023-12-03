package domain

type User struct {
	FIO          string
	PhoneOrEmail string
	Password     string
}

type SignInRepo struct {
	PhoneOrEmail string
	Password     string
}

type UserService interface {
	SignUpUserByPhone(*SignUp) error
	SignUpUserByEmail(*SignUp) error

	SignInUserByPhone(*SignIn) bool
	SignInUserByEmail(*SignIn) bool

	GetUser(phone string) bool

	UpdatePasswordByPhone(*PasswordReset) error
	UpdatePasswordByEmail(*PasswordReset) error


}

type UserRepository interface {
	SignUpUser(*User) error
	SignInUser(*SignInRepo) (string, error)
	GetUser(phone string) bool
	UpdateUser(phoneORemail, password string) error
}

// ------------------------------------------------------------------------
// For end-point
type SignUp struct {
	FIO             string `json:"fio" binding:"required"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Otp             string `json:"otp"`
}

type UserReq struct {
	FIO          string
	PhoneOrEmail string
	Password     string
}

type UserRes struct {
	PhoneOrEmail string `json:"phone_or_email"`
}

type PasswordReset struct {
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Otp             string `json:"otp"`
}

type SignIn struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
}
