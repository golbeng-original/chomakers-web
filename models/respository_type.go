package models

type RepositoryType int

const (
	PotofolioType RepositoryType = 1 + iota
	EssayType
	AboutType
)
