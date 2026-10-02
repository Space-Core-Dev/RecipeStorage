package services

type IdConverter[T any] interface {
	Convert(id int64) T
}
