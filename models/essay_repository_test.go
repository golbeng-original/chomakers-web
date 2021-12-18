package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareTestEssayRepo() (*DBConnection, *EssayRepository, error) {
	dbConnection := getMemoryDbConnect()

	imageRepo := &ImageRepository{DBConnect: dbConnection}
	imageRepo.CreateTable()

	essayRepo := &EssayRepository{DBConnect: dbConnection, ImageRepo: imageRepo}
	err := essayRepo.CreateTable()
	if err != nil {
		return nil, nil, err
	}

	return dbConnection, essayRepo, nil
}

func TestEssayCreateTable(t *testing.T) {
	dbConnection := getMemoryDbConnect()
	defer dbConnection.Close()

	essayRepo := EssayRepository{DBConnect: dbConnection}
	err := essayRepo.CreateTable()
	assert.Nil(t, err, "essayRepo CreateTable err 발생")
}

func TestAddEssay(t *testing.T) {

	dbConnection, repo, err := prepareTestEssayRepo()
	assert.Nil(t, err, "prepareTestEssayRepo() err 발생")

	defer dbConnection.Close()

	_, err = repo.AddEssay("test essay", "essay thumbail", "essay content", []string{"image1", "image2", "image3"})
	assert.Nil(t, err, "AddEssay() err 발생")

	findEssay, err := repo.FindEssay(1)
	assert.Nil(t, err, "FindEssay() err 발생")

	assert.Equal(t, findEssay.Title, "test essay")
	assert.Equal(t, findEssay.ThumbnailImage, "essay thumbail")
	assert.Equal(t, findEssay.EssayContent, "essay content")
	assert.Equal(t, findEssay.Images[0].Path, "image1")
	assert.Equal(t, findEssay.Images[1].Path, "image2")
	assert.Equal(t, findEssay.Images[2].Path, "image3")
}

func prepareTestExistEssayRepo() (*DBConnection, *EssayRepository, error) {
	dbConnection, repo, err := prepareTestEssayRepo()
	if err != nil {
		return nil, nil, err
	}
	completed := false

	defer func(recvCompleted *bool) {
		if *recvCompleted == false {
			dbConnection.Close()
		}
	}(&completed)

	_, err = repo.AddEssay("test essay1", "essay thumbnail1", "essay content1", []string{"image1-1", "image1-2", "image1-3"})
	if err != nil {
		return nil, nil, err
	}

	_, err = repo.AddEssay("test essay2", "essay thumbnail2", "essay content2", []string{"image2-1", "image2-2", "image2-3"})
	if err != nil {
		return nil, nil, err
	}

	_, err = repo.AddEssay("test essay3", "essay thumbnail3", "essay content3", []string{"image3-1", "image3-2", "image3-3"})
	if err != nil {
		return nil, nil, err
	}

	completed = true

	return dbConnection, repo, nil
}

func TestGetEssaies(t *testing.T) {

	dbConnection, repo, err := prepareTestExistEssayRepo()
	assert.Nil(t, err, "prepareTestExistEssayRepo() error")

	defer dbConnection.Close()

	essies, err := repo.GetEssayList()
	assert.Nil(t, err, "GetPotofolioList() err")

	assert.Equal(t, len(essies), 3)

	assert.Equal(t, essies[0].Title, "test essay1")
	assert.Equal(t, essies[0].ThumbnailImage, "essay thumbnail1")

	assert.Equal(t, essies[1].Title, "test essay2")
	assert.Equal(t, essies[1].ThumbnailImage, "essay thumbnail2")

	assert.Equal(t, essies[2].Title, "test essay3")
	assert.Equal(t, essies[2].ThumbnailImage, "essay thumbnail3")
}

func TestUpdateEssay(t *testing.T) {
	dbConnection, repo, err := prepareTestExistEssayRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	_, err = repo.FindEssay(1)
	assert.Nil(t, err, "FindEssay(1) err")

	imageModel, err := repo.ImageRepo.FindImageFromPath(EssayType, 1, "image1-2")
	assert.Nil(t, err)
	assert.NotNil(t, imageModel)

	updateTitle := "update title"
	updateThumbnail := "update thubmnail"
	updateContent := "update content"

	removeImages, err := repo.UpdateEssay(1, &updateTitle, &updateThumbnail, &updateContent, []int64{imageModel.Id}, []string{"update image1-1"})
	assert.Nil(t, err, "UpdateEssay(1) err")

	assert.Equal(t, len(removeImages), 1)
	assert.Equal(t, removeImages[0], "image1-2")

	essay, _ := repo.FindEssay(1)
	assert.Equal(t, essay.Title, "update title")
	assert.Equal(t, essay.ThumbnailImage, "update thubmnail")
	assert.Equal(t, essay.EssayContent, "update content")

	assert.Equal(t, len(essay.Images), 3)
	assert.Equal(t, essay.Images[0].Path, "image1-1")
	assert.Equal(t, essay.Images[1].Path, "image1-3")
	assert.Equal(t, essay.Images[2].Path, "update image1-1")
}

func TestUpdatePartialEssay_1(t *testing.T) {
	dbConnection, repo, err := prepareTestExistEssayRepo()
	assert.Nil(t, err, "prepareTestEssayRepo() error")

	defer dbConnection.Close()

	_, err = repo.FindEssay(1)
	assert.Nil(t, err, "FindEssay(1) err")

	imageModel, err := repo.ImageRepo.FindImageFromPath(EssayType, 1, "image1-2")
	assert.Nil(t, err)
	assert.NotNil(t, imageModel)

	updateContent := "update content1"

	removeImages, err := repo.UpdateEssay(1, nil, nil, &updateContent, []int64{imageModel.Id}, nil)
	assert.Nil(t, err, "UpdatePotofolio(1) err")

	assert.Equal(t, len(removeImages), 1)
	assert.Equal(t, removeImages[0], "image1-2")

	essay, _ := repo.FindEssay(1)
	assert.Equal(t, essay.Title, "test essay1")
	assert.Equal(t, essay.ThumbnailImage, "essay thumbnail1")
	assert.Equal(t, essay.EssayContent, "update content1")

	assert.Equal(t, len(essay.Images), 2)
	assert.Equal(t, essay.Images[0].Path, "image1-1")
	assert.Equal(t, essay.Images[1].Path, "image1-3")
}

func TestUpdatePartialEssay_2(t *testing.T) {
	dbConnection, repo, err := prepareTestExistEssayRepo()
	assert.Nil(t, err, "prepareTestEssayRepo() error")

	defer dbConnection.Close()

	_, err = repo.FindEssay(1)
	assert.Nil(t, err, "FindEssay(1) err")

	imageModel, err := repo.ImageRepo.FindImageFromPath(EssayType, 1, "image1-2")
	assert.Nil(t, err)
	assert.NotNil(t, imageModel)

	updateTitle := "update title"
	updateThumbnail := "update thumbnail"
	removeImages, err := repo.UpdateEssay(1, &updateTitle, &updateThumbnail, nil, nil, []string{"add Image1-1"})
	assert.Nil(t, err, "UpdatePotofolio(1) err")

	assert.Equal(t, len(removeImages), 0)

	potofolio, _ := repo.FindEssay(1)
	assert.Equal(t, potofolio.Title, updateTitle)
	assert.Equal(t, potofolio.ThumbnailImage, updateThumbnail)

	assert.Equal(t, len(potofolio.Images), 4)
	assert.Equal(t, potofolio.Images[0].Path, "image1-1")
	assert.Equal(t, potofolio.Images[1].Path, "image1-2")
	assert.Equal(t, potofolio.Images[2].Path, "image1-3")
	assert.Equal(t, potofolio.Images[3].Path, "add Image1-1")
}

func TestRemoveEssay(t *testing.T) {

	dbConnection, repo, err := prepareTestExistEssayRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	err = repo.RemoveEssay(2)
	assert.Nil(t, err, "RemoveEssay(2) err")

	findEssay, err := repo.FindEssay(2)
	assert.Nil(t, err, "FindEssay() err")

	assert.Nil(t, findEssay)
}
