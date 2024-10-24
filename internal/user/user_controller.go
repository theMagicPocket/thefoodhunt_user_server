package user

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService UserService
}

func New(userService UserService) UserController {
	return UserController{
		UserService: userService,
	}
}

type CreateUserRequest struct {
	AuthId       string `bson:"auth_id" json:"auth_id"`
	Name         string `bson:"name" json:"name"`
	Email        string `bson:"email" json:"email"`
	Phone        string `bson:"phone" json:"phone"`
	ProfilePhoto string `bson:"profile_photo" json:"profile_photo"`
	// UserAddress  []entity.UserAddress `bson:"user_address" json:"user_address"`
	CreatedAt string `bson:"created_at" json:"created_at"`
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var cur CreateUserRequest
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}
	cur.AuthId = userId.(string)
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	cur.CreatedAt = formattedTime
	if err := ctx.ShouldBindJSON(&cur); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newUser := entity.User{
		AuthId:       cur.AuthId,
		Name:         cur.Name,
		Email:        cur.Email,
		Phone:        cur.Phone,
		ProfilePhoto: cur.ProfilePhoto,
		UserAddress:  make([]entity.UserAddress, 0),
		CreatedAt:    cur.CreatedAt,
	}

	user, err := uc.UserService.CreateUser(&newUser)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "User already found", "status": http.StatusConflict})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK, "userid": user})
}

func (uc *UserController) GetUserByAuthId(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := uc.UserService.GetUserByAuthId(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserById(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := uc.UserService.GetUserById(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)

}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user entity.User
	user.Id = ctx.Param("id")
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	user.CreatedAt = formattedTime
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK, "updated": user})

}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	log.Println(userId)
	err := uc.UserService.DeleteUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK})
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	userroute.POST("/create", uc.CreateUser)
	userroute.GET("/get/:id", uc.GetUserById)
	userroute.GET("/getbyauthid/:id", uc.GetUserByAuthId)
	userroute.GET("/getusers", uc.GetAll)
	userroute.PATCH("/update/:id", uc.UpdateUser)
	userroute.DELETE("/delete/:id", uc.DeleteUser)
}
