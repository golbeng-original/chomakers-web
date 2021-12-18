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

type PotofolioTestApiSuite struct {
	suite.Suite
	dbConnection *models.DBConnection

	testServer *httptest.Server
}

func (suite *PotofolioTestApiSuite) getUrl() string {
	return suite.testServer.URL
}

func (suite *PotofolioTestApiSuite) SetupSuite() {

	dbConnection := models.DBConnection{}
	dbConnection.Open("file::memory:?mode=memory&cache=shared")

	suite.dbConnection = &dbConnection

	repositoryConfigure := &models.RepositoryConfigure{}
	repositoryConfigure.Init(&dbConnection)
	repositoryConfigure.IsCheckAuthorize = false

	potofolioRepo := repositoryConfigure.PotofolioRepository

	// shuite Potofolios
	potofolioRepo.AddPotofolio("potofolio1", []string{"/images/a0.jpg", "/images/a1.jpg", "/images/a2.jpg"})
	potofolioRepo.AddPotofolio("potofolio2", []string{"/images/a1.jpg", "/images/a2.jpg"})
	potofolioRepo.AddPotofolio("potofolio3", []string{"/images/a3.jpg"})
	potofolioRepo.AddPotofolio("potofolio4", nil)

	suite.testServer = httptest.NewServer(Setup(repositoryConfigure, "./assets/images"))
}

func (suite *PotofolioTestApiSuite) TearDownSuite() {
	suite.testServer.Close()
	suite.dbConnection.Close()
}

func (suite *PotofolioTestApiSuite) getPotofolioGetList() []apis.ResponsePotofolioElement {
	res, err := http.Get(suite.getUrl() + "/api/potofolio")
	suite.Assert().Nil(err)

	defer res.Body.Close()
	suite.Assert().Equal(res.StatusCode, 200)

	bytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responseJson apis.ResponsePresent
	err = json.Unmarshal(bytes, &responseJson)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseJson.Result, "success")

	var responseData apis.ResponsePotofolioList
	err = json.Unmarshal([]byte(responseJson.Data), &responseData)
	suite.Assert().Nil(err)

	return responseData.List
}

// Test GET /api/potofolio
func (suite *PotofolioTestApiSuite) TestPtofolioGetListApi() {

	potofolios := suite.getPotofolioGetList()

	suite.Assert().Equal(len(potofolios), 4)
	suite.Assert().Equal(potofolios[0].Id, int64(1))
	suite.Assert().Equal(potofolios[0].Title, "potofolio1")
	suite.Assert().Equal(len(potofolios[0].Images), 3)
	suite.Assert().Equal(potofolios[0].Images[0].ImageUrl, "/images/a0.jpg")
}

// Test GET /api/potofolio/:id
func (suite *PotofolioTestApiSuite) TestPtofolioGetApi() {

	res, err := http.Get(suite.getUrl() + "/api/potofolio/1")
	suite.Assert().Nil(err)

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responseJson apis.ResponsePresent
	err = json.Unmarshal(bytes, &responseJson)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responseJson.Result, "success")

	var responseData apis.ResponsePotofolioElement
	err = json.Unmarshal([]byte(responseJson.Data), &responseData)
	suite.Assert().Nil(err)

	suite.Assert().Equal(responseData.Id, int64(1))
	suite.Assert().Equal(responseData.Title, "potofolio1")
	suite.Assert().Equal(len(responseData.Images), 3)
	suite.Assert().Equal(responseData.Images[0].ImageUrl, "/images/a0.jpg")
}

// Test Post /api/potofolio
func (suite *PotofolioTestApiSuite) TestPtofolioPostApi() {

	requestRawImages := []string{"./assets/test_images/a0.jpg", "./assets/test_images/a1.jpg", "./assets/test_images/a2.jpg"}

	requestImages := make([]apis.RequestSaveImage, 0)
	for _, rawImage := range requestRawImages {
		bytes, err := ioutil.ReadFile(rawImage)
		suite.Assert().Nil(err)

		base64Bytes := base64.StdEncoding.EncodeToString(bytes)
		suite.Assert().NotEmpty(base64Bytes)

		_, rawImageFileName := filepath.Split(rawImage)

		requestImages = append(requestImages, apis.RequestSaveImage{Filename: rawImageFileName, Data: base64Bytes})
	}

	requestPotofolio := apis.RequestCreatePotofolio{
		Title:  "new Potofolio",
		Images: requestImages,
	}

	bytes, err := json.Marshal(requestPotofolio)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(bytes))

	res, err := http.Post(suite.getUrl()+"/api/potofolio", "application/json", requestReader)
	suite.Assert().Nil(err)
	suite.Assert().Equal(res.StatusCode, 200)

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	suite.Assert().Nil(err)

	var responsePresent apis.ResponsePresent
	err = json.Unmarshal(bodyBytes, &responsePresent)
	suite.Assert().Nil(err)
	suite.Assert().Equal(responsePresent.Result, "success")

	var responsePotofolio apis.ResponsePotofolioElement
	err = json.Unmarshal([]byte(responsePresent.Data), &responsePotofolio)
	suite.Assert().Nil(err)

	fmt.Println("Create potofolio Result")
	fmt.Println(responsePotofolio)

	responseList := suite.getPotofolioGetList()
	fmt.Println("confirm potofolio list")
	fmt.Println(responseList)
}

// Test Put /api/potofolio/:id
func (suite *PotofolioTestApiSuite) TestPtofolioPutApi() {

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

	title := "update Potofolio"

	//2, 4, 5
	requestPotofolioUpdate := apis.RequestUpdatePotofolio{
		Title:          &title,
		RemoveImageIds: []int64{4, 5},
		AddImages:      requestImages,
	}

	bytes, err := json.Marshal(requestPotofolioUpdate)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(bytes))

	req, err := http.NewRequest(http.MethodPut, suite.getUrl()+"/api/potofolio/2", requestReader)
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

	var responsePotofolioElement apis.ResponsePotofolioElement
	err = json.Unmarshal([]byte(responsePresent.Data), &responsePotofolioElement)
	suite.Assert().Nil(err)
}

// Test Put(Partial) /api/potofolio/:id
func (suite *PotofolioTestApiSuite) TestPtofolioPutPartialApi() {

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

	//2, 4, 5
	requestPotofolioUpdate := apis.RequestUpdatePotofolio{
		Title:          nil,
		RemoveImageIds: nil,
		AddImages:      requestImages,
	}

	bytes, err := json.Marshal(requestPotofolioUpdate)
	suite.Assert().Nil(err)

	requestReader := strings.NewReader(string(bytes))

	req, err := http.NewRequest(http.MethodPut, suite.getUrl()+"/api/potofolio/2", requestReader)
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

	var responsePotofolioElement apis.ResponsePotofolioElement
	err = json.Unmarshal([]byte(responsePresent.Data), &responsePotofolioElement)
	suite.Assert().Nil(err)
}

// Test DEL /api/potofolio/:id
func (suite *PotofolioTestApiSuite) TestPtofolioDeleteApi() {

	req, err := http.NewRequest(http.MethodDelete, suite.getUrl()+"/api/potofolio/2", nil)
	suite.Assert().Nil(err)

	client := http.Client{}
	client.Do(req)

	responseList := suite.getPotofolioGetList()
	fmt.Println("confirm potofolio list")
	fmt.Println(responseList)
}

func TestPotoflioTestApiSuite(t *testing.T) {
	suite.Run(t, new(PotofolioTestApiSuite))
}
