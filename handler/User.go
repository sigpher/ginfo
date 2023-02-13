package handler

import (
	"ginfo/database"
	"ginfo/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateUser(ctx *gin.Context) {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Dept     string `json:"dept"`
		Phone    string `json:"phone"`
	}
	db := database.DB
	user := new(model.User)
	if err := ctx.ShouldBind(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "Review your input",
			"data":   err,
		})
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "Review your input",
			"data":   err,
		})
		return
	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "Couldn't create user",
			"data":   err,
		})
		return
	}
	newUser := NewUser{
		Username: user.Username,
		Email:    user.Email,
		Dept:     user.Dept,
		Phone:    user.Phone,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Created user",
		"data":   newUser,
	})
}

func GetUsers(ctx *gin.Context) {
	var users []model.User
	db := database.DB
	db.Find(&users)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": users,
	})
}

func GetUser(ctx *gin.Context) {
	var user model.User
	id := ctx.Param("id")
	db := database.DB
	db.Find(&user, id)

	if user.Username == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"msg":    "user not found",
			"data":   nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "user found",
		"data":   user,
	})
}
