package models

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtTokenKey = []byte("qudghweb_jwt_secret_key")

const (
	AuthorizedResultSuccess uint32 = iota
	AuthorizedResultRenewalAccessToken
	AuthorizedResultFailed
)

type UserClaims struct {
	UserId int64
	jwt.StandardClaims
}

type AuthorizedResult struct {
	ResultType         uint32
	RenewalAccessToken *string
}

func GenerateToken(userId int64, expireTime time.Duration) (string, error) {

	expireTimeStamp := time.Now().Add(expireTime)

	claims := &UserClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTimeStamp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtTokenKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(token string) (*UserClaims, error) {

	resultClaims := &UserClaims{}

	_, err := jwt.ParseWithClaims(token, resultClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtTokenKey, nil
	})

	return resultClaims, err
}

func AuthorizedFromToken(token string, userRepository *UserRespository, accessTokenExpireTime time.Duration) (*AuthorizedResult, error) {

	userClaims, err := ParseToken(token)
	if userClaims == nil && err != nil {
		return nil, err
	}

	authorizedResult := &AuthorizedResult{
		ResultType: AuthorizedResultSuccess,
	}

	if err != nil {

		// access token 만료 시간 도달
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {

			// refresh token 확인
			userModel, err := userRepository.GetUserModel(userClaims.UserId)
			if err != nil {
				log.Printf("userRepository.GetUserModel error[%s]\n", err)
				return nil, err
			}

			// refresh Token이 없다.
			if userModel.RefreshToken == nil || len(*userModel.RefreshToken) == 0 {
				authorizedResult.ResultType = AuthorizedResultFailed

				log.Printf("userRepository.GetUserModel access token is null\n")
				return authorizedResult, nil
			}

			// refresh Token 만료
			userClaims, err := ParseToken(*userModel.RefreshToken)
			if err != nil && err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {

				log.Printf("refresh token expired\n")

				userRepository.SetRefreshToken(userModel.Id, "")

				authorizedResult.ResultType = AuthorizedResultFailed
				return authorizedResult, nil
			}

			// refresh Token Claims 값에 이상하다.
			if userModel.Id != userClaims.UserId {
				userRepository.SetRefreshToken(userModel.Id, "")

				authorizedResult.ResultType = AuthorizedResultFailed
				return authorizedResult, nil
			}

			// 새 Access-Token 발급
			newAccessToken, err := GenerateToken(userModel.Id, accessTokenExpireTime)
			if err != nil {
				return nil, err
			}

			authorizedResult.ResultType = AuthorizedResultRenewalAccessToken
			authorizedResult.RenewalAccessToken = &newAccessToken
			return authorizedResult, err
		} else {
			// access Token 값이 정상이 아니다.

			authorizedResult.ResultType = AuthorizedResultFailed
			return authorizedResult, err
		}

	} else {

		// userCalims 정보가 정상이 아니다.
		userModel, _ := userRepository.GetUserModel(userClaims.UserId)
		if userModel == nil {
			authorizedResult.ResultType = AuthorizedResultFailed
			return authorizedResult, nil
		}
	}

	return authorizedResult, nil
}
