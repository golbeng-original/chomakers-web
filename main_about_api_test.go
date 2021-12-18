package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/golbeng-original/chomakers-web/apis"
	"github.com/golbeng-original/chomakers-web/models"
)

type AboutTestApiSuite struct {
	suite.Suite
	dbConnection    *models.DBConnection
	aboutRepository *models.AboutRepository
	testServer      *httptest.Server
}

func (suite *AboutTestApiSuite) getUrl() string {
	return suite.testServer.URL
}

func (suite *AboutTestApiSuite) SetupSuite() {

	dbConnection := models.DBConnection{}
	dbConnection.Open("file::memory:?mode=memory&cache=shared")

	suite.dbConnection = &dbConnection

	repositoryConfigure := &models.RepositoryConfigure{}
	repositoryConfigure.Init(&dbConnection)
	repositoryConfigure.IsCheckAuthorize = false

	repositoryConfigure.AboutRepository.CreateTable()

	suite.aboutRepository = repositoryConfigure.AboutRepository

	suite.testServer = httptest.NewServer(Setup(repositoryConfigure, "./assets/images"))
}

func (suite *AboutTestApiSuite) TearDownSuite() {
	suite.testServer.Close()
	suite.dbConnection.Close()
}

func (suite *AboutTestApiSuite) BeforeTest(suiteName, testName string) {

	if testName == "TestAboutGetApi" ||
		testName == "TestAboutPostApi_3" ||
		testName == "TestAboutHistoryPostApi_2" ||
		testName == "TestAboutHistoryPostApi_3" {

		profileImage := "/iamge/a0.jpg"
		profileName := "test name"
		profileContact := "test contact"
		profileIntroduce := "test introduce"

		suite.aboutRepository.UpdateAbout(&profileImage, &profileName, &profileContact, &profileIntroduce)

		addHistoryInfos := make([]models.AboutHistoryContent, 0)
		addHistoryInfos = append(addHistoryInfos, models.AboutHistoryContent{
			Category: "category 1",
			Duration: "duraction 1",
			Content:  "content 1",
		})

		addHistoryInfos = append(addHistoryInfos, models.AboutHistoryContent{
			Category: "category 2",
			Duration: "duraction 2",
			Content:  "content 2",
		})

		addHistoryInfos = append(addHistoryInfos, models.AboutHistoryContent{
			Category: "category 3",
			Duration: "duraction 3",
			Content:  "content 3",
		})

		suite.aboutRepository.UpdateAboutHistory(nil, nil, addHistoryInfos)
	}

}

func (suite *AboutTestApiSuite) TestAboutGetApi() {

	res, err := http.Get(suite.getUrl() + "/api/about")
	suite.Assert().Nil(err)

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responseJson apis.ResponsePresent
	err = json.Unmarshal(bytes, &responseJson)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseJson.Result, "success")

	var responseData apis.ResponseAbout
	err = json.Unmarshal([]byte(responseJson.Data), &responseData)
	suite.Assert().Nil(err)

	fmt.Println(responseData)
}

// Test Post /api/essay
func (suite *AboutTestApiSuite) TestAboutPostApi_1() {

	bytes, err := ioutil.ReadFile("./assets/test_images/a2.jpg")
	suite.Assert().Nil(err)

	base64Bytes := base64.StdEncoding.EncodeToString(bytes)
	suite.Assert().NotEmpty(base64Bytes)

	profileName := "test Name 이름"
	profileContact := "연락처"
	profileIntroduce := "자기 소개!!!<br>자기 소개!!"

	reqAbout := apis.RequestUpdateAbout{
		ProfileImage: &apis.RequestSaveImage{
			Filename: "a2.jpg",
			Data:     base64Bytes,
		},
		ProfileName:      &profileName,
		Contact:          &profileContact,
		IntroduceContent: &profileIntroduce,
	}

	requestBytes, err := json.Marshal(reqAbout)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(requestBytes))

	res, err := http.Post(suite.getUrl()+"/api/about", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseAbout apis.ResponseAbout
	err = json.Unmarshal([]byte(responsePresent.Data), &responseAbout)
	suite.Assert().Nil(err)

	fmt.Println("Post Request About Result")
	fmt.Println(responseAbout)

	suite.Assert().NotEmpty(responseAbout.ProfileImage)
	suite.Assert().Equal(responseAbout.ProfileName, profileName)
	suite.Assert().Equal(responseAbout.Contact, profileContact)
	suite.Assert().Equal(responseAbout.IntroduceContent, profileIntroduce)

}

func (suite *AboutTestApiSuite) TestAboutPostApi_2() {
	bytes, err := ioutil.ReadFile("./assets/test_images/a2.jpg")
	suite.Assert().Nil(err)

	base64Bytes := base64.StdEncoding.EncodeToString(bytes)
	suite.Assert().NotEmpty(base64Bytes)

	//profileName := "test Name 이름"
	//profileContact := "연락처"
	profileIntroduce := "자기 소개!!!<br>자기 소개!!"

	reqAbout := apis.RequestUpdateAbout{
		ProfileImage: &apis.RequestSaveImage{
			Filename: "a2.jpg",
			Data:     base64Bytes,
		},
		ProfileName:      nil,
		Contact:          nil,
		IntroduceContent: &profileIntroduce,
	}

	requestBytes, err := json.Marshal(reqAbout)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(requestBytes))

	res, err := http.Post(suite.getUrl()+"/api/about", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseAbout apis.ResponseAbout
	err = json.Unmarshal([]byte(responsePresent.Data), &responseAbout)
	suite.Assert().Nil(err)

	fmt.Println("Post Request About Result")
	fmt.Println(responseAbout)

	suite.Assert().NotEmpty(responseAbout.ProfileImage)
	suite.Assert().Empty(responseAbout.ProfileName)
	suite.Assert().Empty(responseAbout.Contact)
	suite.Assert().Equal(responseAbout.IntroduceContent, profileIntroduce)
}

func (suite *AboutTestApiSuite) TestAboutPostApi_3() {

	profileIntroduce := "자기 소개!!!<br>자기 소개!!"

	reqAbout := apis.RequestUpdateAbout{
		ProfileImage:     nil,
		ProfileName:      nil,
		Contact:          nil,
		IntroduceContent: &profileIntroduce,
	}

	requestBytes, err := json.Marshal(reqAbout)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(requestBytes))

	res, err := http.Post(suite.getUrl()+"/api/about", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseAbout apis.ResponseAbout
	err = json.Unmarshal([]byte(responsePresent.Data), &responseAbout)
	suite.Assert().Nil(err)

	fmt.Println("Post Request About Result")
	fmt.Println(responseAbout)

	suite.Assert().Equal(responseAbout.ProfileImage, "/iamge/a0.jpg")
	suite.Assert().Equal(responseAbout.ProfileName, "test name")
	suite.Assert().Equal(responseAbout.Contact, "test contact")
	suite.Assert().Equal(responseAbout.IntroduceContent, profileIntroduce)
}

func (suite *AboutTestApiSuite) TestAboutHistoryPostApi_1() {

	reqAboutHistory := apis.RequestUpdateAboutHistory{
		RemoveIds: nil,
		AppendHistories: []apis.RequesAboutHistoryElement{
			{
				Category: "category1",
				Duration: "duraction1",
				Content:  "content1",
			},
			{
				Category: "category2",
				Duration: "duraction2",
				Content:  "content2",
			},
			{
				Category: "category3",
				Duration: "duraction3",
				Content:  "content3",
			},
			{
				Category: "category4",
				Duration: "duraction4",
				Content:  "content4",
			},
		},
	}

	requestBytes, err := json.Marshal(reqAboutHistory)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(requestBytes))
	suite.Assert().NotNil(requestReader)

	res, err := http.Post(suite.getUrl()+"/api/about-history", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseAbout apis.ResponseAbout
	err = json.Unmarshal([]byte(responsePresent.Data), &responseAbout)
	suite.Assert().Nil(err)

	fmt.Println("Post Request About Result")
	fmt.Println(responseAbout)

}

func (suite *AboutTestApiSuite) TestAboutHistoryPostApi_2() {

	reqAboutHistory := apis.RequestUpdateAboutHistory{
		RemoveIds: []int64{
			1,
		},
		AppendHistories: []apis.RequesAboutHistoryElement{
			{
				Category: "new category1",
				Duration: "new duraction1",
				Content:  "new content1",
			},
			{
				Category: "new category2",
				Duration: "new duraction2",
				Content:  "new content2",
			},
			{
				Category: "new category3",
				Duration: "new duraction3",
				Content:  "new content3",
			},
			{
				Category: "new category4",
				Duration: "new duraction4",
				Content:  "new content4",
			},
		},
		UpdateHistories: nil,
	}

	requestBytes, err := json.Marshal(reqAboutHistory)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(requestBytes))
	suite.Assert().NotNil(requestReader)

	res, err := http.Post(suite.getUrl()+"/api/about-history", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseAbout apis.ResponseAbout
	err = json.Unmarshal([]byte(responsePresent.Data), &responseAbout)
	suite.Assert().Nil(err)

	fmt.Println("Post Request About Result")
	fmt.Println(responseAbout)
}

func (suite *AboutTestApiSuite) TestAboutHistoryPostApi_3() {
	reqAboutHistory := apis.RequestUpdateAboutHistory{
		RemoveIds: []int64{
			1,
		},
		AppendHistories: []apis.RequesAboutHistoryElement{
			{
				Category: "new category1",
				Duration: "new duraction1",
				Content:  "new content1",
			},
			{
				Category: "new category2",
				Duration: "new duraction2",
				Content:  "new content2",
			},
			{
				Category: "new category3",
				Duration: "new duraction3",
				Content:  "new content3",
			},
			{
				Category: "new category4",
				Duration: "new duraction4",
				Content:  "new content4",
			},
		},
		UpdateHistories: []apis.RequestUpdateAboutHistoryElement{
			{
				Id: 2,
				RequesAboutHistoryElement: apis.RequesAboutHistoryElement{
					Category: "change category 2",
					Duration: "change duration 2",
					Content:  "change content 2",
				},
			},
		},
	}

	requestBytes, err := json.Marshal(reqAboutHistory)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(requestBytes))
	suite.Assert().NotNil(requestReader)

	res, err := http.Post(suite.getUrl()+"/api/about-history", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseAbout apis.ResponseAbout
	err = json.Unmarshal([]byte(responsePresent.Data), &responseAbout)
	suite.Assert().Nil(err)

	fmt.Println("Post Request About Result")
	fmt.Println(responseAbout)
}

func TestAboutTestApiSuite(t *testing.T) {
	suite.Run(t, &AboutTestApiSuite{})
}
