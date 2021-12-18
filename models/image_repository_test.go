package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getMemoryDbConnect() *DBConnection {
	dbConnect := DBConnection{}
	err := dbConnect.Open("file::memory:?mode=memory&cache=shared")
	if err != nil {
		fmt.Printf("getMemoryDbConnect Error [%v]\n", err)
		return nil
	}

	return &dbConnect
}

func prepareTestImageRepo() (*DBConnection, *ImageRepository, error) {
	dbConnection := getMemoryDbConnect()

	potofolioRepo := &ImageRepository{DBConnect: dbConnection}
	err := potofolioRepo.CreateTable()
	if err != nil {
		return nil, nil, err
	}

	return dbConnection, potofolioRepo, nil
}

func TestImageCreateTable(t *testing.T) {
	dbConnection := getMemoryDbConnect()
	defer dbConnection.Close()

	potofolioRepo := ImageRepository{DBConnect: dbConnection}
	err := potofolioRepo.CreateTable()
	assert.Nil(t, err, "imageRepo CreateTable err 발생")
}

func TestAddImage(t *testing.T) {

	dbConnection, repo, err := prepareTestImageRepo()
	assert.Nil(t, err, "getPotofolioTestPrepare() err 발생")

	defer dbConnection.Close()

	err = repo.AddImges(PotofolioType, 1, []string{"image1", "image2", "image3"})
	assert.Nil(t, err, "AddImges() err 발생")

	findImages, err := repo.GetImages(PotofolioType, 1)
	assert.Nil(t, err, "findImages() err 발생")

	assert.Equal(t, findImages[0].Path, "image1")
	assert.Equal(t, findImages[1].Path, "image2")
	assert.Equal(t, findImages[2].Path, "image3")
}

func prepareTestExistImageRepo() (*DBConnection, *ImageRepository, error) {
	dbConnection, repo, err := prepareTestImageRepo()
	if err != nil {
		return nil, nil, err
	}

	err = repo.AddImges(PotofolioType, 1, []string{"potofolio image1-1", "potofolio image1-2", "potofolio image1-3"})
	if err != nil {
		return nil, nil, err
	}

	err = repo.AddImges(PotofolioType, 2, []string{"potofolio image2-1", "potofolio image2-2", "potofolio image2-3"})
	if err != nil {
		return nil, nil, err
	}

	err = repo.AddImges(EssayType, 1, []string{"essay image2-1", "essay image2-2", "essay image2-3"})
	if err != nil {
		return nil, nil, err
	}

	err = repo.AddImges(EssayType, 2, []string{"essay image2-1", "essay image2-2", "essay image2-3"})
	if err != nil {
		return nil, nil, err
	}

	return dbConnection, repo, nil
}

func TestGetImages(t *testing.T) {

	dbConnection, repo, err := prepareTestExistImageRepo()
	assert.Nil(t, err, "prepareTestExistImageRepo() error")

	defer dbConnection.Close()

	images, err := repo.GetImages(PotofolioType, 1)
	assert.Nil(t, err, "GetImages() err")

	assert.Equal(t, len(images), 3)

	assert.Equal(t, images[0].Path, "potofolio image1-1")
	assert.Equal(t, images[1].Path, "potofolio image1-2")
	assert.Equal(t, images[2].Path, "potofolio image1-3")

	images, err = repo.GetImages(EssayType, 2)
	assert.Nil(t, err, "GetImages() err")

	assert.Equal(t, len(images), 3)

	assert.Equal(t, images[0].Path, "essay image2-1")
	assert.Equal(t, images[1].Path, "essay image2-2")
	assert.Equal(t, images[2].Path, "essay image2-3")
}

func TestRemoveImageAll(t *testing.T) {

	dbConnection, repo, err := prepareTestExistImageRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	err = repo.RemoveImages(EssayType, 1)
	assert.Nil(t, err, "RemoveImages(models.EssayType, 1) err")

	images, err := repo.GetImages(EssayType, 1)
	assert.Nil(t, err, "FindPotofolio() err")

	assert.Equal(t, len(images), 0)
}

func TestRemoveImageElement(t *testing.T) {
	dbConnection, repo, err := prepareTestExistImageRepo()
	assert.Nil(t, err, "prepareTestExistPotofolioRepo() error")

	defer dbConnection.Close()

	err = repo.RemoveImageFromPath(EssayType, 1, "essay image2-2")
	assert.Nil(t, err, "RemoveImages(models.EssayType, 1) err")

	images, err := repo.GetImages(EssayType, 1)
	assert.Nil(t, err, "FindPotofolio() err")

	assert.Equal(t, len(images), 2)
}
