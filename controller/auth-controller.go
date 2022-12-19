package controller

import (
	"net/http"
	"strconv"

	"github.com/KuraoHikari/library-app/dto"
	"github.com/KuraoHikari/library-app/entity"
	"github.com/KuraoHikari/library-app/helper"
	"github.com/KuraoHikari/library-app/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController{
	return &authController{
		authService: authService,
		jwtService: jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context){
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email,loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID,10))
		v.Token = generatedToken
		res := helper.BuildResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, res)
		return
	}
	res :=helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
}
func (c *authController) Register(ctx *gin.Context){
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if !c.authService.IsDuplicateEmail(registerDTO.Email){
		res := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, res)
	}else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		res := helper.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, res)
	}
}