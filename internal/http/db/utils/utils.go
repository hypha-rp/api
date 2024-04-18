package utils

type DatabaseOperations interface {
	Create(value interface{}) error
	First(out interface{}, where ...interface{}) error
}
