package apis

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golbeng-original/chomakers-web/models"
)

var accessTokenExpire time.Duration
var userRepository *models.UserRespository

func CheckAuthentication(c *gin.Context) (bool, error) {

	cookies := c.Request.Cookies()
	if len(cookies) == 0 {
		return false, nil
	}

	cookie, err := c.Request.Cookie("access-token")
	if err != nil {
		return false, err
	}

	result, err := models.AuthorizedFromToken(cookie.Value, userRepository, accessTokenExpire)
	if err != nil {
		log.Printf("models.AuthorizedFromToken err [%s]", err)
		return false, err
	}

	if result.ResultType == models.AuthorizedResultFailed {
		log.Printf("models.AuthorizedResultFailed")
		return false, nil
	}

	if result.ResultType == models.AuthorizedResultRenewalAccessToken {
		c.SetCookie("access-token", *result.RenewalAccessToken, 2147483647, "/", "", false, true)
	}

	return true, nil
}

func LoginApis(api *gin.RouterGroup, repositoryConfigure *models.RepositoryConfigure) {

	userRepository = repositoryConfigure.UserRepository
	accessTokenExpire = repositoryConfigure.AccessTokenExpireTime

	api.POST("/login", func(c *gin.Context) {

		var reqLogin RequestLogin
		c.ShouldBindJSON(&reqLogin)

		userModel, err := userRepository.GetUserModelFromUserName(reqLogin.UserName)

		if err != nil && !errors.Is(err, &models.UserNotExistError{}) {
			errorMessage := fmt.Sprintf("login error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		responseLogin := ResponseLogin{}
		if errors.Is(err, &models.UserNotExistError{}) {
			responseLogin.LoginResult = 1
		} else if userModel.Password != reqLogin.Password {
			responseLogin.LoginResult = 2
		} else {
			responseLogin.LoginResult = 0
		}

		var resultAccessToken string
		if responseLogin.LoginResult == 0 {

			refreshToken, refreshTokenErr := models.GenerateToken(userModel.Id, repositoryConfigure.RefreshTokenExpireTime)
			accessToken, accessTokenErr := models.GenerateToken(userModel.Id, repositoryConfigure.AccessTokenExpireTime)
			if refreshTokenErr != nil || accessTokenErr != nil {
				responseLogin.LoginResult = 3
			} else {
				// refreshToken 등록
				userRepository.SetRefreshToken(userModel.Id, refreshToken)

				resultAccessToken = accessToken
			}
		}

		if responseLogin.LoginResult == 0 {
			c.SetCookie("access-token", resultAccessToken, 2147483647, "/", "", false, true)
		}

		responsePresent, err := SuccessResponsePresent(c, responseLogin)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusNotFound, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.GET("/authentication", func(c *gin.Context) {

		isAuthentication, err := CheckAuthentication(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, FailedResponsePreset(err.Error()))
			c.Abort()
			return
		}

		if !isAuthentication {
			c.JSON(http.StatusUnauthorized, FailedResponsePreset(""))
			c.Abort()
			return
		}

		responsePresent, err := SuccessResponsePresent(c, nil)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})

	api.GET("/logout", func(c *gin.Context) {

		// token 확인
		accessToken, err := c.Cookie("access-token")
		if err != nil {
			responsePresent, _ := SuccessResponsePresent(c, nil)
			c.JSON(http.StatusOK, responsePresent)
			return
		}

		// accessToken 제거 하기
		c.SetCookie("access-token", "", 2147483647, "/", "", false, true)

		// refreshToken 제거 하기
		userClaims, err := models.ParseToken(accessToken)
		if err == nil {
			userRepository.ClearRefreshToken(userClaims.UserId)
		}

		responsePresent, err := SuccessResponsePresent(c, nil)
		if err != nil {
			errorMessage := fmt.Sprintf("create SuccessResponsePresent error [%v]", err)
			c.JSON(http.StatusInternalServerError, FailedResponsePreset(errorMessage))
			return
		}

		c.JSON(http.StatusOK, responsePresent)
	})
}
