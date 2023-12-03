package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"portfolio/internal/domain"
	random "portfolio/internal/pkg/utils/password"
	"time"

	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserController struct {
	Service domain.UserService
}

// @Summary 	Sign-Up user by email
// @Description This api can Sign-Up new user by phone
// @Tags 		Sign-Up
// @Accept 		json
// @Produce 	json
// @Param body  body domain.SignUp true "Sign"
// @Failure 400 string Error response
// @Router /v1/sign-up-phone [post]
func (u *UserController) SignUpUserByPhone(ctx *gin.Context) {
	var req domain.SignUp
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	ok := u.Service.GetUser(req.PhoneNumber)
	if ok {
		ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud")
		return
	}

	err := u.Service.SignUpUserByPhone(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Sign-Up user by email
// @Description This api can Sign-Up new user by email
// @Tags 		Sign-Up
// @Accept 		json
// @Produce 	json
// @Param body  body domain.SignUp true "Sign"
// @Failure 400 string Error response
// @Router /v1/sign-up-email [post]
func (u *UserController) SignUpUserByEmail(ctx *gin.Context) {
	var req domain.SignUp
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	ok := u.Service.GetUser(req.Email)
	if ok {
		ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud")
		return
	}

	err := u.Service.SignUpUserByEmail(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Sign-in user by phone
// @Description This api can sign-in user by phone
// @Tags 		Sign-In
// @Accept 		json
// @Produce 	json
// @Param body  body domain.SignIn true "Sign"
// @Failure 400 string Error response
// @Router /v1/sign-in-phone [post]
func (u *UserController) SignInPhone(ctx *gin.Context) {
	var req domain.SignIn
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	ok := u.Service.SignInUserByPhone(&req)
	if !ok {
		ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud emas")
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Sign-in user by email
// @Description This api can sign-in user by email
// @Tags 		Sign-In
// @Accept 		json
// @Produce 	json
// @Param body  body domain.SignIn true "Sign"
// @Failure 400 string Error response
// @Router /v1/sign-in-email [post]
func (u *UserController) SignInEmail(ctx *gin.Context) {
	var req domain.SignIn
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	ok := u.Service.SignInUserByEmail(&req)
	if !ok {
		ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud emas")
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Check user and send code telegram
// @Description This api can Check user and send code telegram
// @Tags 		Password-reset
// @Accept 		json
// @Produce 	json
// @Param body  body domain.PasswordReset true "PasswordReset"
// @Failure 400 string Error response
// @Router /v1/check-user [post]
func (u *UserController) CheckUser(ctx *gin.Context) {
	var req domain.PasswordReset
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	ok := u.Service.GetUser(req.PhoneNumber)
	if !ok {
		ok := u.Service.GetUser(req.Email)
		if !ok {
			ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud emas")
			return
		}
	}

	//---------------------------Vohrtincha

	token := "6939554264:AAFnIu-3UHPHBQh_FJVTD-6LVgje83D32uw"

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create bot"})
		return
	}

	chatId := int64(6843604455)
	code := random.RandomPassword()

	msg := tgbotapi.NewMessage(chatId, code)

	_, err = bot.Send(msg)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to send message"})
		return
	}

	redis := newRedisClient()

	err = redis.Set(context.Background(), req.Email, code, time.Minute*5).Err()
	if err != nil {
		err = redis.Set(context.Background(), req.PhoneNumber, code, time.Minute*5).Err()
		if err != nil {
			log.Println("Error setting value in Redis:", err)
			return
		}
	}

	//--------------------------

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Reset password
// @Description This api can reset password
// @Tags 		Password-reset
// @Accept 		json
// @Produce 	json
// @Param body  body domain.PasswordReset true "PasswordReset"
// @Failure 400 string Error response
// @Router /v1/check-code [post]
func (u *UserController) CheckCodePhone(ctx *gin.Context) {
	var req domain.PasswordReset
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	redis := newRedisClient()

	value, err := redis.Get(context.Background(), req.PhoneNumber).Result()
	if err != nil {
		value, _ := redis.Get(context.Background(), req.Email).Result()
		if req.Otp == value {
			ctx.JSON(http.StatusOK, "ok")
			return
		}
	}

	if req.Otp == value {
		ctx.JSON(http.StatusOK, "ok")
		return
	}

	ctx.JSON(http.StatusBadRequest, "qode noto'g'ri kiritildi")
}

// @Summary 	Update password by phone
// @Description This api can update password by phone
// @Tags 		Password-reset
// @Accept 		json
// @Produce 	json
// @Param body  body domain.PasswordReset true "PasswordReset"
// @Failure 400 string Error response
// @Router /v1/update-password-phone [post]
func (u *UserController) UpdatePasswordByPhone(ctx *gin.Context) {
	var req domain.PasswordReset
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	if req.ConfirmPassword != req.NewPassword {
		ctx.JSON(http.StatusBadRequest, "Qodlar bir biriga mos emas")
		return
	}

	err := u.Service.UpdatePasswordByPhone(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud emas")
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Update password by phone
// @Description This api can update password by phone
// @Tags 		Password-reset
// @Accept 		json
// @Produce 	json
// @Param body  body domain.PasswordReset true "PasswordReset"
// @Failure 400 string Error response
// @Router /v1/update-password-email [post]
func (u *UserController) UpdatePasswordByEmail(ctx *gin.Context) {
	var req domain.PasswordReset
	json.NewDecoder(ctx.Request.Body).Decode(&req)

	if req.ConfirmPassword != req.NewPassword {
		ctx.JSON(http.StatusBadRequest, "Qodlar bir biriga mos emas")
		return
	}

	err := u.Service.UpdatePasswordByEmail(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bunday foydalanuvchi mavjud emas")
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

// @Summary 	Send code telegram
// @Description This api can send code telegram
// @Tags 		Send
// @Accept 		json
// @Produce 	json
// @Param 		phone_email query string true "phone_email"
// @Failure 400 string Error response
// @Router /v1/send-code [post]
func (u *UserController) SendCode(ctx *gin.Context) {
	phonrORemail := ctx.Query("phone_email")
	token := "6939554264:AAFnIu-3UHPHBQh_FJVTD-6LVgje83D32uw"

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create bot"})
		return
	}

	chatId := int64(6843604455)
	code := random.RandomPassword()

	msg := tgbotapi.NewMessage(chatId, code)

	_, err = bot.Send(msg)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to send message"})
		return
	}

	redis := newRedisClient()

	err = redis.Set(context.Background(), phonrORemail, code, time.Minute*1).Err()
	if err != nil {
		fmt.Println("Error setting value in Redis:", err)
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return client
}
