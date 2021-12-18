package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareTestPotofolioRepo() (*DBConnection, *PotofolioRepository, error) {
	dbConnection := getMemoryDbConnect()

	imageRepo := &ImageRepository{DBConnect: dbConnection}
	imageRepo.CreateTable()

	potofolioRepo := &PotofolioRepository{DBConnect: dbConnection, ImageRepo: imageRepo}
	err := potofolioRepo.CreateTable()
	if err != nil {
		return nil, nil, err
	}

	return dbConnection, potofolioRepo, nil
}

func TestDBConnect(t *testing.T) {

	dbConnect := DBConnection{}
	err := dbConnect.Open("test.db")
	defer dbConnect.Close()

	assert.Nil(t, err, "dbConnect.Open 실패")
}

func TestPotofolioCreateTable(t *testing.T) {
	dbConnection := getMemoryDbConnect()
	defer dbConnection.Close()

	potofolioRepo := PotofolioRepository{DBConnect: dbConnection}
	err := potofolioRepo.CreateTable()
	assert.Nil(t, err, "potofolioRepo CreateTable err 발생")
}

func TestAddPotofolio(t *testing.T) {

	dbConnection, repo, err := prepareTestPotofolioRepo()
	assert.Nil(t, err, "getPotofolioTestPrepare() err 발생")

	defer dbConnection.Close()

	_, err = repo.AddPotofolio("test potooflio", []string{"image1", "image2", "image3"})
	assert.Nil(t, err, "AddPotofolio() err 발생")

	findPotofolio, err := repo.FindPotofolio(1)
	assert.Nil(t, err, "FindPotofolio() err 발생")

	assert.Equal(t, findPotofolio.Title, "test potooflio")
	assert.Equal(t, findPotofolio.Images[0].Path, "image1")
	assert.Equal(t, findPotofolio.Images[1].Path, "image2")
	assert.Equal(t, findPotofolio.Images[2].Path, "image3")
}

func prepareTestExistPotofolioRepo() (*DBConnection, *PotofolioRepository, error) {
	dbConnection, repo, err := prepareTestPotofolioRepo()
	if err != nil {
		return nil, nil, err
	}

	_, err = repo.AddPotofolio("test title1", []string{"image1-1", "image1-2", "image1-3"})
	if err != nil {
		return nil, nil, err
	}

	_, err = repo.AddPotofolio("test title2", []string{"image2-1", "image2-2", "image2-3"})
	if err != nil {
		return nil, nil, err
	}

	_, err = repo.AddPotofolio("test title3", []string{"image3-1", "image3-2", "image3-3"})
	if err != nil {
		return nil, nil, err
	}

	return dbConnection, repo, nil
}

func TestGetPotofolios(t *testing.T) {

	dbConnection, repo, err := prepareTestExistPotofolioRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	potofolies, err := repo.GetPotofolioList()
	assert.Nil(t, err, "GetPotofolioList() err")

	assert.Equal(t, len(potofolies), 3)

	assert.Equal(t, potofolies[0].Title, "test title1")
	assert.Equal(t, potofolies[0].Images[0].Path, "image1-1")
	assert.Equal(t, potofolies[0].Images[1].Path, "image1-2")
	assert.Equal(t, potofolies[0].Images[2].Path, "image1-3")
}

func TestUpdatePotofolio(t *testing.T) {
	dbConnection, repo, err := prepareTestExistPotofolioRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	_, err = repo.FindPotofolio(1)
	assert.Nil(t, err, "FindPotofolio(1) err")

	imageModel, err := repo.ImageRepo.FindImageFromPath(PotofolioType, 1, "image1-2")
	assert.Nil(t, err)
	assert.NotNil(t, imageModel)

	updateTitle := "update title"

	removeImages, err := repo.UpdatePotofolio(1, &updateTitle, []int64{imageModel.Id}, []string{"update image1-1"})
	assert.Nil(t, err, "UpdatePotofolio(1) err")

	assert.Equal(t, len(removeImages), 1)
	assert.Equal(t, removeImages[0], "image1-2")

	potofolio, _ := repo.FindPotofolio(1)
	assert.Equal(t, potofolio.Title, "update title")
	assert.Equal(t, len(potofolio.Images), 3)
	assert.Equal(t, potofolio.Images[0].Path, "image1-1")
	assert.Equal(t, potofolio.Images[1].Path, "image1-3")
	assert.Equal(t, potofolio.Images[2].Path, "update image1-1")
}

func TestUpdatePartialPotofolio_1(t *testing.T) {
	dbConnection, repo, err := prepareTestExistPotofolioRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	_, err = repo.FindPotofolio(1)
	assert.Nil(t, err, "FindPotofolio(1) err")

	imageModel, err := repo.ImageRepo.FindImageFromPath(PotofolioType, 1, "image1-2")
	assert.Nil(t, err)
	assert.NotNil(t, imageModel)

	removeImages, err := repo.UpdatePotofolio(1, nil, []int64{imageModel.Id}, nil)
	assert.Nil(t, err, "UpdatePotofolio(1) err")

	assert.Equal(t, len(removeImages), 1)
	assert.Equal(t, removeImages[0], "image1-2")

	potofolio, _ := repo.FindPotofolio(1)
	assert.Equal(t, potofolio.Title, "test title1")
	assert.Equal(t, len(potofolio.Images), 2)
	assert.Equal(t, potofolio.Images[0].Path, "image1-1")
	assert.Equal(t, potofolio.Images[1].Path, "image1-3")
}

func TestUpdatePartialPotofolio_2(t *testing.T) {
	dbConnection, repo, err := prepareTestExistPotofolioRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	_, err = repo.FindPotofolio(1)
	assert.Nil(t, err, "FindPotofolio(1) err")

	imageModel, err := repo.ImageRepo.FindImageFromPath(PotofolioType, 1, "image1-2")
	assert.Nil(t, err)
	assert.NotNil(t, imageModel)

	updateTitle := "update title"
	removeImages, err := repo.UpdatePotofolio(1, &updateTitle, nil, []string{"add Image1-1"})
	assert.Nil(t, err, "UpdatePotofolio(1) err")

	assert.Equal(t, len(removeImages), 0)

	potofolio, _ := repo.FindPotofolio(1)
	assert.Equal(t, potofolio.Title, "update title")
	assert.Equal(t, len(potofolio.Images), 4)
	assert.Equal(t, potofolio.Images[0].Path, "image1-1")
	assert.Equal(t, potofolio.Images[1].Path, "image1-2")
	assert.Equal(t, potofolio.Images[2].Path, "image1-3")
	assert.Equal(t, potofolio.Images[3].Path, "add Image1-1")
}

func TestRemovePotofolio(t *testing.T) {

	dbConnection, repo, err := prepareTestExistPotofolioRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	err = repo.RemovePotofolio(2)
	assert.Nil(t, err, "RemovePotofolio(2) err")

	findPotofolio, err := repo.FindPotofolio(2)
	assert.Nil(t, err, "FindPotofolio() err")

	assert.Nil(t, findPotofolio)
}
