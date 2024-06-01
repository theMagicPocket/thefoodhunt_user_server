package user

import (
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

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user entity.User
	user.ID = entity.GenerateID()
	istLocation,_:= time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	user.CreatedAt = formattedTime
	if err :=  ctx.ShouldBindJSON(&user); err != nil {
	ctx.JSON(http.StatusBadRequest,gin.H{"message": err.Error()})
	return
	};
 err := uc.UserService.CreateUser(&user)
 if err != nil{
  ctx.JSON(http.StatusBadGateway,gin.H{"message": err.Error()})
  return
 }
  ctx.JSON(http.StatusOK, gin.H{"message":"success","status": http.StatusOK,"userid":user.ID})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := uc.UserService.GetUser((&userId))
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
	user.ID = ctx.Param("id")
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	istLocation,_:= time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	user.CreatedAt = formattedTime
	ctx.JSON(http.StatusOK, gin.H{"message": "success","status":http.StatusOK,"updated":user})

}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	log.Println(userId)
	err := uc.UserService.DeleteUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success","status":http.StatusOK})
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	userroute.POST("/create", uc.CreateUser)
	userroute.GET("/get/:id", uc.GetUser)
	userroute.GET("/getusers", uc.GetAll)
	userroute.PATCH("/update/:id", uc.UpdateUser)
	userroute.DELETE("/delete/:id", uc.DeleteUser)
}
