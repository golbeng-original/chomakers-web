package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	dbConnection *DBConnection

	userRepository *UserRespository

	tempPasswordMd5 string
}

func (suite *UserRepositoryTestSuite) SetupSuite() {

	dbConnection := DBConnection{}
	dbConnection.Open("file::memory:?mode=memory&cache=shared")

	suite.dbConnection = &dbConnection

	suite.userRepository = &UserRespository{DBConnect: &dbConnection}
	suite.userRepository.CreateTable()
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	suite.dbConnection.Close()
}

func (suite *UserRepositoryTestSuite) BeforeTest(suiteName, testName string) {

	if testName == "TestAddExistUser" ||
		testName == "TestIsExistUser" ||
		testName == "TestIsEqualUserPassword" {
		suite.userRepository.AddUser("root", "password")

		userModel, _ := suite.userRepository.GetUserModelFromUserName("root")
		suite.tempPasswordMd5 = userModel.Password
	}
}

func (suite *UserRepositoryTestSuite) TestAddUser() {

	err := suite.userRepository.AddUser("root", "password")
	suite.Assert().Nil(err)
}

func (suite *UserRepositoryTestSuite) TestAddExistUser() {
	err := suite.userRepository.AddUser("root", "password")
	suite.Assert().NotNil(err)
}

func (suite *UserRepositoryTestSuite) TestIsExistUser() {
	isExist, err := suite.userRepository.IsExist("root")
	suite.Assert().Nil(err)
	suite.Assert().Equal(isExist, true)
}

func (suite *UserRepositoryTestSuite) TestIsNotExistUser() {
	isExist, err := suite.userRepository.IsExist("root")
	suite.Assert().Nil(err)
	suite.Assert().Equal(isExist, false)
}

func (suite *UserRepositoryTestSuite) TestIsEqualUserPassword() {
	userModel, err := suite.userRepository.GetUserModelFromUserName("root")
	suite.Assert().Nil(err)

	suite.Assert().Equal(userModel.Password, suite.tempPasswordMd5)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
