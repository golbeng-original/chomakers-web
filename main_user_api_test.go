package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/golbeng-original/chomakers-web/apis"
	"github.com/golbeng-original/chomakers-web/models"
)

type UserTestApiSuite struct {
	suite.Suite

	dbConnection *models.DBConnection
	testServer   *httptest.Server

	userRepository   *models.UserRespository
	testUserPassword string
}

func (suite *UserTestApiSuite) getUrl() string {
	return suite.testServer.URL
}

func (suite *UserTestApiSuite) SetupSuite() {
	dbConnection := models.DBConnection{}
	dbConnection.Open("file::memory:?mode=memory&cache=shared")

	suite.dbConnection = &dbConnection

	repositoryConfigure := models.RepositoryConfigure{}
	repositoryConfigure.Init(&dbConnection)
	repositoryConfigure.IsCheckAuthorize = false

	suite.userRepository = repositoryConfigure.UserRepository
	suite.testServer = httptest.NewServer(Setup(&repositoryConfigure, "./assets/images"))
}

func (suite *UserTestApiSuite) TearDownSuite() {
	suite.testServer.Close()
	suite.dbConnection.Close()
}

func (suite *UserTestApiSuite) BeforeTest(suiteName, testName string) {

	if testName == "TestLogin_2" ||
		testName == "TestLogin_3" ||
		testName == "TestLoginWithToken" ||
		testName == "TestLoginWithLogout" {

		suite.userRepository.AddUser("root", "1234")
		userModel, _ := suite.userRepository.GetUserModelFromUserName("root")

		if userModel != nil {
			suite.testUserPassword = userModel.Password
		}

	}

}

func (suite *UserTestApiSuite) TestLogin_1() {

	reqLogin := apis.RequestLogin{
		UserName: "root",
		Password: "1234",
	}

	reqLoginJson, err := json.Marshal(reqLogin)
	suite.Assert().Nil(err)

	reqLoginReader := strings.NewReader(string(reqLoginJson))

	res, err := http.Post(suite.getUrl()+"/api/login", "application/json", reqLoginReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(responseBody, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseData apis.ResponseLogin
	err = json.Unmarshal([]byte(responsePresent.Data), &responseData)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseData.LoginResult, 1)

	fmt.Println(responseData)
}

func (suite *UserTestApiSuite) TestLogin_2() {

	reqLogin := apis.RequestLogin{
		UserName: "root",
		Password: "1234",
	}

	reqLoginJson, err := json.Marshal(reqLogin)
	suite.Assert().Nil(err)

	reqLoginReader := strings.NewReader(string(reqLoginJson))

	res, err := http.Post(suite.getUrl()+"/api/login", "application/json", reqLoginReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(responseBody, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseData apis.ResponseLogin
	err = json.Unmarshal([]byte(responsePresent.Data), &responseData)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseData.LoginResult, 2)

	fmt.Println(responseData)

}

func (suite *UserTestApiSuite) TestLogin_3() {

	reqLogin := apis.RequestLogin{
		UserName: "root",
		Password: suite.testUserPassword,
	}

	reqLoginJson, err := json.Marshal(reqLogin)
	suite.Assert().Nil(err)

	reqLoginReader := strings.NewReader(string(reqLoginJson))

	res, err := http.Post(suite.getUrl()+"/api/login", "application/json", reqLoginReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(responseBody, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseData apis.ResponseLogin
	err = json.Unmarshal([]byte(responsePresent.Data), &responseData)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseData.LoginResult, 0)

	fmt.Println(responseData)

}

func (suite *UserTestApiSuite) TestLoginWithToken() {
	reqLogin := apis.RequestLogin{
		UserName: "root",
		Password: suite.testUserPassword,
	}

	reqLoginJson, err := json.Marshal(reqLogin)
	suite.Assert().Nil(err)

	reqLoginReader := strings.NewReader(string(reqLoginJson))

	res, err := http.Post(suite.getUrl()+"/api/login", "application/json", reqLoginReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(responseBody, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseData apis.ResponseLogin
	err = json.Unmarshal([]byte(responsePresent.Data), &responseData)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseData.LoginResult, 0)

	var tokenValue string

	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "access-token" {
			tokenValue = cookie.Value
		}
	}

	suite.Assert().NotEmpty(tokenValue)

	fmt.Printf("tokenValue = %s\n", tokenValue)

	fmt.Println(responseData)
}

func (suite *UserTestApiSuite) TestLoginWithLogout() {
	reqLogin := apis.RequestLogin{
		UserName: "root",
		Password: suite.testUserPassword,
	}

	reqLoginJson, err := json.Marshal(reqLogin)
	suite.Assert().Nil(err)

	reqLoginReader := strings.NewReader(string(reqLoginJson))

	res, err := http.Post(suite.getUrl()+"/api/login", "application/json", reqLoginReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(responseBody, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseData apis.ResponseLogin
	err = json.Unmarshal([]byte(responsePresent.Data), &responseData)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseData.LoginResult, 0)

	var tokenValue string

	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "access-token" {
			tokenValue = cookie.Value
		}
	}

	suite.Assert().NotEmpty(tokenValue)

	request, err := http.NewRequest("GET", suite.getUrl()+"/api/logout", nil)
	suite.Assert().Nil(err)

	accessCookie := &http.Cookie{
		Name:  "access-token",
		Value: tokenValue,
	}

	request.AddCookie(accessCookie)

	client := &http.Client{}
	res, err = client.Do(request)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)
}

func TestUserTestApiSuite(t *testing.T) {
	suite.Run(t, new(UserTestApiSuite))
}
