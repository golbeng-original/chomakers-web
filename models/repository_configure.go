package models

import "time"

type RepositoryConfigure struct {
	PotofolioRepository *PotofolioRepository
	EssayRepository     *EssayRepository
	AboutRepository     *AboutRepository
	UserRepository      *UserRespository

	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration

	IsCheckAuthorize bool
}

func (repositoryConfigure *RepositoryConfigure) Init(dbConnection *DBConnection) {

	imageRepository := &ImageRepository{DBConnect: dbConnection}
	imageRepository.CreateTable()

	repositoryConfigure.PotofolioRepository = &PotofolioRepository{DBConnect: dbConnection, ImageRepo: imageRepository}
	repositoryConfigure.PotofolioRepository.CreateTable()

	repositoryConfigure.EssayRepository = &EssayRepository{DBConnect: dbConnection, ImageRepo: imageRepository}
	repositoryConfigure.EssayRepository.CreateTable()

	repositoryConfigure.AboutRepository = &AboutRepository{DBConnect: dbConnection}
	repositoryConfigure.AboutRepository.CreateTable()

	repositoryConfigure.UserRepository = &UserRespository{DBConnect: dbConnection}
	repositoryConfigure.UserRepository.CreateTable()

	repositoryConfigure.AccessTokenExpireTime = 1 * time.Minute
	repositoryConfigure.RefreshTokenExpireTime = 14 * 24 * 60 * time.Minute

	repositoryConfigure.IsCheckAuthorize = true
}
