package handler

import (
	"errors"
	"ginfo/database"
	"ginfo/model"
	"ginfo/util"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	err := db.Where(&model.User{Email: e}).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	var user model.User
	db := database.DB
	err := db.Where(&model.User{Username: u}).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// func AuthHandler(ctx *gin.Context) {
// 	var user model.User
// 	err := ctx.ShouldBind(&user)
// 	if err != nil {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"msg": "invalid token",
// 		})
// 		return
// 	}
// 	uname := user.Username
// 	db := database.DB
// 	db.Where("username = ?", user.Username).First(&user)
// 	if user.Username == uname {
// 		tokenString, _ := util.GenToken(user.ID, user.Username)
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"msg":  "success",
// 			"data": gin.H{"token": tokenString},
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"code": 2002,
// 		"msg":  "鉴权失败",
// 	})
// }

func Login(ctx *gin.Context) {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Dept     string `json:"dept"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	input := new(LoginInput)
	var ud UserData

	if err := ctx.ShouldBind(input); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"msg":    "Error on login request",
			"data":   nil,
		})
		return
	}

	identity := input.Identity
	pass := input.Password

	user, email, err := new(model.User), new(model.User), *new(error)

	if valid(identity) {
		email, err = getUserByEmail(identity)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"msg":    "Error on email",
				"data":   err,
			})
			return
		}

	} else {
		user, err = getUserByUsername(identity)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"msg":    "Error on username",
				"data":   err,
			})
			return
		}
		if email == nil && user == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "User not found",
				"data":    err,
			})
		}

	}

	if email != nil {
		ud = UserData{
			ID:       email.ID,
			Username: email.Username,
			Email:    email.Email,
			Dept:     email.Dept,
			Phone:    email.Phone,
			Password: email.Password,
		}

	}
	if user != nil {
		ud = UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Dept:     user.Dept,
			Phone:    user.Phone,
			Password: user.Password,
		}
	}
	if !CheckPasswordHash(pass, ud.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid password",
			"data":    nil,
		})
		return
	}

	shortTokenString, longTokenString := util.GenDoubleToken(user.ID, user.Username)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "gennerate token",
		"data": gin.H{
			"shortToken": shortTokenString,
			"longToken":  longTokenString,
		},
	})

	// tokenString, _ := util.GenToken(user.ID, user.Username)
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"status": "success",
	// 	"msg":    "gennerate token",
	// 	"data":   gin.H{"token": tokenString},
	// })
}
