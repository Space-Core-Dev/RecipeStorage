package entity

type Id int64

type Entity interface {
	Id() Id
	// any additional methods
}


