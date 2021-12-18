package models

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
)

type UserTokenTestSuite struct {
	suite.Suite
}

func (suite *UserTokenTestSuite) SetupSuite() {

}

func (suite *UserTokenTestSuite) TearDownSuite() {

}

func (suite *UserTokenTestSuite) TestExpiredToken() {

	token, err := GenerateToken(100, 1*time.Second)
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(token)

	time.Sleep(3 * time.Second)

	UserClaims, err := ParseToken(token)
	suite.Assert().NotNil(err)
	suite.Assert().NotNil(UserClaims)

	validateErr := err.(*jwt.ValidationError)
	suite.Assert().Equal(validateErr.Errors, jwt.ValidationErrorExpired)

}

func (suite *UserTokenTestSuite) TestNormalToken() {

	token, err := GenerateToken(100, 1*time.Minute)
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(token)

	time.Sleep(1 * time.Second)

	userClaims, err := ParseToken(token)
	suite.Assert().Nil(err)

	suite.Assert().Equal(userClaims.UserId, int64(100))
}

func TestUserTokenSuite(t *testing.T) {
	suite.Run(t, new(UserTokenTestSuite))
}
