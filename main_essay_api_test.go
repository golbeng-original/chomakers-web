package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"

	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/golbeng-original/chomakers-web/apis"
	"github.com/golbeng-original/chomakers-web/models"
)

type EssayTestApiSuite struct {
	suite.Suite
	dbConnection *models.DBConnection

	testServer      *httptest.Server
	essayRepository *models.EssayRepository
}

func (suite *EssayTestApiSuite) getUrl() string {
	return suite.testServer.URL
}

func (suite *EssayTestApiSuite) SetupSuite() {

	dbConnection := models.DBConnection{}
	dbConnection.Open("file::memory:?mode=memory&cache=shared")

	suite.dbConnection = &dbConnection

	repositoryConfigure := &models.RepositoryConfigure{}
	repositoryConfigure.Init(&dbConnection)
	repositoryConfigure.IsCheckAuthorize = false

	suite.essayRepository = repositoryConfigure.EssayRepository

	// shuite Essay
	suite.essayRepository.AddEssay("essay1", "/image/a0.jpg", "essay content1\nessay content1", []string{"/image/a0.jpg", "/image/a1.jpg", "/image/a2.jpg"})
	suite.essayRepository.AddEssay("essay2", "/image/a3.jpg", "essay content2\nessay content2", []string{"/image/a1.jpg", "/image/a2.jpg"})
	suite.essayRepository.AddEssay("essay3", "/image/a4.jpg", "essay content3\nessay content3", []string{"/image/a2.jpg"})
	suite.essayRepository.AddEssay("essay4", "/image/a5.jpg", "essay content4\nessay content4", nil)

	suite.testServer = httptest.NewServer(Setup(repositoryConfigure, "./assets/images"))
}

func (suite *EssayTestApiSuite) TearDownSuite() {
	suite.testServer.Close()
	suite.dbConnection.Close()
}

func (suite *EssayTestApiSuite) BeforeTest(suiteName, testName string) {

}

func (suite *EssayTestApiSuite) getEssayGetList() []apis.ResponseEssayThumbnailElement {
	res, err := http.Get(suite.getUrl() + "/api/essay")
	suite.Assert().Nil(err)

	defer res.Body.Close()
	suite.Assert().Equal(res.StatusCode, 200)

	bytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responseJson apis.ResponsePresent
	err = json.Unmarshal(bytes, &responseJson)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseJson.Result, "success")

	var responseData []apis.ResponseEssayThumbnailElement
	err = json.Unmarshal([]byte(responseJson.Data), &responseData)
	suite.Assert().Nil(err)

	return responseData
}

// Test GET /api/essay
func (suite *EssayTestApiSuite) TestEssayGetListApi() {

	potofolios := suite.getEssayGetList()

	suite.Assert().Equal(len(potofolios), 4)
	suite.Assert().Equal(potofolios[0].Id, int64(1))
	suite.Assert().Equal(potofolios[0].Title, "essay1")
	suite.Assert().Equal(potofolios[0].ThumbnailImage, "/image/a0.jpg")
}

// Test GET /api/essay/:id
func (suite *EssayTestApiSuite) TestEssayGetApi() {

	res, err := http.Get(suite.getUrl() + "/api/essay/1")
	suite.Assert().Nil(err)

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responseJson apis.ResponsePresent
	err = json.Unmarshal(bytes, &responseJson)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseJson.Result, "success")

	var responseData apis.ResponseEssayElement
	err = json.Unmarshal([]byte(responseJson.Data), &responseData)
	suite.Assert().Nil(err)

	suite.Assert().Equal(responseData.Id, int64(1))
	suite.Assert().Equal(responseData.Title, "essay1")
	suite.Assert().Equal(responseData.EssayContent, "essay content1\nessay content1")
	suite.Assert().Equal(responseData.ThumbnailImage, "/image/a0.jpg")
	suite.Assert().Equal(len(responseData.Images), 3)
	suite.Assert().Equal(responseData.Images[0].ImageUrl, "/image/a0.jpg")
}

// Test Post /api/essay
func (suite *EssayTestApiSuite) TestEssayPostApi() {

	bytes, err := ioutil.ReadFile("./assets/test_images/a2.jpg")
	suite.Assert().Nil(err)

	base64Bytes := base64.StdEncoding.EncodeToString(bytes)
	suite.Assert().NotEmpty(base64Bytes)

	requestThumbnailImage := apis.RequestSaveImage{
		Filename: "a2.jpg",
		Data:     base64Bytes,
	}

	requestRawImages := []string{
		"./assets/test_images/a0.jpg",
		"./assets/test_images/a1.jpg",
		"./assets/test_images/a2.jpg"}

	requestImages := make([]apis.RequestSaveImage, 0)
	for _, rawImage := range requestRawImages {
		bytes, err = ioutil.ReadFile(rawImage)
		suite.Assert().Nil(err)

		base64Bytes = base64.StdEncoding.EncodeToString(bytes)
		suite.Assert().NotEmpty(base64Bytes)

		_, rawImageFileName := filepath.Split(rawImage)

		requestImages = append(requestImages, apis.RequestSaveImage{Filename: rawImageFileName, Data: base64Bytes})
	}

	requestCreateEssay := apis.RequestCreateEssay{
		Title:          "new Essay",
		EssayContent:   "new Essay Content<b>sdfsd</b></br>new new new new essay",
		ThumbnailImage: requestThumbnailImage,
		Images:         requestImages,
	}

	bytes, err = json.Marshal(requestCreateEssay)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(bytes))

	res, err := http.Post(suite.getUrl()+"/api/essay", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseEssay apis.ResponseEssayElement
	err = json.Unmarshal([]byte(responsePresent.Data), &responseEssay)
	suite.Assert().Nil(err)

	fmt.Println("Create essay Result")
	fmt.Println(responseEssay)

	responseList := suite.getEssayGetList()
	fmt.Println("confirm essay list")
	fmt.Println(responseList)
}

// Test Put /api/essay/:id
func (suite *EssayTestApiSuite) TestEssayPutApi() {

	bytes, err := ioutil.ReadFile("./assets/test_images/a2.jpg")
	suite.Assert().Nil(err)

	base64Bytes := base64.StdEncoding.EncodeToString(bytes)
	suite.Assert().NotEmpty(base64Bytes)

	requestThumbnailImage := apis.RequestSaveImage{
		Filename: "a2.jpg",
		Data:     base64Bytes,
	}

	addRawImages := []string{"./assets/test_images/a0.jpg", "./assets/test_images/a3.jpg"}

	requestImages := make([]apis.RequestSaveImage, 0)
	for _, rawImage := range addRawImages {
		bytes, err := ioutil.ReadFile(rawImage)
		suite.Assert().Nil(err)

		base64Bytes := base64.StdEncoding.EncodeToString(bytes)
		suite.Assert().NotEmpty(base64Bytes)

		_, rawImageFileName := filepath.Split(rawImage)

		requestImages = append(requestImages, apis.RequestSaveImage{Filename: rawImageFileName, Data: base64Bytes})
	}

	title := "update Essay"
	content := "update Essay Content"

	//2, 4, 5
	requestUpdatEssay := apis.RequestUpdateEssay{
		Title:          &title,
		EssayContent:   &content,
		NewThumbnail:   &requestThumbnailImage,
		RemoveImageIds: []int64{4, 5},
		AddImages:      requestImages,
	}

	bytes, err = json.Marshal(requestUpdatEssay)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(bytes))

	req, err := http.NewRequest(http.MethodPut, suite.getUrl()+"/api/essay/2", requestReader)
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	res, err := client.Do(req)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseEssayElement apis.ResponseEssayElement
	err = json.Unmarshal([]byte(responsePresent.Data), &responseEssayElement)
	suite.Assert().Nil(err)
}

// Test Put(Partial) /api/potofolio/:id
func (suite *EssayTestApiSuite) TestPtofolioPutPartialApi() {

	addRawImages := []string{"./assets/test_images/a0.jpg", "./assets/test_images/a3.jpg"}

	requestImages := make([]apis.RequestSaveImage, 0)
	for _, rawImage := range addRawImages {
		bytes, err := ioutil.ReadFile(rawImage)
		suite.Assert().Nil(err)

		base64Bytes := base64.StdEncoding.EncodeToString(bytes)
		suite.Assert().NotEmpty(base64Bytes)

		_, rawImageFileName := filepath.Split(rawImage)

		requestImages = append(requestImages, apis.RequestSaveImage{Filename: rawImageFileName, Data: base64Bytes})
	}

	title := "update Essay"

	//2, 4, 5
	requestUpdatEssay := apis.RequestUpdateEssay{
		Title:          &title,
		EssayContent:   nil,
		NewThumbnail:   nil,
		RemoveImageIds: nil,
		AddImages:      requestImages,
	}

	bytes, err := json.Marshal(requestUpdatEssay)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(bytes))

	req, err := http.NewRequest(http.MethodPut, suite.getUrl()+"/api/essay/2", requestReader)
	suite.Assert().Nil(err)

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	res, err := client.Do(req)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responseEssayElement apis.ResponseEssayElement
	err = json.Unmarshal([]byte(responsePresent.Data), &responseEssayElement)
	suite.Assert().Nil(err)
}

// Test DEL /api/potofolio/:id
func (suite *EssayTestApiSuite) TestEssayDeleteApi() {

	req, err := http.NewRequest(http.MethodDelete, suite.getUrl()+"/api/essay/4", nil)
	suite.Assert().Nil(err)

	client := http.Client{}
	client.Do(req)

	responseList := suite.getEssayGetList()
	fmt.Println("confirm essay list")
	fmt.Println(responseList)
}

func TestEssayTestApiSuite(t *testing.T) {
	suite.Run(t, new(EssayTestApiSuite))
}
