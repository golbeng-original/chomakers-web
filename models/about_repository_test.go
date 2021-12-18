package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AboutRepositryTestSuite struct {
	suite.Suite
	dbConnect *DBConnection

	aboutRepository *AboutRepository
}

func (suite *AboutRepositryTestSuite) SetupSuite() {

	suite.dbConnect = &DBConnection{}
	suite.dbConnect.Open("file::memory:?mode=memory&cache=shared")

	suite.aboutRepository = &AboutRepository{DBConnect: suite.dbConnect}
	suite.aboutRepository.CreateTable()
}

func (suite *AboutRepositryTestSuite) TearDownSuite() {
	suite.dbConnect.Close()
}

func (suite *AboutRepositryTestSuite) BeforeTest(suiteName, testName string) {

	if testName == "TestGetAboutAll" {
		profileImage := "/image/profile.jpg"
		profileName := "name data"
		contact := "contact data"
		introduceContent := "introduce data"

		suite.aboutRepository.UpdateAbout(&profileImage, &profileName, &contact, &introduceContent)
	} else if testName == "TestGetAboutPartial" {
		profileImage := "/image/profile.jpg"
		introduceContent := "introduce data"

		suite.aboutRepository.UpdateAbout(&profileImage, nil, nil, &introduceContent)
	} else if testName == "TestUpdateAbout_2" {
		profileImage := "/image/profile.jpg"
		introduceContent := "introduce data"

		suite.aboutRepository.UpdateAbout(&profileImage, nil, nil, &introduceContent)
	} else if testName == "TestGetAboutHistory" ||
		testName == "TestUpdateAboutHistory_1" ||
		testName == "TestUpdateAboutHistory_2" ||
		testName == "TestUpdateAboutHistory_3" {

		addHistories := make([]AboutHistoryContent, 0)
		addHistories = append(addHistories, AboutHistoryContent{
			Category: "category1",
			Duration: "durection1",
			Content:  "content1",
		})
		addHistories = append(addHistories, AboutHistoryContent{
			Category: "category2",
			Duration: "durection2",
			Content:  "content2",
		})
		addHistories = append(addHistories, AboutHistoryContent{
			Category: "category3",
			Duration: "durection3",
			Content:  "content3",
		})

		suite.aboutRepository.UpdateAboutHistory(nil, nil, addHistories)
	}
}

func (suite *AboutRepositryTestSuite) TestGetAboutEmpty() {
	aboutModel, err := suite.aboutRepository.GetAbout()
	suite.Nil(err)

	suite.Nil(aboutModel.ProfileImage)
	suite.Nil(aboutModel.ProfileName)
	suite.Nil(aboutModel.Contact)
	suite.Nil(aboutModel.IntroduceContent)
}

func (suite *AboutRepositryTestSuite) TestGetAboutAll() {

	aboutModel, err := suite.aboutRepository.GetAbout()
	suite.Nil(err)

	suite.Equal(*aboutModel.ProfileImage, "/image/profile.jpg")
	suite.Equal(*aboutModel.ProfileName, "name data")
	suite.Equal(*aboutModel.Contact, "contact data")
	suite.Equal(*aboutModel.IntroduceContent, "introduce data")
}

func (suite *AboutRepositryTestSuite) TestGetAboutPartial() {
	aboutModel, err := suite.aboutRepository.GetAbout()
	suite.Nil(err)

	suite.Equal(*aboutModel.ProfileImage, "/image/profile.jpg")
	suite.Nil(aboutModel.ProfileName)
	suite.Nil(aboutModel.Contact)
	suite.Equal(*aboutModel.IntroduceContent, "introduce data")
}

func (suite *AboutRepositryTestSuite) TestUpdateAbout_1() {

	updateProfileIamge := "update ProfileImage"
	updateProfileName := "update ProfileName"
	updateContact := "update Contact"
	updateintroduce := "update Introduce"

	_, err := suite.aboutRepository.UpdateAbout(&updateProfileIamge, &updateProfileName, &updateContact, &updateintroduce)
	suite.Nil(err)

	aboutModel, err := suite.aboutRepository.GetAbout()
	suite.Nil(err)

	suite.Equal(*aboutModel.ProfileImage, updateProfileIamge)
	suite.Equal(*aboutModel.ProfileName, updateProfileName)
	suite.Equal(*aboutModel.Contact, updateContact)
	suite.Equal(*aboutModel.IntroduceContent, updateintroduce)
}

func (suite *AboutRepositryTestSuite) TestUpdateAbout_2() {

	updateProfileImage := "update ProfileImage"
	updateProfileName := "update ProfileName"
	updateContact := "update Contact"
	updateintroduce := "update Introduce"

	prevPorfileImage, err := suite.aboutRepository.UpdateAbout(&updateProfileImage, &updateProfileName, &updateContact, &updateintroduce)
	suite.Nil(err)
	suite.Equal(*prevPorfileImage, "/image/profile.jpg")

	aboutModel, err := suite.aboutRepository.GetAbout()
	suite.Nil(err)

	suite.Equal(*aboutModel.ProfileImage, updateProfileImage)
	suite.Equal(*aboutModel.ProfileName, updateProfileName)
	suite.Equal(*aboutModel.Contact, updateContact)
	suite.Equal(*aboutModel.IntroduceContent, updateintroduce)
}

func (suite *AboutRepositryTestSuite) TestGetAboutHistory() {

	result, err := suite.aboutRepository.GetHistory()
	suite.Nil(err)
	suite.Equal(len(result), 3)

	suite.Equal(result[0].Id, int64(1))
	suite.Equal(result[0].Category, "category1")
	suite.Equal(result[0].Duration, "durection1")
	suite.Equal(result[0].Content, "content1")

	fmt.Printf("%v\n", result)
}

func (suite *AboutRepositryTestSuite) TestUpdateAboutHistory_1() {

	suite.aboutRepository.UpdateAboutHistory([]int64{1, 2}, nil, nil)

	result, err := suite.aboutRepository.GetHistory()
	suite.Nil(err)
	suite.Equal(len(result), 1)

	suite.Equal(result[0].Id, int64(3))
	suite.Equal(result[0].Category, "category3")
	suite.Equal(result[0].Duration, "durection3")
	suite.Equal(result[0].Content, "content3")

	fmt.Printf("%v\n", result)
}

func (suite *AboutRepositryTestSuite) TestUpdateAboutHistory_2() {

	addHistories := make([]AboutHistoryContent, 0)
	addHistories = append(addHistories, AboutHistoryContent{
		Category: "new category1",
		Duration: "new durection1",
		Content:  "new content1",
	})

	suite.aboutRepository.UpdateAboutHistory([]int64{1, 2}, nil, addHistories)

	result, err := suite.aboutRepository.GetHistory()
	suite.Nil(err)
	suite.Equal(len(result), 2)

	suite.Equal(result[0].Id, int64(3))
	suite.Equal(result[0].Category, "category3")
	suite.Equal(result[0].Duration, "durection3")
	suite.Equal(result[0].Content, "content3")

	suite.Equal(result[1].Id, int64(4))
	suite.Equal(result[1].Category, "new category1")
	suite.Equal(result[1].Duration, "new durection1")
	suite.Equal(result[1].Content, "new content1")

	fmt.Printf("%v\n", result)
}

func (suite *AboutRepositryTestSuite) TestUpdateAboutHistory_3() {

	addHistories := make([]AboutHistoryContent, 0)
	addHistories = append(addHistories, AboutHistoryContent{
		Category: "new category1",
		Duration: "new durection1",
		Content:  "new content1",
	})

	updateHistories := make([]AboutHistoryIdContent, 0)
	updateHistories = append(updateHistories, AboutHistoryIdContent{
		Id: 2,
		AboutHistoryContent: AboutHistoryContent{
			Category: "update category2",
			Duration: "update duration2",
			Content:  "update content2",
		},
	})

	suite.aboutRepository.UpdateAboutHistory([]int64{1}, updateHistories, addHistories)

	result, err := suite.aboutRepository.GetHistory()
	suite.Nil(err)
	suite.Equal(len(result), 3)

	suite.Equal(result[0].Id, int64(2))
	suite.Equal(result[0].Category, "update category2")
	suite.Equal(result[0].Duration, "update duration2")
	suite.Equal(result[0].Content, "update content2")

	suite.Equal(result[1].Id, int64(3))
	suite.Equal(result[1].Category, "category3")
	suite.Equal(result[1].Duration, "durection3")
	suite.Equal(result[1].Content, "content3")

	suite.Equal(result[2].Id, int64(4))
	suite.Equal(result[2].Category, "new category1")
	suite.Equal(result[2].Duration, "new durection1")
	suite.Equal(result[2].Content, "new content1")

	fmt.Printf("%v\n", result)
}

func TestAboutRepositorySuite(t *testing.T) {
	suite.Run(t, new(AboutRepositryTestSuite))
}
