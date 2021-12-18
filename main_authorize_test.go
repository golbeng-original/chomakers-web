package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/golbeng-original/chomakers-web/apis"
	"github.com/golbeng-original/chomakers-web/models"
)

func getCookieValue(response *http.Response, name string) string {
	var tokenValue string
	cookies := response.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			tokenValue = cookie.Value
			break
		}
	}

	return tokenValue
}

type AuthroizeTestSuite struct {
	suite.Suite

	dbConnection *models.DBConnection
	testServer   *httptest.Server

	userRepository     *models.UserRespository
	repositoryCongiure *models.RepositoryConfigure
	testUserName       string
	testPassword       string
	testAccessToken    string
}

func (suite *AuthroizeTestSuite) getUrl() string {
	return suite.testServer.URL
}

func (suite *AuthroizeTestSuite) SetupSuite() {

	dbConnection := models.DBConnection{}
	dbConnection.Open("file::memory:?mode=memory&cache=shared")

	suite.dbConnection = &dbConnection

	repositoryConfigure := &models.RepositoryConfigure{}
	repositoryConfigure.Init(&dbConnection)
	repositoryConfigure.IsCheckAuthorize = true

	suite.repositoryCongiure = repositoryConfigure
	suite.userRepository = repositoryConfigure.UserRepository

	suite.testUserName = "root"
	suite.userRepository.AddUser(suite.testUserName, "1234")
	userModel, err := suite.userRepository.GetUserModelFromUserName(suite.testUserName)
	suite.Assert().Nil(err)

	suite.testPassword = userModel.Password

	suite.testServer = httptest.NewServer(Setup(repositoryConfigure, "./assets/images"))
}

func (suite *AuthroizeTestSuite) TearDownSuite() {
	suite.testServer.Close()
	suite.dbConnection.Close()
}

func (suite *AuthroizeTestSuite) BeforeTest(suiteName, testName string) {

	if testName == "TestExpireAccess" {
		suite.repositoryCongiure.AccessTokenExpireTime = time.Second
	} else if testName == "TestExpireBothToken" {
		suite.repositoryCongiure.AccessTokenExpireTime = time.Second
		suite.repositoryCongiure.RefreshTokenExpireTime = time.Second
	}

	if testName == "TestLoginAccess" ||
		testName == "TestExpireAccess" ||
		testName == "TestExpireBothToken" ||
		testName == "TestLoginAfterLogoutState" {

		reqLogin := apis.RequestLogin{
			UserName: suite.testUserName,
			Password: suite.testPassword,
		}

		bytes, err := json.Marshal(reqLogin)
		suite.Assert().Nil(err)

		res, err := http.Post(suite.getUrl()+"/api/login", "application/json", strings.NewReader(string(bytes)))
		suite.Assert().Nil(err)
		suite.Assert().Equal(res.StatusCode, 200)

		tokenValue := getCookieValue(res, "access-token")

		suite.Assert().NotEmpty(tokenValue)
		suite.testAccessToken = tokenValue
	}

}

func (suite *AuthroizeTestSuite) TestNotLoginAccess() {

	reqCreatePotofolio := apis.RequestCreatePotofolio{}

	bytes, err := json.Marshal(reqCreatePotofolio)
	suite.Assert().Nil(err)

	res, err := http.Post(suite.getUrl()+"/api/potofolio", "application/json", strings.NewReader(string(bytes)))
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, http.StatusUnauthorized)

}

func (suite *AuthroizeTestSuite) TestLoginAccess() {
	reqCreatePotofolio := apis.RequestCreatePotofolio{}

	bytes, err := json.Marshal(reqCreatePotofolio)
	suite.Assert().Nil(err)

	req, err := http.NewRequest(http.MethodPost, suite.getUrl()+"/api/potofolio", strings.NewReader(string(bytes)))
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access-token",
		Value: suite.testAccessToken,
	})

	client := &http.Client{}
	res, err := client.Do(req)

	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, http.StatusOK)
}

func (suite *AuthroizeTestSuite) TestExpireAccessToken() {
	reqCreatePotofolio := apis.RequestCreatePotofolio{}

	bytes, err := json.Marshal(reqCreatePotofolio)
	suite.Assert().Nil(err)

	req, err := http.NewRequest(http.MethodPost, suite.getUrl()+"/api/potofolio", strings.NewReader(string(bytes)))
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access-token",
		Value: suite.testAccessToken,
	})

	client := &http.Client{}
	res, err := client.Do(req)

	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, http.StatusOK)
}

func (suite *AuthroizeTestSuite) TestExpireBothToken() {
	reqCreatePotofolio := apis.RequestCreatePotofolio{}

	bytes, err := json.Marshal(reqCreatePotofolio)
	suite.Assert().Nil(err)

	req, err := http.NewRequest(http.MethodPost, suite.getUrl()+"/api/potofolio", strings.NewReader(string(bytes)))
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access-token",
		Value: suite.testAccessToken,
	})

	client := &http.Client{}
	res, err := client.Do(req)

	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, http.StatusUnauthorized)
}

func (suite *AuthroizeTestSuite) TestLoginAfterLogoutState() {

	req, err := http.NewRequest(http.MethodGet, suite.getUrl()+"/api/logout", nil)
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access-token",
		Value: suite.testAccessToken,
	})

	client := &http.Client{}
	res, err := client.Do(req)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, http.StatusOK)

	tokenValue := getCookieValue(res, "access-token")
	suite.Assert().Empty(tokenValue)

	reqCreatePotofolio := apis.RequestCreatePotofolio{}

	bytes, err := json.Marshal(reqCreatePotofolio)
	suite.Assert().Nil(err)

	req, err = http.NewRequest(http.MethodPost, suite.getUrl()+"/api/potofolio", strings.NewReader(string(bytes)))
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access-token",
		Value: tokenValue,
	})

	client = &http.Client{}
	res, err = client.Do(req)

	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, http.StatusUnauthorized)
}

func TestAuthroizeTestSuite(t *testing.T) {
	suite.Run(t, new(AuthroizeTestSuite))
}
